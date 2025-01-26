package main

import (
	"github.com/aumer-amr/gatus-operator/v2/internal/gatus-operator/config"
	"github.com/aumer-amr/gatus-operator/v2/internal/gatus-operator/controller"
	"github.com/aumer-amr/gatus-operator/v2/internal/gatus-operator/manager"
	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	config := config.Generate()

	logLevel, err := zapcore.ParseLevel(config.LogLevel)
	if err != nil {
		logLevel = zapcore.InfoLevel
	}

	ctrl.SetLogger(zap.New(zap.Level(logLevel), zap.UseDevMode(true)))

	mgr := manager.Run()
	controller.Run(mgr)
}
