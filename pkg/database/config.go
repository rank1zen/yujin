package database

type Config struct {
        Url string
}

func (c *Config) DSN() string {
        return c.Url
}

func NewConfig(connString string) *Config {
        return &Config{Url: connString}
}
