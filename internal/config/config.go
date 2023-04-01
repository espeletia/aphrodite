package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig  ServerConfig
	ServiceConfig ServiceConfig
}

func LoadConfig() *Config {
	config := &Config{
		ServerConfig:  loadServerConfig(),
		ServiceConfig: loadServiceConfig(),
	}
	return config
}

func configViper(configName string) *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configurations/")
	return v
}
