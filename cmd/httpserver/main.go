package main

import (
	"HttpServerAvito/internal/app/httpserver"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"time"
)

var (
	configPath string
)

func InitConfig() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to server conf")
}

func main() {
	fmt.Println("точка входа")
	time.Sleep(time.Second * 5)
	InitConfig()
	flag.Parse()

	config := httpserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("перед стартом")
	err = httpserver.Start(config)
	if err != nil {
		fmt.Println("fail start server")
	}
}
