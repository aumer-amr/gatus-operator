package controller

import (
	"context"
	"os"

	gatusiov1alpha1 "github.com/aumer-amr/gatus-operator/v2/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	logger = ctrl.Log.WithName("controller")
)

type ReconcileGatus struct {
	Client     client.Client
	Manager    ctrl.Manager
	Controller controller.Controller
}

func Run(mgr ctrl.Manager) {
	logger.Info("setting up controller")

	controller, err := controller.New("gatus-controller", mgr, controller.Options{
		Reconciler: &ReconcileGatus{
			Client:  mgr.GetClient(),
			Manager: mgr,
		},
	})

	if err != nil {
		logger.Error(err, "unable to set up controller")
		os.Exit(1)
	}

	rg := &ReconcileGatus{
		Manager:    mgr,
		Client:     mgr.GetClient(),
		Controller: controller,
	}

	rg.WatchResource(mgr, &gatusiov1alpha1.Gatus{})
}

func (r *ReconcileGatus) WatchResource(mgr ctrl.Manager, resourceType client.Object) {
	if err := r.Controller.Watch(source.Kind(mgr.GetCache(), resourceType, &handler.EnqueueRequestForObject{})); err != nil {
		logger.Error(err, "unable to watch resource")
		os.Exit(1)
	}

	if err := r.Controller.Watch(source.Kind(mgr.GetCache(), resourceType,
		handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), resourceType, handler.OnlyControllerOwner()))); err != nil {
		logger.Error(err, "unable to watch Pods")
		os.Exit(1)
	}
}

func (r *ReconcileGatus) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger.Info("reconciling", "request", req)

	gatus := &gatusiov1alpha1.Gatus{}
	err, cacheMiss := r.CheckCache(ctx, req.NamespacedName.String(), gatus)

	if err != nil {
		if cacheMiss {
			logger.Error(err, "unable to fetch Ingress from cache")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "unable to fetch Ingress")
		return ctrl.Result{}, err
	}

	r.GatusReconcile(gatus)

	return ctrl.Result{}, nil
}

func (r *ReconcileGatus) CheckCache(ctx context.Context, namespacedName string, typed client.Object) (error, bool) {
	err := r.Client.Get(ctx, client.ObjectKey{Name: namespacedName}, typed)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Cache miss", "namespacedName", namespacedName)
			return err, true
		}
		logger.Error(err, "Failed to get object from cache", "namespacedName", namespacedName)
		return err, false
	}

	return nil, false
}
