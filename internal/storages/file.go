package storages

import (
	"go-metrics/internal/logger"
	"os"
	"path/filepath"
)

type FileStorageConfig interface {
	GetFileStoragePath() string
}

func NewFile(config FileStorageConfig) (*os.File, error) {
	logger.Logger.Infow("Opening file storage", "path", config.GetFileStoragePath())
	dir := filepath.Dir(config.GetFileStoragePath())
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Logger.Errorw("Failed to create directories", "error", err)
		return nil, err
	}
	file, err := os.OpenFile(config.GetFileStoragePath(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		logger.Logger.Errorw("Failed to open file", "error", err)
		return nil, err
	}
	return file, nil
}
