package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	IntervalMs int `mapstructure:"INTERVAL_MS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("sysinfo")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}