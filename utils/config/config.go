package config

type Config struct {
	SMTPHost          string
	SMTPPort          string
	AuthEmail         string
	AuthEmailPassword string
	Port              string
	WithJWTAuth       bool
}

func ReadConfig() (*Config, error) {
	panic("not implemented")
}
