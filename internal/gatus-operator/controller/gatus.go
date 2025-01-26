package controller

import (
	"context"
	"fmt"

	gatusiov1alpha1 "github.com/aumer-amr/gatus-operator/v2/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

var (
	baseAnnotation = "gatus.io/"
)

func (r *ReconcileGatus) GatusReconcile(gatus *gatusiov1alpha1.Gatus) {
	logger.Info("reconciling Gatus", "name", gatus.Name, "namespace", gatus.Namespace)

	configMapName := fmt.Sprintf("%s-gatus-config", gatus.Name)
	configMapExists, err := r.checkConfigMapExists(context.Background(), gatus.Namespace, configMapName)

	if err != nil {
		logger.Error(err, "error checking if ConfigMap exists")
		return
	}

	if !configMapExists {
		logger.Info("ConfigMap does not exist, creating", "name", configMapName)
		err := r.createConfigMap(gatus, configMapName)
		if err != nil {
			logger.Error(err, "error creating ConfigMap", "error")
			return
		}
	}
}

func (r *ReconcileGatus) createConfigMap(gatus *gatusiov1alpha1.Gatus, configMapName string) error {
	yamlBytes, err := yaml.Marshal(gatus.Spec.Endpoint)
	if err != nil {
		return fmt.Errorf("Error converting struct to YAML: %v\n", err)
	}

	yamlString := string(yamlBytes)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: gatus.Namespace,
		},
		Data: map[string]string{
			"gatus.yaml": yamlString,
		},
	}

	err = r.Client.Create(context.Background(), configMap)
	if err != nil {
		return fmt.Errorf("error creating ConfigMap: %w", err)
	}

	return nil
}

func (r *ReconcileGatus) checkConfigMapExists(ctx context.Context, namespace string, name string) (bool, error) {
	configMap := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, configMap)

	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			return false, nil
		}
		return false, fmt.Errorf("error checking ConfigMap: %w", err)
	}
	return true, nil
}
