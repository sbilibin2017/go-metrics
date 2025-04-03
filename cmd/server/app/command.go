package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	c "go-metrics/pkg/context"
	"go-metrics/pkg/log"
)

const (
	DefaultAddress         = "localhost:8080"
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "data/metrics.json"
	DefaultRestore         = true

	FlagAddress         = "address"
	FlagStoreInterval   = "store-interval"
	FlagFileStoragePath = "file-storage-path"
	FlagRestore         = "restore"
	FlagDatabaseDSN     = "database-dsn"

	ShortFlagAddress         = "a"
	ShortFlagStoreInterval   = "i"
	ShortFlagFileStoragePath = "f"
	ShortFlagRestore         = "r"
	ShortFlagDatabaseDSN     = "d"

	EnvAddress         = "ADDRESS"
	EnvStoreInterval   = "STORE_INTERVAL"
	EnvFileStoragePath = "FILE_STORAGE_PATH"
	EnvRestore         = "RESTORE"
	EnvDatabaseDSN     = "DATABASE_DSN"

	DescriptionAddress         = "Address of the HTTP server endpoint"
	DescriptionStoreInterval   = "Interval in seconds to store metrics to disk"
	DescriptionFileStoragePath = "Path to the file to store metrics"
	DescriptionRestore         = "Whether to load previously saved values on server startup"
	DescriptionDatabaseDSN     = "Database DSN"
)

func NewCommand() *cobra.Command {
	viper.AutomaticEnv()

	cmd := &cobra.Command{
		Use:   "app",
		Short: "HTTP Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Init(log.LevelInfo)
			config := &Config{
				Address:         viper.GetString(EnvAddress),
				DatabaseDSN:     viper.GetString(EnvDatabaseDSN),
				StoreInterval:   viper.GetInt(EnvStoreInterval),
				FileStoragePath: viper.GetString(EnvFileStoragePath),
				Restore:         viper.GetBool(EnvRestore),
			}
			container, err := NewContainer(config)
			if err != nil {
				return err
			}
			worker := NewWorker(config, container)
			server := NewServer(config, container, worker)
			ctx, cancel := c.NewContext()
			defer cancel()
			return server.Start(ctx)
		},
	}

	cmd.PersistentFlags().StringP(FlagAddress, ShortFlagAddress, DefaultAddress, DescriptionAddress)
	cmd.PersistentFlags().IntP(FlagStoreInterval, ShortFlagStoreInterval, DefaultStoreInterval, DescriptionStoreInterval)
	cmd.PersistentFlags().StringP(FlagFileStoragePath, ShortFlagFileStoragePath, DefaultFileStoragePath, DescriptionFileStoragePath)
	cmd.PersistentFlags().BoolP(FlagRestore, ShortFlagRestore, DefaultRestore, DescriptionRestore)
	cmd.PersistentFlags().StringP(FlagDatabaseDSN, ShortFlagDatabaseDSN, "", DescriptionDatabaseDSN)

	viper.BindPFlag(EnvAddress, cmd.PersistentFlags().Lookup(FlagAddress))
	viper.BindPFlag(EnvStoreInterval, cmd.PersistentFlags().Lookup(FlagStoreInterval))
	viper.BindPFlag(EnvFileStoragePath, cmd.PersistentFlags().Lookup(FlagFileStoragePath))
	viper.BindPFlag(EnvRestore, cmd.PersistentFlags().Lookup(FlagRestore))
	viper.BindPFlag(EnvDatabaseDSN, cmd.PersistentFlags().Lookup(FlagDatabaseDSN))

	return cmd
}
