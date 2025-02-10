package manager

import (
	"flag"
	"fmt"

	gatusiov1alpha1 "github.com/aumer-amr/gatus-operator/v2/api/v1alpha1"
	config "github.com/aumer-amr/gatus-operator/v2/internal/gatus-operator/config"
	"github.com/aumer-amr/gatus-operator/v2/internal/gatus-operator/controller"
	corev1 "k8s.io/api/core/v1"
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

func Run() error {
	logger.Info("setting up")

	configuration := config.Generate()

	if configuration.DevMode {
		logger.Info("running in dev mode, setting up local kubeconfig")

		var kubeConfig string
		flagSet := flag.NewFlagSet("kubeconfig", flag.ExitOnError)
		flagSet.StringVar(&kubeConfig, "kubeconfig", "../../kubeconfig", "Path to the kubeconfig file to use for CLI requests.")

		ctrl.RegisterFlags(flagSet)
	}

	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: configuration.ProbeAddr,
		LeaderElection:         false,
		Metrics: server.Options{
			BindAddress: configuration.MetricsAddr,
		},
	})
	if err != nil {
		logger.Error(err, "unable to start")
		return err
	}

	if err := manager.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		panic(fmt.Errorf("unable to add healthz check: %w", err))
	}
	if err := manager.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		panic(fmt.Errorf("unable to add readyz check: %w", err))
	}

	logger.Info(fmt.Sprintf("endpoint defaults found: %v", configuration.HasDefaults()))

	err = ctrl.NewControllerManagedBy(manager).
		For(&gatusiov1alpha1.Gatus{}).
		Owns(&gatusiov1alpha1.Gatus{}).
		Owns(&corev1.ConfigMap{}).
		Complete(&controller.ReconcileGatus{
			Client: manager.GetClient(),
			Config: configuration,
		})
	if err != nil {
		logger.Error(err, "unable to setup controller")
		return err
	}

	logger.Info("starting")
	if err := manager.Start(ctrl.SetupSignalHandler()); err != nil {
		logger.Error(err, "problem running")
		return err
	}
	return nil
}
