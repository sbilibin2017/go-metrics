package app

import (
	"flag"
	"os"
	"strconv"
)

const (
	FlagAddress         = "a"
	FlagStoreInterval   = "i"
	FlagFileStoragePath = "f"
	FlagRestore         = "r"

	DescriptionAddress         = "Адрес эндпоинта HTTP-сервера"
	DescriptionStoreInterval   = "Интервал времени в секундах для сохранения показаний на диск"
	DescriptionFileStoragePath = "Путь до файла для хранения показаний"
	DescriptionRestore         = "Загружать ли ранее сохранённые значения при старте сервера"

	DefaultAddress         = "localhost:8080"
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "data/metrics.json"
	DefaultRestore         = false
)

func ParseFlags() *Config {
	serverAddr := DefaultAddress
	storeInterval := DefaultStoreInterval
	fileStoragePath := DefaultFileStoragePath
	restore := DefaultRestore
	databaseDSN := ""

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		serverAddr = envAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if val, err := strconv.Atoi(envStoreInterval); err == nil {
			storeInterval = val
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if val, err := strconv.ParseBool(envRestore); err == nil {
			restore = val
		}
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		databaseDSN = envDatabaseDSN
	}

	flag.StringVar(&serverAddr, FlagAddress, serverAddr, DescriptionAddress)
	flag.IntVar(&storeInterval, FlagStoreInterval, storeInterval, DescriptionStoreInterval)
	flag.StringVar(&fileStoragePath, FlagFileStoragePath, fileStoragePath, DescriptionFileStoragePath)
	flag.BoolVar(&restore, FlagRestore, restore, DescriptionRestore)
	flag.StringVar(&databaseDSN, "d", databaseDSN, "database dsn")
	flag.Parse()

	return &Config{
		Address:         serverAddr,
		DatabaseDSN:     databaseDSN,
		StoreInterval:   storeInterval,
		FileStoragePath: fileStoragePath,
		Restore:         restore,
	}
}
