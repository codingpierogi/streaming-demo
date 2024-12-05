package config

import "github.com/spf13/viper"

type Config struct {
	Mode             string `mapstructure:"MODE"`
	Addr             string `mapstructure:"ADDR"`
	BootstrapServers string `mapstructure:"BOOTSTRAP_SERVERS"`
	Topic            string `mapstructure:"TOPIC"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("notify")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
