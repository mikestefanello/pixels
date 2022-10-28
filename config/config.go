package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Environment  string
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
	// TODO: Per env
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}
