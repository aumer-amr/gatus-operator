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

func (r *ReconcileGatus) GatusReconcile(gatus *gatusiov1alpha1.Gatus) error {
	logger.Info("reconciling Gatus", "name", gatus.Name, "namespace", gatus.Namespace)

	configMapList := corev1.ConfigMapList{}
	err := r.Manager.GetClient().List(context.Background(), &configMapList, client.MatchingLabels{
		"app.kubernetes.io/managed-by": "gatus.io",
		"gatus.io/parent-uid":          string(gatus.ObjectMeta.UID),
	})

	if err != nil {
		return err
	}

	if !gatus.ObjectMeta.DeletionTimestamp.IsZero() {
		if len(configMapList.Items) == 1 {
			configMap := configMapList.Items[0]
			return r.Delete(gatus, configMap)
		}
		return nil
	}

	if len(configMapList.Items) == 0 {
		return r.Create(gatus)
	}

	if len(configMapList.Items) == 1 {
		configMap := configMapList.Items[0]
		return r.Update(gatus, configMap)
	}

	return nil
}

func (r *ReconcileGatus) Create(gatus *gatusiov1alpha1.Gatus) error {
	yamlString, err := r.GetEndpointsYaml(gatus)
	if err != nil {
		return fmt.Errorf("error getting YAML: %w", err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: r.GenerateMetaData(gatus),
		Data: map[string]string{
			"gatus.yaml": yamlString,
		},
	}

	r.Manager.GetClient().Create(context.Background(), configMap)

	return nil
}

func (r *ReconcileGatus) Update(gatus *gatusiov1alpha1.Gatus, configMap corev1.ConfigMap) error {
	yamlString, err := r.GetEndpointsYaml(gatus)
	if err != nil {
		return fmt.Errorf("error getting YAML: %w", err)
	}

	configMap.Data["gatus.yaml"] = yamlString
	r.Manager.GetClient().Update(context.Background(), &configMap)

	return nil
}

func (r *ReconcileGatus) Delete(gatus *gatusiov1alpha1.Gatus, configMap corev1.ConfigMap) error {
	r.Manager.GetClient().Delete(context.Background(), &configMap)

	return nil
}

func (r *ReconcileGatus) GetEndpointsYaml(gatus *gatusiov1alpha1.Gatus) (string, error) {
	type GatusEndpoint struct {
		Endpoints []gatusiov1alpha1.EndpointEndpoint `json:"endpoints"`
	}

	yamlBytes, err := yaml.Marshal(GatusEndpoint{Endpoints: []gatusiov1alpha1.EndpointEndpoint{gatus.Spec.Endpoint}})
	if err != nil {
		return "", fmt.Errorf("error converting struct to YAML: %v", err)
	}

	yamlString := string(yamlBytes)

	return yamlString, nil
}

func (r *ReconcileGatus) GenerateMetaData(gatus *gatusiov1alpha1.Gatus) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", gatus.Name, "gatus-config"),
		Namespace: gatus.Namespace,
		Labels: map[string]string{
			"app.kubernetes.io/managed-by": "gatus.io",
			"gatus.io/enabled":             "enabled",
			"gatus.io/parent-uid":          string(gatus.ObjectMeta.UID),
		},
	}
}
