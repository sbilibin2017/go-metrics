package app

import (
	"flag"
	"go-metrics/internal/configs"
)

func ParseFlags() *configs.ServerConfig {
	serverAddr := flag.String("a", "localhost:8080", "Адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)")
	flag.Parse()
	return &configs.ServerConfig{
		Address: *serverAddr,
	}
}
