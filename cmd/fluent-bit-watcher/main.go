package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/sync/errgroup"
	"kubesphere.io/fluentbit-operator/pkg/filenotify"
)

const (
	defaultBinPath      = "/fluent-bit/bin/fluent-bit"
	defaultCfgPath      = "/fluent-bit/etc/fluent-bit.conf"
	defaultWatchDir     = "/fluent-bit/config"
	defaultPollInterval = 1 * time.Second
)

func main() {
	var configPath string
	var binPath string
	var watchPath string
	var poll bool
	var pollInterval time.Duration
	flag.StringVar(&binPath, "b", defaultBinPath, "The fluent bit binary path.")
	flag.StringVar(&configPath, "c", defaultCfgPath, "The config file path.")
	flag.StringVar(&watchPath, "watch-path", defaultWatchDir, "The path to watch.")
	flag.BoolVar(&poll, "poll", false, "Use poll watcher instead of inotify.")
	flag.DurationVar(&pollInterval, "poll-interval", defaultPollInterval, "Poll interval if using poll watcher.")

	flag.Parse()

	signalCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// First, launch the fluent-bit process.
	cmd := exec.Command(binPath, "--enable-hot-reload", "-c", configPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to start fluent-bit: %v", err)
	}

	grp, grpCtx := errgroup.WithContext(context.Background())
	grp.Go(func() error {
		// Watch the process. If it exits, we want to crash immediately.
		defer cancel()
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("failed to run fluent-bit: %w", err)
		}
		return nil
	})
	grp.Go(func() error {
		// Watch the config as it's loaded into the pod and trigger a config reload.
		var watcher filenotify.FileWatcher
		if poll {
			watcher = filenotify.NewPollingWatcher(pollInterval)
		} else {
			var err error
			watcher, err = filenotify.NewEventWatcher()
			if err != nil {
				return fmt.Errorf("failed to open event watcher: %w", err)
			}
		}

		if err := watcher.Add(watchPath); err != nil {
			return fmt.Errorf("failed to watch path %q: %w", watchPath, err)
		}

		for {
			select {
			case <-signalCtx.Done():
				return nil
			case <-grpCtx.Done():
				return nil
			case event := <-watcher.Events():
				if !isValidEvent(event) {
					continue
				}
				log.Print("Config file changed, reloading...")
				if err := cmd.Process.Signal(syscall.SIGHUP); err != nil {
					return fmt.Errorf("failed to reload config: %w", err)
				}
			case err := <-watcher.Errors():
				return fmt.Errorf("failed the watcher: %w", err)
			}
		}
	})

	select {
	case <-signalCtx.Done():
	case <-grpCtx.Done():
	}

	// Always try to gracefully shut down fluent-bit. This will allow `cmd.Wait` above to finish
	// and thus allow `grp.Wait` below to return.
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
		log.Printf("Failed to send SIGTERM to fluent-bit: %v", err)
		// Do not exit on error here. The process might've died and that's okay.
	}

	if err := grp.Wait(); err != nil {
		log.Fatalf("Failure during the run time of fluent-bit: %v", err)
	}
}

// Inspired by https://github.com/jimmidyson/configmap-reload
func isValidEvent(event fsnotify.Event) bool {
	return event.Op == fsnotify.Create || event.Op == fsnotify.Write
}
