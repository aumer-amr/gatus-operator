package config

import (
	"fmt"
	"os"

	"github.com/aumer-amr/gatus-operator/v2/api/v1alpha1"
	"gopkg.in/yaml.v3"
)

type DefaultsConfig struct {
	Defaults map[string]v1alpha1.EndpointEndpoint `yaml:"defaults,omitempty"`
}

type OperatorConfig struct {
	K8sSidecarAnnotation string `yaml:"k8s-sidecar-annotation,omitempty"`
}

type Config struct {
	MetricsAddr    string
	ProbeAddr      string
	LogLevel       string
	ConfigPath     string
	DevMode        bool
	DefaultsConfig DefaultsConfig `yaml:"defaults,omitempty"`
	OperatorConfig OperatorConfig `yaml:"operator,omitempty"`
}

func Generate() *Config {
	config := &Config{
		MetricsAddr: func() string {
			if addr := os.Getenv("METRICS_ADDR"); addr != "" {
				return addr
			}
			return ":8080"
		}(),
		ProbeAddr: func() string {
			if addr := os.Getenv("PROBE_ADDR"); addr != "" {
				return addr
			}
			return ":8081"
		}(),
		LogLevel: func() string {
			if level := os.Getenv("LOG_LEVEL"); level != "" {
				return level
			}
			return "info"
		}(),
		ConfigPath: func() string {
			if path := os.Getenv("CONFIG_PATH"); path != "" {
				return path
			}
			return "/config/config.yaml"
		}(),
		DevMode: func() bool {
			if devMode := os.Getenv("DEV_MODE"); devMode != "" {
				return devMode == "true"
			}
			return false
		}(),
		DefaultsConfig: DefaultsConfig{},
		OperatorConfig: OperatorConfig{},
	}

	if config.ConfigPath != "" {
		readConfig, err := readConfigFile(config.ConfigPath)
		if err != nil {
			panic(fmt.Errorf("failed to read config file: %v", err))
		}

		config.DefaultsConfig = readConfig.DefaultsConfig
		config.OperatorConfig = readConfig.OperatorConfig
	}

	return config
}

func readConfigFile(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config = Config{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %v", err)
	}

	return &config, nil
}
