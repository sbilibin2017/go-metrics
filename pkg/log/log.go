package log

import (
	"sync"

	"go.uber.org/zap"
)

type Level string

const (
	LevelInfo  Level = "info"
	LevelError Level = "error"
)

var logger *zap.SugaredLogger
var mu sync.Mutex

func Init(level Level) error {
	mu.Lock()
	defer mu.Unlock()
	zapConfig := zap.NewProductionConfig()
	switch level {
	case LevelInfo:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case LevelError:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logInstance, err := zapConfig.Build()
	if err != nil {
		return err
	}
	logInstance = logInstance.WithOptions(zap.WithCaller(false))
	logger = logInstance.Sugar()
	return nil
}

func Info(msg string, args ...any) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.Infow(msg, args...)
	}
}

func Error(msg string, args ...any) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.Errorw(msg, args...)
	}
}

func Sync() {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.Sync()
	}
}
