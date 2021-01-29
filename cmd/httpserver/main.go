package main

import (
	"HttpServerAvito/internal/app/httpserver"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configPath string
)

func InitConfig() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to server conf")
}

func main() {
	InitConfig()
	flag.Parse()

	config := httpserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	err = httpserver.Start(config)
	if err != nil {
		fmt.Println("fail start server")
	}
}
