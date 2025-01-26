package manager

import (
	"flag"
	"fmt"
	"os"

	gatusiov1alpha1 "github.com/aumer-amr/gatus-operator/v2/api/v1alpha1"
	"github.com/aumer-amr/gatus-operator/v2/internal/gatus-operator/config"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	logger = ctrl.Log.WithName("manager")
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(gatusiov1alpha1.AddToScheme(scheme))
}

func Run() ctrl.Manager {
	logger.Info("setting up manager")

	config := config.Generate()

	if config.DevMode {
		logger.Info("running in dev mode, setting up local kubeconfig")

		var kubeConfig string
		flagSet := flag.NewFlagSet("kubeconfig", flag.ExitOnError)
		flagSet.StringVar(&kubeConfig, "kubeconfig", "../../kubeconfig", "Path to the kubeconfig file to use for CLI requests.")

		ctrl.RegisterFlags(flagSet)
	}

	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: config.ProbeAddr,
		LeaderElection:         false,
		Metrics: server.Options{
			BindAddress: config.MetricsAddr,
		},
	})

	if err != nil {
		logger.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err := manager.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		panic(fmt.Errorf("unable to add healthz check: %w", err))
	}
	if err := manager.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		panic(fmt.Errorf("unable to add readyz check: %w", err))
	}

	logger.Info("starting manager")
	if err := manager.Start(ctrl.SetupSignalHandler()); err != nil {
		logger.Error(err, "problem running manager")
		os.Exit(1)
	}

	return manager
}
