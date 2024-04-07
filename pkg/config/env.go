package config

import (
	"log"

	"github.com/spf13/viper"
)

type envConfigs struct {
	DataPath string `mapstructure:"DATA_PATH"`
}

var EnvConfigs *envConfigs

func init() {
	EnvConfigs = loadEnvs()
}

func loadEnvs() (config *envConfigs) {
	viper.AddConfigPath(".")

	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return config
}
