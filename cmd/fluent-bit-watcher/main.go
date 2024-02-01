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
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// The main config file includes this file. This is where we write our generated config.
	configFileName = "/fluent-bit/config/fluent-bit.conf"
	// The main parsers file includes this file. This is where we write our generated config.
	parsersFileName = "/fluent-bit/config/parsers.conf"
)

func main() {
	ctrl.SetLogger(klog.NewKlogr())

	scheme := runtime.NewScheme()
	utilruntime.Must(v1alpha2.AddToScheme(scheme))

	ns := os.Getenv("NAMESPACE")
	if ns == "" {
		log.Fatal("NAMESPACE environment variable must be set")
	}

	signalCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	kubeConfig, err := ctrl.GetConfig()
	if err != nil {
		log.Fatalf("Failed to get kubernetes config: %v", err)
	}

	// First, generate a config from the current state in the cluster.
	nonCachedClient, err := client.New(kubeConfig, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalf("Failed to create non-cached client: %v", err)
	}
	mainConfig, parserConfig, err := regenerateConfigs(signalCtx, nonCachedClient, ns)
	if err != nil {
		log.Fatalf("Failed to generate initial config: %v", err)
	}
	if err := overrideFile(configFileName, mainConfig); err != nil {
		log.Fatalf("Failed to write initial config file: %v", err)
	}
	if err := overrideFile(parsersFileName, parserConfig); err != nil {
		log.Fatalf("Failed to write initial parser file: %v", err)
	}

	// Launch the actual fluent-bit process.
	cmd := exec.Command("/fluent-bit/bin/fluent-bit", "--enable-hot-reload", "-c", "/fluent-bit/etc/fluent-bit.conf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start fluent-bit: %v", err)
	}

	// Create the controller that'll watch for config changes and cause a
	// config reload.
	mgr, err := ctrl.NewManager(kubeConfig, ctrl.Options{
		Scheme: scheme,
		Cache: cache.Options{
			DefaultNamespaces: map[string]cache.Config{ns: {}},
		},
		Metrics: server.Options{BindAddress: "0"},
	})
	if err != nil {
		log.Fatalf("Failed to create new manager: %v", err)
	}

	singletonEnqueue := handler.EnqueueRequestsFromMapFunc(func(context.Context, client.Object) []reconcile.Request {
		return []reconcile.Request{{}}
	})
	if err := ctrl.NewControllerManagedBy(mgr).
		Named("config").
		Watches(&v1alpha2.Input{}, singletonEnqueue).
		Watches(&v1alpha2.Filter{}, singletonEnqueue).
		Watches(&v1alpha2.Output{}, singletonEnqueue).
		Watches(&v1alpha2.Parser{}, singletonEnqueue).
		Complete(&configReconciler{
			fluentbitProcess: cmd,
			apiClient:        mgr.GetClient(),
			namespace:        ns,
			lastMainConfig:   mainConfig,
			lastParserConfig: parserConfig,
		}); err != nil {
		log.Fatalf("Failed to create controller: %v", err)
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
		if err := mgr.Start(signalCtx); err != nil {
			return fmt.Errorf("failed to run manager: %w", err)
		}
		return nil
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

// configReconciler watches inputs, filters, outputs and parsers and regenerates the respective
// config files on changes. If there were any changes, the fluent-bit process is kicked to reload
// its configuration.
type configReconciler struct {
	fluentbitProcess *exec.Cmd
	apiClient        client.Client
	namespace        string

	lastMainConfig   []byte
	lastParserConfig []byte
}

func (r *configReconciler) Reconcile(ctx context.Context, _ reconcile.Request) (reconcile.Result, error) {
	log.Print("Reconciling configuration")

	mainConfig, parserConfig, err := regenerateConfigs(ctx, r.apiClient, r.namespace)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to regenerate configs: %w", err)
	}

	var changed bool
	if !bytes.Equal(mainConfig, r.lastMainConfig) {
		if err := overrideFile(configFileName, mainConfig); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to override config file: %w", err)
		}
		changed = true
		r.lastMainConfig = mainConfig
	}
	if !bytes.Equal(parserConfig, r.lastParserConfig) {
		if err := overrideFile(parsersFileName, parserConfig); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to override parser file: %w", err)
		}
		changed = true
		r.lastParserConfig = parserConfig
	}

	if !changed {
		log.Print("No changes to the config files")
		return reconcile.Result{}, nil
	}

	log.Print("Reloading fluent-bit config")
	if err := r.fluentbitProcess.Process.Signal(syscall.SIGHUP); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to signal fluent-bit process: %w", err)
	}
	return reconcile.Result{}, nil
}

// regenerateConfigs fetches the current in-cluster state and generates a new
// main and parsers config.
func regenerateConfigs(ctx context.Context, apiClient client.Client, ns string) ([]byte, []byte, error) {
	var inputs v1alpha2.InputList
	if err := apiClient.List(ctx, &inputs, client.InNamespace(ns)); err != nil {
		return nil, nil, fmt.Errorf("failed to list inputs: %w", err)
	}

	var filters v1alpha2.FilterList
	if err := apiClient.List(ctx, &filters, client.InNamespace(ns)); err != nil {
		return nil, nil, fmt.Errorf("failed to list filters: %w", err)
	}

	var outputs v1alpha2.OutputList
	if err := apiClient.List(ctx, &outputs, client.InNamespace(ns)); err != nil {
		return nil, nil, fmt.Errorf("failed to list outputs: %w", err)
	}

	var parsers v1alpha2.ParserList
	if err := apiClient.List(ctx, &parsers, client.InNamespace(ns)); err != nil {
		return nil, nil, fmt.Errorf("failed to list parsers: %w", err)
	}

	// Note: We're not using the secretloader for anything right now, so this is a noop.
	// We keep it in the call anyway to facilitate some level of compatibility to the
	// upstream CRD code.
	sl := plugins.NewSecretLoader(apiClient, ns)
	mainConfig, err := renderMainConfig(sl, inputs, filters, outputs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to render main config: %w", err)
	}

	parserSections, err := parsers.Load(sl)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to render parser config: %w", err)
	}

	return mainConfig, []byte(parserSections), nil
}

// renderMainConfig renders the main fluent-bit config from inputs, filters and outputs.
func renderMainConfig(
	sl plugins.SecretLoader,
	inputs v1alpha2.InputList,
	filters v1alpha2.FilterList,
	outputs v1alpha2.OutputList,
) ([]byte, error) {
	var buf bytes.Buffer

	inputSections, err := inputs.Load(sl)
	if err != nil {
		return nil, fmt.Errorf("failed to render input sections: %w", err)
	}

	filterSections, err := filters.Load(sl)
	if err != nil {
		return nil, fmt.Errorf("failed to render filter sections: %w", err)
	}

	outputSections, err := outputs.Load(sl)
	if err != nil {
		return nil, fmt.Errorf("failed to render output sections: %w", err)
	}

	if inputSections != "" && outputSections == "" {
		outputSections = `[Output]
    Name    null
    Match   *`
	}

	buf.WriteString(`[Service]
    Daemon    false
    Flush    1
    Grace    60
    Http_Server    true
    Log_Level    warning
    Parsers_File    /fluent-bit/etc/parsers.conf
`)
	buf.WriteString(inputSections)
	buf.WriteString(filterSections)
	buf.WriteString(outputSections)

	return buf.Bytes(), nil
}

// overrideFile creates or truncates the file at the given name and writes the given content to it.
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
