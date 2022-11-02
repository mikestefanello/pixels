package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	EnvLocal      Environment = "local"
	EnvProduction Environment = "prod"
)

type Environment string

type Config struct {
	App          string
	Environment  Environment
	Project      string
	Subscription string
	Topic        string
	HTTP         struct {
		Port    int
		Address string
	}
}

func GetConfig() (Config, error) {
	var c Config

	// Load the config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	// Load env variables
	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.BindEnv("project", "GCP_PROJECT"); err != nil {
		return c, nil
	}

	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}
