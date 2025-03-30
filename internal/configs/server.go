package configs

type ServerConfig struct {
	Address         string
	DatabaseDSN     string
	StoreInterval   string
	FileStoragePath string
	Restore         string
}
