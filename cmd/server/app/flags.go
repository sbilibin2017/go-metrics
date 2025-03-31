package app

import (
	"flag"
	"go-metrics/internal/configs"
	"os"
)

func ParseFlags() *configs.ServerConfig {
	serverAddr := "localhost:8080"
	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		serverAddr = envAddr
	}
	flag.StringVar(&serverAddr, "a", serverAddr, "Адрес эндпоинта HTTP-сервера")
	flag.Parse()
	return &configs.ServerConfig{
		Address: serverAddr,
	}
}
