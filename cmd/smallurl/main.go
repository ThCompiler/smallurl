package main

import (
	"flag"
	"log"
	"smallurl/internal/app"
	"smallurl/internal/app/config"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/localhost-config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
