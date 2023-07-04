package app

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Env  string `mapstructure:"ENV"`
	Port string `mapstructure:"PORT"`

	OpenAIOrganizationId string `mapstructure:"OPENAI_ORGANIZATION_ID"`
	OpenAIAPIKey         string `mapstructure:"OPENAI_API_KEY"`

	DbUrl  string `mapstructure:"DB_URL"`
	DbHost string `mapstructure:"DB_HOST"`
	DbPort string `mapstructure:"DB_PORT"`
	DbName string `mapstructure:"DB_NAME"`
	DbUser string `mapstructure:"DB_USER"`
	DbPwd  string `mapstructure:"DB_PWD"`
	DbDsn  string
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	config.setDBConfig()

	return
}

func (c *Config) setDBConfig() {
	var ssl string
	if c.Env == "local" {
		ssl = "sslmode=disable"
	} else {
		//TODO: add ssl cert for dev and prod
		ssl = "sslmode=require"
	}

	c.DbDsn = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s %s", c.DbHost, c.DbPort, c.DbName, c.DbUser, c.DbPwd, ssl)
}
