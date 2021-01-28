package httpserver

import "HttpServerAvito/store"

type Config struct {
	BindAddr 	string `toml:"bind_addr"`
	StoreConfig *store.Config ``
}

func NewConfig () *Config  {
	return &Config{
		BindAddr: ":8080",
		StoreConfig: store.NewConfig(),
	}
}