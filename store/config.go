package store

type Config struct {
	DataBaseUrl string `toml:"database_url"`
}

func NewConfig() *Config{
	return &Config{
		DataBaseUrl: "host=database port=5432 user=postgres dbname=avito_httpserv password=password sslmode=disable",
	}
}
