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

	FlagAddress        = "address"
	FlagReportInterval = "report-interval"
	FlagPollInterval   = "poll-interval"

	ShortFlagAddress        = "a"
	ShortFlagReportInterval = "r"
	ShortFlagPollInterval   = "p"

	EnvAddress        = "ADDRESS"
	EnvReportInterval = "REPORT_INTERVAL"
	EnvPollInterval   = "POLL_INTERVAL"

	DescriptionAddress        = "Address of the HTTP server endpoint"
	DescriptionReportInterval = "Interval in seconds for sending metrics to the server"
	DescriptionPollInterval   = "Interval in seconds for polling metrics from the runtime package"
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

	viper.BindPFlag(EnvAddress, cmd.PersistentFlags().Lookup(FlagAddress))
	viper.BindPFlag(EnvReportInterval, cmd.PersistentFlags().Lookup(FlagReportInterval))
	viper.BindPFlag(EnvPollInterval, cmd.PersistentFlags().Lookup(FlagPollInterval))

	return cmd
}
