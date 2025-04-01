package app

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	FlagAddress         = "a"
	FlagStoreInterval   = "i"
	FlagFileStoragePath = "f"
	FlagRestore         = "r"
	FlagDatabaseDSN     = "d"

	DescriptionAddress         = "Адрес эндпоинта HTTP-сервера"
	DescriptionStoreInterval   = "Интервал времени в секундах для сохранения показаний на диск"
	DescriptionFileStoragePath = "Путь до файла для хранения показаний"
	DescriptionRestore         = "Загружать ли ранее сохранённые значения при старте сервера"
	DescriptionDatabaseDSN     = "Строка подключения к базе данных"

	DefaultAddress         = "localhost:8080"
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "data/metrics.json"
	DefaultRestore         = true
	DefaultDatabaseDSN     = ""
)

func ParseFlags() *Config {
	serverAddr := DefaultAddress
	storeInterval := DefaultStoreInterval
	fileStoragePath := DefaultFileStoragePath
	restore := DefaultRestore
	databaseDSN := DefaultDatabaseDSN

	flag.StringVar(&serverAddr, FlagAddress, serverAddr, DescriptionAddress)
	flag.IntVar(&storeInterval, FlagStoreInterval, storeInterval, DescriptionStoreInterval)
	flag.StringVar(&fileStoragePath, FlagFileStoragePath, fileStoragePath, DescriptionFileStoragePath)
	flag.BoolVar(&restore, FlagRestore, restore, DescriptionRestore)
	flag.StringVar(&databaseDSN, FlagDatabaseDSN, databaseDSN, DescriptionDatabaseDSN)

	flag.Parse()

	if envVal, exists := os.LookupEnv("ADDRESS"); exists {
		serverAddr = envVal
	}
	if envVal, exists := os.LookupEnv("STORE_INTERVAL"); exists {
		if val, err := strconv.Atoi(envVal); err == nil {
			storeInterval = val
		} else {
			fmt.Fprintf(os.Stderr, "Ошибка парсинга STORE_INTERVAL: %v\n", err)
		}
	}
	if envVal, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists {
		fileStoragePath = envVal
	}
	if envVal, exists := os.LookupEnv("RESTORE"); exists {
		if val, err := strconv.ParseBool(envVal); err == nil {
			restore = val
		} else {
			fmt.Fprintf(os.Stderr, "Ошибка парсинга RESTORE: %v\n", err)
		}
	}
	if envVal, exists := os.LookupEnv("DATABASE_DSN"); exists {
		databaseDSN = envVal
	}

	return &Config{
		Address:         serverAddr,
		StoreInterval:   storeInterval,
		FileStoragePath: fileStoragePath,
		Restore:         restore,
		DatabaseDSN:     databaseDSN,
	}
}
