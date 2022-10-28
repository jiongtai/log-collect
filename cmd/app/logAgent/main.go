package main

import (
	"log"
	"log-collect/config"
	"log-collect/internal/app/logAgent"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	logAgent.Run(cfg)
}
