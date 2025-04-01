package app

import (
	"flag"
	"os"
	"strconv"
)

func ParseFlags() *Config {
	addr := "localhost:8080"
	reportInterval := 10
	pollInterval := 2
	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		addr = envAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if ri, err := strconv.Atoi(envReportInterval); err == nil {
			reportInterval = ri
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if pi, err := strconv.Atoi(envPollInterval); err == nil {
			pollInterval = pi
		}
	}
	flag.StringVar(&addr, "a", addr, "Адрес эндпоинта HTTP-сервера")
	flag.IntVar(&reportInterval, "r", reportInterval, "Частота отправки метрик на сервер в секундах")
	flag.IntVar(&pollInterval, "p", pollInterval, "Частота опроса метрик из пакета runtime в секундах")
	flag.Parse()
	return &Config{
		Address:        addr,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
	}
}
