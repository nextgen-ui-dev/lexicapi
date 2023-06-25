package app

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env  string `mapstructure:"ENV"`
	Port string `mapstructure:"PORT"`
}

func SetConfig() (config Config, err error) {
	viper.AddConfigPath("../")
	viper.SetConfigName(".env")
	viper.SetEnvPrefix("")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
