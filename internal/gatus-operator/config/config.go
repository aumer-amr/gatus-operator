package config

import (
	"os"
)

type Config struct {
	MetricsAddr string
	ProbeAddr   string
	LogLevel    string
	ConfigPath  string
	DevMode     bool
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
	}
	return config
}
