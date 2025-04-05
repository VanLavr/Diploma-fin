package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/VanLavr/Diploma-fin/utils/errors"
)

type Config struct {
	SMTPHost          string `env:"SMTPHOST"`
	SMTPPort          string `env:"SMTPPORT"`
	AuthEmail         string `env:"AUTHEMAIL"`
	AuthEmailPassword string `env:"AUTHEMAILPASSWORD"`
	Port              string `env:"PORT"`
	WithJWTAuth       bool   `env:"WITHJWTAUTH"`
	DbString          string `env:"DBSTRING"`
	SMTP2OAuthCode    string `env:"SMTP2OAUTHCODE"`
	Secret            string `env:"SECRET"`
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
		SMTP2OAuthCode:    v.GetString("SMTP2OAUTHCODE"),
		Secret:            v.GetString("SECRET"),
	}
	fmt.Println(config)

	return config, nil
}
