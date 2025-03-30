package configs

type AgentConfig struct {
	Address        string
	PollInterval   int
	ReportInterval int
}

func (c *AgentConfig) GetAddress() string {
	return c.Address
}
