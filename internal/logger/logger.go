package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func Init() error {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	zapLogger, err := config.Build()
	if err != nil {
		return err
	}
	Logger = zapLogger.Sugar()
	return nil
}
