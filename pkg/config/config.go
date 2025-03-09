package config

type Config struct {
	SMTPHost          string
	SMTPPort          string
	AuthEmail         string
	AuthEmailPassword string
}

func (c Config) ReadConfig() error {
	panic("not implemented")
}
