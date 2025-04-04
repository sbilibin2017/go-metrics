package app

import (
	"go-metrics/pkg/context"
	"go-metrics/pkg/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultAddress        = "localhost:8080"
	DefaultReportInterval = 10
	DefaultPollInterval   = 2
	DefaultRateLimit      = 1

	FlagAddress        = "address"
	FlagReportInterval = "report-interval"
	FlagPollInterval   = "poll-interval"
	FlagKey            = "key"
	FlagRateLimit      = "rate-limit"

	ShortFlagAddress        = "a"
	ShortFlagReportInterval = "r"
	ShortFlagPollInterval   = "p"
	ShortFlagKey            = "k"
	ShortFlagRateLimit      = "l"

	EnvAddress        = "ADDRESS"
	EnvReportInterval = "REPORT_INTERVAL"
	EnvPollInterval   = "POLL_INTERVAL"
	EnvKey            = "KEY"
	EnvRateLimit      = "RATE_LIMIT"

	DescriptionAddress        = "Address of the HTTP server endpoint"
	DescriptionReportInterval = "Interval in seconds for sending metrics to the server"
	DescriptionPollInterval   = "Interval in seconds for polling metrics from the runtime package"
	DescriptionKey            = "Secret key for data signing"
	DescriptionRateLimit      = "Limit the number of concurrent outgoing requests"
)

func NewCommand() *cobra.Command {
	viper.AutomaticEnv()

	cmd := &cobra.Command{
		Use:   "metrics-agent",
		Short: "Metrics Agent for collecting and sending metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Init(log.LevelInfo)
			config := &Config{
				Address:        viper.GetString(EnvAddress),
				ReportInterval: viper.GetInt(EnvReportInterval),
				PollInterval:   viper.GetInt(EnvPollInterval),
				Key:            viper.GetString(EnvKey),
				RateLimit:      viper.GetInt(EnvRateLimit),
			}
			agent := NewMetricAgent(config)
			ctx, cancel := context.NewContext()
			defer cancel()
			return agent.Start(ctx)
		},
	}

	cmd.PersistentFlags().StringP(FlagAddress, ShortFlagAddress, DefaultAddress, DescriptionAddress)
	cmd.PersistentFlags().IntP(FlagReportInterval, ShortFlagReportInterval, DefaultReportInterval, DescriptionReportInterval)
	cmd.PersistentFlags().IntP(FlagPollInterval, ShortFlagPollInterval, DefaultPollInterval, DescriptionPollInterval)
	cmd.PersistentFlags().StringP(FlagKey, ShortFlagKey, "", DescriptionKey)
	cmd.PersistentFlags().IntP(FlagRateLimit, ShortFlagRateLimit, DefaultRateLimit, DescriptionRateLimit)

	viper.BindPFlag(EnvAddress, cmd.PersistentFlags().Lookup(FlagAddress))
	viper.BindPFlag(EnvReportInterval, cmd.PersistentFlags().Lookup(FlagReportInterval))
	viper.BindPFlag(EnvPollInterval, cmd.PersistentFlags().Lookup(FlagPollInterval))
	viper.BindPFlag(EnvKey, cmd.PersistentFlags().Lookup(FlagKey))
	viper.BindPFlag(EnvRateLimit, cmd.PersistentFlags().Lookup(FlagRateLimit))

	return cmd
}
