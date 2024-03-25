package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Driver     string `mapstructure:"DRIVER"`
	DataSource string `mapstructure:"DATA_SOURCE"`
	Key1       string `mapstructure:"AUTH_KEY1"`
	Key2       string `mapstructure:"AUTH_KEY2"`
}

func LoadConfig() (config Config, er error) {
	viper.SetConfigFile("./env/app.env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if er = viper.ReadInConfig(); er != nil {
		return
	}

	er = viper.Unmarshal(&config)

	return
}
