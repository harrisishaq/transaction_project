package config

import (
	"log"

	"github.com/spf13/viper"
)

var configStruct = map[string]interface{}{
	"gorm-config": &GormConfig,
	"app-config":  &AppConfig,
}

type EnvAppConfig struct {
	DefaultPassword string `mapstructure:"DEFAULT_PASSWORD" validate:"required"`
}

var AppConfig EnvAppConfig

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	for key, value := range configStruct {
		log.Println("Loading config: ", key)
		if err := viper.Unmarshal(value); err != nil {
			log.Printf("Error loading config %s, cause: %+v\n", key, err)
			log.Fatal(err)
		}
		log.Printf("%s: %+v\n", key, value)
		log.Printf("Config %s loaded successfully", key)
	}
}
