package app

type Config struct {
	Address        string
	PollInterval   int
	ReportInterval int
}

func (c *Config) GetAddress() string {
	return c.Address
}
