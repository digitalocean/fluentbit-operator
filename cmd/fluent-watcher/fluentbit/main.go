package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fluent/fluent-operator/v2/pkg/copy"
	"github.com/fluent/fluent-operator/v2/pkg/filenotify"
	"github.com/fluent/fluent-operator/v2/pkg/gziputil"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
)

const (
	defaultBinPath  = "/fluent-bit/bin/fluent-bit"
	defaultCfgPath  = "/fluent-bit/etc/fluent-bit.conf"
	defaultWatchDir = "/fluent-bit/config"
	// decompressed config in scratch space
	scratchCfgDir  = "/tmp/fluent-bit"
	scratchCfgFile = "fluent-bit.scratch.conf"

	defaultPollInterval = 1 * time.Second
	defaultFlbTimeout   = 30 * time.Second

	MaxDelayTime = 5 * time.Minute
	ResetTime    = 10 * time.Minute
)

var (
	logger        log.Logger
	cmd           *exec.Cmd
	flbTerminated chan bool
	mutex         sync.Mutex
	restartTimes  int32
	timerCtx      context.Context
	timerCancel   context.CancelFunc
)

var configPath string
var externalPluginPath string
var binPath string
var watchPath string
var poll bool
var exitOnFailure bool
var pollInterval time.Duration
var flbTerminationTimeout time.Duration

func main() {

	flag.StringVar(&binPath, "b", defaultBinPath, "The fluent bit binary path.")
	flag.StringVar(&configPath, "c", defaultCfgPath, "The config file path.")
	flag.StringVar(&externalPluginPath, "e", "", "Path to external plugin (shared lib)")
	flag.BoolVar(&exitOnFailure, "exit-on-failure", false, "If fluentbit exits with failure, also exit the watcher.")
	flag.StringVar(&watchPath, "watch-path", defaultWatchDir, "The path to watch.")
	flag.BoolVar(&poll, "poll", false, "Use poll watcher instead of ionotify.")
	flag.DurationVar(&pollInterval, "poll-interval", defaultPollInterval, "Poll interval if using poll watcher.")
	flag.DurationVar(&flbTerminationTimeout, "flb-timeout", defaultFlbTimeout, "Time to wait for FluentBit to gracefully terminate before sending SIGKILL.")

	flag.Parse()

	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "time", log.TimestampFormat(time.Now, time.RFC3339))

	timerCtx, timerCancel = context.WithCancel(context.Background())

	var g run.Group
	{
		// Termination handler.
		g.Add(run.SignalHandler(context.Background(), os.Interrupt, syscall.SIGTERM))
	}
	{
		// Watch the Fluent bit, if the Fluent bit not exists or stopped, restart it.
		cancel := make(chan struct{})
		g.Add(
			func() error {

				for {
					select {
					case <-cancel:
						return nil
					default:
					}

					// Start fluent bit if it does not existed.
					start()
					// Wait for the fluent bit exit.
					err := wait()
					if exitOnFailure && err != nil {
						_ = level.Error(logger).Log("msg", "Fluent bit exited with error; exiting watcher")
						return err
					}

					timerCtx, timerCancel = context.WithCancel(context.Background())

					// After the fluent bit exit, fluent bit watcher restarts it with an exponential
					// back-off delay (1s, 2s, 4s, ...), that is capped at five minutes.
					backoff()
				}
			},
			func(err error) {
				close(cancel)
				stop()
				resetTimer()
			},
		)
	}
	{
		// Watch the config file, if the config file changed, stop Fluent bit.
		watcher, err := newWatcher(poll, pollInterval)
		if err != nil {
			_ = level.Error(logger).Log("err", err)
			return
		}

		// Start watcher.
		err = watcher.Add(watchPath)
		if err != nil {
			_ = level.Error(logger).Log("err", err)
			return
		}

		cancel := make(chan struct{})
		g.Add(
			func() error {

				for {
					select {
					case <-cancel:
						return nil
					case event := <-watcher.Events():
						if !isValidEvent(event) {
							continue
						}
						// After the config file changed, it should stop the fluent bit,
						// and resets the restart backoff timer.
						if cmd != nil {
							_ = level.Info(logger).Log("msg", "Config file changed, stopping Fluent Bit")
							stop()
							resetTimer()
							_ = level.Info(logger).Log("msg", "Config file changed, stopped Fluent Bit")
						}
					case <-watcher.Errors():
						_ = level.Error(logger).Log("msg", "Watcher stopped")
						return nil
					}
				}
			},
			func(err error) {
				_ = watcher.Close()
				close(cancel)
			},
		)
	}

	if err := g.Run(); err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	_ = level.Info(logger).Log("msg", "See you next time!")
}

func newWatcher(poll bool, interval time.Duration) (filenotify.FileWatcher, error) {
	var err error
	var watcher filenotify.FileWatcher

	if poll {
		watcher = filenotify.NewPollingWatcher(interval)
	} else {
		watcher, err = filenotify.New(interval)
	}

	if err != nil {
		return nil, err
	}

	return watcher, nil
}

// Inspired by https://github.com/jimmidyson/configmap-reload
func isValidEvent(event fsnotify.Event) bool {
	return event.Op == fsnotify.Create || event.Op == fsnotify.Write
}

