package controller

import (
	"context"

	gatusiov1alpha1 "github.com/aumer-amr/gatus-operator/v2/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

var logger = ctrl.Log.WithName("controller")

type ReconcileGatus struct {
	client.Client
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

	r.GatusReconcile(ctx, gatus)

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
