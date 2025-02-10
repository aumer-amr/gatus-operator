package config

import (
	"encoding/json"

	"github.com/aumer-amr/gatus-operator/v2/api/v1alpha1"
)

func (c Config) HasDefaults() bool {
	defaults := c.DefaultsConfig.Defaults
	return defaults == nil
}

func (c Config) ApplyDefaults(override v1alpha1.EndpointEndpoint) v1alpha1.EndpointEndpoint {
	base := c.DefaultsConfig

	overrideMap := structToMap(override)
	normalizeEmptyObjects(overrideMap)

	mergedMap := make(map[string]interface{})
	if globalDefaults, hasGlobal := base.Defaults["global"]; hasGlobal {
		mergedMap = structToMap(globalDefaults)
	}

	if group, hasGroup := base.Defaults[override.Group]; hasGroup {
		groupMap := structToMap(group)
		mergedMap = mergeMaps(mergedMap, groupMap)
	}

	mergedMap = mergeMaps(mergedMap, overrideMap)

	var result v1alpha1.EndpointEndpoint
	mapToStruct(mergedMap, &result)

	return result
}

func structToMap(input interface{}) map[string]interface{} {
	data, _ := json.Marshal(input)
	var result map[string]interface{}
	_ = json.Unmarshal(data, &result)
	return result
}

func mapToStruct(input map[string]interface{}, output interface{}) {
	data, _ := json.Marshal(input)
	_ = json.Unmarshal(data, output)
}

func mergeMaps(base, override map[string]interface{}) map[string]interface{} {
	for k, v := range override {
		base[k] = v
	}
	return base
}

func normalizeEmptyObjects(data map[string]interface{}) {
	for k, v := range data {
		if subMap, ok := v.(map[string]interface{}); ok && len(subMap) == 0 {
			data[k] = nil
		}
	}
}
