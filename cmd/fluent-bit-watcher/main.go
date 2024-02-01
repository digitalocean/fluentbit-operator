package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/runtime"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// TODO: Make this a flag or pass it via env.
	ns = "apps-external-log"

	// The main config file includes this file. This is where we write our generated config.
	configFileName = "/fluent-bit/config/fluent-bit.conf"
	// The main parsers file includes this file. This is where we write our generated config.
	parsersFileName = "/fluent-bit/config/parsers.conf"
)

func main() {
	signalCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// First, launch the fluent-bit process.
	cmd := exec.Command("/fluent-bit/bin/fluent-bit", "--enable-hot-reload", "-c", "/fluent-bit/etc/fluent-bit.conf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start fluent-bit: %v", err)
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

	scheme := runtime.NewScheme()
	v1alpha2.AddToScheme(scheme)
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Cache: cache.Options{
			DefaultNamespaces: map[string]cache.Config{ns: {}},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create new manager: %v", err)
	}
	apiClient := mgr.GetClient()
	ctrl.NewControllerManagedBy(mgr).
		Watches(&v1alpha2.Input{}, &handler.EnqueueRequestForObject{}).
		Watches(&v1alpha2.Filter{}, &handler.EnqueueRequestForObject{}).
		Watches(&v1alpha2.Output{}, &handler.EnqueueRequestForObject{}).
		Watches(&v1alpha2.Parser{}, &handler.EnqueueRequestForObject{}).
		Complete(reconcile.Func(func(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
			var inputs v1alpha2.InputList
			if err = apiClient.List(ctx, &inputs, client.InNamespace(ns)); err != nil {
				return ctrl.Result{}, err
			}

			var filters v1alpha2.FilterList
			if err = apiClient.List(ctx, &filters, client.InNamespace(ns)); err != nil {
				return ctrl.Result{}, err
			}

			var outputs v1alpha2.OutputList
			if err = apiClient.List(ctx, &outputs, client.InNamespace(ns)); err != nil {
				return ctrl.Result{}, err
			}

			var parsers v1alpha2.ParserList
			if err = apiClient.List(ctx, &parsers, client.InNamespace(ns)); err != nil {
				return ctrl.Result{}, err
			}

			// TODO: Get rid of this? Maybe it's fine?
			sl := plugins.NewSecretLoader(apiClient, ns)
			mainConfig, err := renderMainConfig(sl, inputs, filters, outputs)
			if err != nil {
				return ctrl.Result{}, err
			}

			parserConfig, err := renderParserConfig(sl, parsers)
			if err != nil {
				return ctrl.Result{}, err
			}

			if err := overrideFile(configFileName, mainConfig); err != nil {
				return ctrl.Result{}, err
			}
			if err := overrideFile(parsersFileName, parserConfig); err != nil {
				return ctrl.Result{}, err
			}

			return reconcile.Result{}, nil
		}))

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

func renderMainConfig(
	sl plugins.SecretLoader,
	inputs v1alpha2.InputList,
	filters v1alpha2.FilterList,
	outputs v1alpha2.OutputList,
) ([]byte, error) {
	var buf bytes.Buffer

	// The Service defines the global behaviour of the Fluent Bit engine.
	buf.WriteString(`[Service]
    Daemon    false
    Flush    1
    Grace    60
    Http_Server    true
    Log_Level    warning
    Parsers_File    /fluent-bit/etc/parsers.conf`)

	inputSections, err := inputs.Load(sl)
	if err != nil {
		return nil, err
	}

	filterSections, err := filters.Load(sl)
	if err != nil {
		return nil, err
	}

	outputSections, err := outputs.Load(sl)
	if err != nil {
		return nil, err
	}

	if inputSections != "" && outputSections == "" {
		outputSections = `[Output]
    Name    null
    Match   *`
	}

	buf.WriteString(inputSections)
	buf.WriteString(filterSections)
	buf.WriteString(outputSections)

	return buf.Bytes(), nil
}

func renderParserConfig(sl plugins.SecretLoader, parsers v1alpha2.ParserList) ([]byte, error) {
	parserSections, err := parsers.Load(sl)
	if err != nil {
		return nil, err
	}

	return []byte(parserSections), nil
}

func overrideFile(fileName string, content []byte) error {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		return fmt.Errorf("failed to write content to file: %w", err)
	}
	return nil
}