func start() {
	mutex.Lock()
	defer mutex.Unlock()

	if cmd != nil {
		return
	}

	var err error

	isCustomConfigSet := !(defaultCfgPath == configPath)
	configFilePath := configPath
	compressed := false

	fileToCheck := configFilePath

	if !isCustomConfigSet {
		// if default is used we should check the fluent-bit.conf in watch dir (/fluent-bit/config/)
		// default is /fluent-bit/etc/fluent-bit.conf
		// which comes from Dockerfile which is copied from the conf source folder
		// and all it does is @INCLUDE /fluent-bit/config/fluent-bit.conf

		fileToCheck = filepath.Join(defaultWatchDir, "fluent-bit.conf")
	}

	compressed, err = gziputil.IsCompressed(fileToCheck)
	if err != nil {
		_ = level.Error(logger).Log("msg", "checking "+fileToCheck, "error", err)
		return
	}

	if compressed {
		// there may be references in service config to local relative files (ie parsers.conf)
		// we should copy all to scratch space first
		if err := copy.CopyFilesWithFilter("/fluent-bit/etc/", scratchCfgDir, copyFilterfunc); err != nil {
			_ = level.Error(logger).Log("msg", "copying parsers config file", "error", err)
			return
		}

		if isCustomConfigSet {
			// if custom config used, also copy files in folder of custom config
			baseDir := filepath.Dir(fileToCheck)
			if err := copy.CopyFilesWithFilter(baseDir, scratchCfgDir, copyFilterfunc); err != nil {
				_ = level.Error(logger).Log("msg", "copying parsers config file", "error", err)
				return
			}
		}

		// decompress
		newConfig := filepath.Join(scratchCfgDir, scratchCfgFile)

		_ = level.Info(logger).Log("msg", fmt.Sprintf("%s is compressed. Decompressing to %s", fileToCheck, newConfig))

		if err := gziputil.Decompress(fileToCheck, newConfig); err != nil {
			_ = level.Error(logger).Log("msg", "decompress", "error", err)
			return
		}

		// finally set it
		configFilePath = newConfig
	} else {
		_ = level.Info(logger).Log("msg", fmt.Sprintf("%s is not compressed.", configFilePath))
	}

	if externalPluginPath != "" {
		cmd = exec.Command(binPath, "-c", configFilePath, "-e", externalPluginPath)
	} else {
		cmd = exec.Command(binPath, "-c", configFilePath)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	flbTerminated = make(chan bool, 1)
	if err := cmd.Start(); err != nil {
		_ = level.Error(logger).Log("msg", "start Fluent bit error", "error", err)
		cmd = nil
		return
	}

	_ = level.Info(logger).Log("msg", "Fluent bit started")
}

func wait() error {
	mutex.Lock()
	if cmd == nil {
		mutex.Unlock()
		return nil
	}
	mutex.Unlock()

	startTime := time.Now()

	err := cmd.Wait()
	if err != nil {
		_ = level.Error(logger).Log("msg", "Fluent bit exited", "error", err)
	}
	cmd = nil
	flbTerminated <- true

	// Once the fluent bit has executed for 10 minutes without any problems,
	// it should resets the restart backoff timer.
	if time.Since(startTime) >= ResetTime {
		atomic.StoreInt32(&restartTimes, 0)
	}

	return err
}

func backoff() {

	delayTime := time.Duration(math.Pow(2, float64(atomic.LoadInt32(&restartTimes)))) * time.Second
	if delayTime >= MaxDelayTime {
		delayTime = MaxDelayTime
	}

	_ = level.Info(logger).Log("msg", "backoff", "delay", delayTime)

	startTime := time.Now()

	timer := time.NewTimer(delayTime)
	defer timer.Stop()

	select {
	case <-timerCtx.Done():
		_ = level.Info(logger).Log("msg", "context cancel", "actual", time.Since(startTime), "expected", delayTime)

		atomic.StoreInt32(&restartTimes, 0)

		return
	case <-timer.C:
		_ = level.Info(logger).Log("msg", "backoff timer done", "actual", time.Since(startTime), "expected", delayTime)

		atomic.AddInt32(&restartTimes, 1)

		return
	}
}

func stop() {

	mutex.Lock()
	defer mutex.Unlock()

	if cmd == nil || cmd.Process == nil {
		_ = level.Info(logger).Log("msg", "Fluent Bit not running. No process to stop.")
		return
	}

	// Send SIGTERM, if fluent-bit doesn't terminate in the specified timeframe, send SIGKILL
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		_ = level.Info(logger).Log("msg", "Error while terminating FluentBit", "error", err)
	} else {
		_ = level.Info(logger).Log("msg", "Sent SIGTERM to FluentBit, waiting max "+flbTerminationTimeout.String())
	}

	select {
	case <-time.After(flbTerminationTimeout):
		_ = level.Info(logger).Log("msg", "FluentBit failed to terminate gracefully, killing process")
		cmd.Process.Kill()
		<-flbTerminated
	case <-flbTerminated:
		_ = level.Info(logger).Log("msg", "FluentBit terminated successfully")
	}
}

func resetTimer() {
	timerCancel()
	atomic.StoreInt32(&restartTimes, 0)
}

func copyFilterfunc(fi os.FileInfo) bool {
	return strings.HasSuffix(fi.Name(), ".conf") ||
		strings.HasSuffix(fi.Name(), ".lua") ||
		strings.HasSuffix(fi.Name(), ".yaml") ||
		strings.HasSuffix(fi.Name(), ".yml")
}
