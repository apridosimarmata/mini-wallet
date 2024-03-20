package infrastructure

type Config struct {
	POSTGRES_DB       string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
}

func GetConfig() Config {
	return Config{}
}
