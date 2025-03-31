package app

import (
	"flag"
	"go-metrics/internal/configs"
)

func ParseFlags() *configs.AgentConfig {
	addr := flag.String("a", "localhost:8080", "Адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)")
	reportInterval := flag.Int("r", 10, "Частота отправки метрик на сервер в секундах (по умолчанию 10 секунд)")
	pollInterval := flag.Int("p", 2, "Частота опроса метрик из пакета runtime в секундах (по умолчанию 2 секунды)")

	flag.Parse()

	return &configs.AgentConfig{
		Address:        *addr,
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
	}
}
