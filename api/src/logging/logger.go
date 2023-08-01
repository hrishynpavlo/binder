package logging

import (
	"go.uber.org/zap"
)

func ProviderLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment(zap.AddCaller())
	return logger
}
