package infrastructure

import "github.com/ory/viper"

type Config struct {
	POSTGRES_DB       string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string

	REDIS_HOST string
	REDIS_PORT string
}

func GetConfig() Config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return Config{
		POSTGRES_DB:       viper.GetString("POSTGRES_DB"),
		POSTGRES_HOST:     viper.GetString("POSTGRES_HOST"),
		POSTGRES_PORT:     viper.GetString("POSTGRES_PORT"),
		POSTGRES_USER:     viper.GetString("POSTGRES_USER"),
		POSTGRES_PASSWORD: viper.GetString("POSTGRES_PASSWORD"),
		REDIS_HOST:        viper.GetString("REDIS_HOST"),
		REDIS_PORT:        viper.GetString("REDIS_PORT"),
	}
}
