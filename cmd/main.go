package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "finance_project/docs"
	"finance_project/internal/app"
	"finance_project/internal/config"
	"finance_project/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Загружаем .env файл
	err := godotenv.Load("/home/abylay/finance_project/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Запуск миграций
	migrationsDir := "./migrations"
	appliedMigrations, err := database.RunMigrations(db, migrationsDir)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// Вывод выполненных миграций
	if len(appliedMigrations) > 0 {
		fmt.Println("Applied migrations:")
		for _, migration := range appliedMigrations {
			fmt.Println("-", migration)
		}
	} else {
		fmt.Println("No new migrations were applied.")
	}

	// Запуск приложения
	app.Run(cfg)
}
