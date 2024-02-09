package postgresql

type Config struct {
	Name     string
	User     string
	Host     string
	Port     string
	SSLMode  string
	Password string
}

func (c *Config) ConnectionURL() string {
	if c == nil {
		return ""
	}

	return ""
}
