package config

import (
	"fmt"

	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/spf13/viper"
)

type Config struct {
	SMTPHost          string `env:"SMTPHOST"`
	SMTPPort          string `env:"SMTPPORT"`
	AuthEmail         string `env:"AUTHEMAIL"`
	AuthEmailPassword string `env:"AUTHEMAILPASSWORD"`
	Port              string `env:"PORT"`
	WithJWTAuth       bool   `env:"WITHJWTAUTH"`
	DbString          string `env:"DBSTRING"`
}

func ReadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile("./configs/app.env")
	err := v.ReadInConfig()
	errors.FatalOnError(err)

	config := &Config{
		SMTPHost:          v.GetString("SMTPHOST"),
		SMTPPort:          v.GetString("SMTPPORT"),
		AuthEmail:         v.GetString("AUTHEMAIL"),
		AuthEmailPassword: v.GetString("AUTHEMAILPASSWORD"),
		Port:              v.GetString("PORT"),
		WithJWTAuth:       v.GetBool("WITHJWTAUTH"),
		DbString:          v.GetString("DBSTRING"),
	}
	fmt.Println(config)

	return config, nil
}
