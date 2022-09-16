package log

import (
	"bus-api/core/config"

	"go.uber.org/zap"
)

func ConfigLogger() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		config.ReadConfig("log.output"),
	}
	level := config.ReadConfig("log.level")
	logLevel := zap.DebugLevel
	switch level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	default:
		logLevel = zap.InfoLevel
	}
	cfg.Level = zap.NewAtomicLevelAt(logLevel)
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
	// logger.Info("logger construction succeeded")
}

func Debug(message string) {
	zap.L().Debug(message)
}

func Info(message string) {
	zap.L().Info(message)
}

func Error(message string) {
	zap.L().Error(message)
}

func Fatal(message string) {
	zap.L().Fatal(message)
}
