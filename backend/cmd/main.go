package main

import (
	"log"

	_ "finance_project/docs"
	"finance_project/internal/app"
	"finance_project/internal/config"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("/home/abylay/finance_project/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app.Run(cfg)

}
