package config

import "github.com/spf13/viper"

type Config struct {
	REDIS_ADDR     string `mapstructure:"REDIS_ADDR"`
	REDIS_PASSWORD string `mapstructure:"REDIS_PASSWORD"`
	REDIS_DB       int    `mapstructure:"REDIS_DB"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
