package controller

import (
	"context"
	"fmt"

	gatusiov1alpha1 "github.com/aumer-amr/gatus-operator/v2/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/yaml"
)

const FINALIZER_NAME = "gatus-operator.aumer.io/finalizer"

func (r *ReconcileGatus) gatusReconcile(ctx context.Context, gatus *gatusiov1alpha1.Gatus) error {
	logger.Info("reconciling Gatus", "name", gatus.Name, "namespace", gatus.Namespace)

	configMapList := corev1.ConfigMapList{}
	err := r.List(ctx, &configMapList, client.MatchingLabels{
		"app.kubernetes.io/managed-by":       "gatus-operator",
		"gatus-operator.aumer.io/parent-uid": string(gatus.ObjectMeta.UID),
	})
	if err != nil {
		return err
	}

	if gatus.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(gatus, FINALIZER_NAME) {
			controllerutil.AddFinalizer(gatus, FINALIZER_NAME)
			r.Update(ctx, gatus)
		}
	} else {
		hasFinalizer := controllerutil.ContainsFinalizer(gatus, FINALIZER_NAME)

		if len(configMapList.Items) == 1 {
			configMap := configMapList.Items[0]
			err := r.deleteConfigMap(ctx, configMap)
			if err != nil {
				return err
			}

			if hasFinalizer {
				controllerutil.RemoveFinalizer(gatus, FINALIZER_NAME)
				if err := r.Update(ctx, gatus); err != nil {
					return err
				}
			}
		}

		if hasFinalizer {
			controllerutil.RemoveFinalizer(gatus, "gatus-operator.aumer.io/finalizer")
			if err := r.Update(ctx, gatus); err != nil {
				return err
			}
		}
		return nil
	}

	if len(configMapList.Items) == 0 {
		return r.createConfigMap(ctx, gatus)
	}

	if len(configMapList.Items) == 1 {
		configMap := configMapList.Items[0]
		return r.updateConfigMap(ctx, gatus, configMap)
	}

	return nil
}

func (r *ReconcileGatus) createConfigMap(ctx context.Context, gatus *gatusiov1alpha1.Gatus) error {
	yamlString, err := r.getEndpointsYaml(gatus)
	if err != nil {
		return fmt.Errorf("error getting YAML: %w", err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: r.generateMetaData(gatus),
		Data: map[string]string{
			"gatus.yaml": yamlString,
		},
	}

	r.Create(ctx, configMap)

	return nil
}

func (r *ReconcileGatus) updateConfigMap(ctx context.Context, gatus *gatusiov1alpha1.Gatus, configMap corev1.ConfigMap) error {
	yamlString, err := r.getEndpointsYaml(gatus)
	if err != nil {
		return fmt.Errorf("error getting YAML: %w", err)
	}

	configMap.Data["gatus.yaml"] = yamlString

	return r.Update(ctx, &configMap)
}

func (r *ReconcileGatus) deleteConfigMap(ctx context.Context, configMap corev1.ConfigMap) error {
	return r.Delete(ctx, &configMap)
}

func (r *ReconcileGatus) getEndpointsYaml(gatus *gatusiov1alpha1.Gatus) (string, error) {
	type GatusEndpoint struct {
		Endpoints []gatusiov1alpha1.EndpointEndpoint `json:"endpoints"`
	}

	gatusEndpoint := r.Config.ApplyDefaults(gatus.Spec.Endpoint)

	yamlBytes, err := yaml.Marshal(GatusEndpoint{Endpoints: []gatusiov1alpha1.EndpointEndpoint{gatusEndpoint}})
	if err != nil {
		return "", fmt.Errorf("error converting struct to YAML: %v", err)
	}

	yamlString := string(yamlBytes)

	return yamlString, nil
}

func (r *ReconcileGatus) generateMetaData(gatus *gatusiov1alpha1.Gatus) metav1.ObjectMeta {
	labels := map[string]string{
		"app.kubernetes.io/managed-by":       "gatus-operator",
		"gatus-operator.aumer.io/parent-uid": string(gatus.ObjectMeta.UID),
	}

	if r.Config.OperatorConfig.K8sSidecarAnnotation != "" {
		labels[r.Config.OperatorConfig.K8sSidecarAnnotation] = "enabled"
	}

	return metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", gatus.Name, "gatus-config"),
		Namespace: gatus.Namespace,
		Labels:    labels,
	}
}
