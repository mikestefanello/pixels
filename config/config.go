package config

import (
	"fmt"

	"github.com/spf13/pflag"
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

	// Determine the environment
	// You could determine this based on the GCP project env variable
	pflag.String("environment", "local", "the execution environment")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return c, err
	}

	// Load the config file per environment
	viper.SetConfigName(fmt.Sprintf("config-%s", viper.Get("Environment")))
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
