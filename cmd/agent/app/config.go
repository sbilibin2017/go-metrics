package app

type Config struct {
	Address        string
	PollInterval   int
	ReportInterval int
	Key            string
}

func (c *Config) GetAddress() string {
	return c.Address
}
