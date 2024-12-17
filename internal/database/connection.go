package database

import (
	"database/sql"
	"finance_project/internal/config"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName, // Исправлено: заменено cfg.Name на cfg.DBName
		cfg.SSLMode,
	)
	return sql.Open("postgres", connStr)
}
