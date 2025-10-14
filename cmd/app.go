package main

import (
	"log"

	"github.com/alpinesboltltd/boltz-ai/internal/app"
	"github.com/alpinesboltltd/boltz-ai/internal/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	godotenv.Load(".env")
	var cfg config.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Web interface
	app.Run(&cfg)
}
