package store

type Config struct {
	DataBaseUrl string `toml:"database_url"`
}

func NewConfig() *Config{
	return &Config{
		DataBaseUrl: "host=localhost dbname=avito_httpserv sslmode=disable",
	}
}
