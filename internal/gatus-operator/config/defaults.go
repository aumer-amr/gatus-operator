package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

	jsonOverride, err := json.Marshal(override)
	if err != nil {
		return override
	}

	err = json.NewDecoder(strings.NewReader(string(jsonOverride))).Decode(&base)
	if err != nil {
		panic(err)
	}

	return *base
}
