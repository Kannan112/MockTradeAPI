package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
}

// name of envs and used to read from system envs
var envsNames = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD",
}

func LoadConfig() (config Config, err error) {

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()

	if err != nil {
		for _, env := range envsNames {
			if err := viper.BindEnv(env); err != nil {
				return config, err
			}
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(config); err != nil {
		return config, err
	}
	return config, nil
}
