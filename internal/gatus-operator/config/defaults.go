package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/aumer-amr/gatus-operator/v2/api/v1alpha1"
	"gopkg.in/yaml.v3"
)

func getDefaults() (*v1alpha1.EndpointEndpoint, error) {
	configuration := Generate()
	filePath := configuration.ConfigPath

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var defaults = v1alpha1.EndpointEndpoint{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&defaults); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %v", err)
	}

	return &defaults, nil
}

func HasDefaults() bool {
	_, err := getDefaults()
	return err == nil
}

func ApplyDefaults(override v1alpha1.EndpointEndpoint) v1alpha1.EndpointEndpoint {
	base, err := getDefaults()
	if err != nil {
		return override
	}

	baseVal := reflect.ValueOf(base).Elem()
	overrideVal := reflect.ValueOf(override)

	for i := 0; i < baseVal.NumField(); i++ {
		baseField := baseVal.Field(i)
		overrideField := overrideVal.Field(i)

		if !overrideField.IsZero() {
			baseField.Set(overrideField)
		}
	}

	return *base
}
