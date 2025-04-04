package app

type Config struct {
	Address         string
	DatabaseDSN     string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	Key             string
}

func (c *Config) GetAddress() string {
	return c.Address
}

func (c *Config) GetFileStoragePath() string {
	return c.FileStoragePath
}

func (c *Config) GetDatabaseDSN() string {
	return c.DatabaseDSN
}

func (c *Config) GetKey() string {
	return c.Key
}
