package logging

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ProviderLogger() *zap.Logger {

	encoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())

	file, _ := os.Create(fmt.Sprintf("binder_api_%d.log", time.Now().Unix()))

	fileLogger := zapcore.NewCore(encoder, zapcore.AddSync(file), zapcore.DebugLevel)
	consoleLogger := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	core := zapcore.NewTee(consoleLogger, fileLogger)
	logger := zap.New(core, zap.AddCaller())

	return logger
}
