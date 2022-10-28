package main

import (
	"log"
	"log-collect/config"
	"log-collect/internal/app/logTransfer"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	logTransfer.Run(cfg)
}
