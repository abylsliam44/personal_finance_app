package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations выполняет все SQL-миграции из папки migrations.
func RunMigrations(db *sql.DB, migrationsDir string) error {
	// Чтение списка файлов в директории миграций.
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Фильтрация только SQL-файлов.
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Проверка, есть ли миграции.
	if len(migrationFiles) == 0 {
		log.Println("No migration files found.")
		return nil
	}

	// Сортировка файлов по имени.
	sort.Strings(migrationFiles)

	// Проверка, существует ли таблица migrations
	if err := ensureMigrationsTableExists(db); err != nil {
		return err
	}

	// Применение каждой миграции.
	for _, file := range migrationFiles {
		// Проверка, была ли уже применена эта миграция
		if err := applyMigrationIfNotApplied(db, file, migrationsDir); err != nil {
			return err
		}
	}

	log.Println("All migrations applied successfully.")
	return nil
}

// ensureMigrationsTableExists проверяет существование таблицы для хранения миграций.
func ensureMigrationsTableExists(db *sql.DB) error {
	var exists bool
	// Проверяем, существует ли таблица migrations
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'migrations')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check migrations table existence: %w", err)
	}

	// Если таблица не существует, возвращаем ошибку
	if !exists {
		return fmt.Errorf("migrations table does not exist")
	}
	return nil
}

// applyMigrationIfNotApplied проверяет, была ли уже применена миграция, и если нет - применяет ее.
func applyMigrationIfNotApplied(db *sql.DB, file string, migrationsDir string) error {
	// Проверяем, была ли уже применена эта миграция
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE migration_name = $1", file).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if migration has been applied: %w", err)
	}

	// Если миграция уже применена, пропускаем ее
	if count > 0 {
		log.Printf("Migration %s already applied, skipping.", file)
		return nil
	}

	// Чтение содержимого миграции
	path := filepath.Join(migrationsDir, file)
	query, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", file, err)
	}

	// Применение миграции
	log.Printf("Applying migration: %s", file)
	_, err = db.Exec(string(query))
	if err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", file, err)
	}

	// Запись миграции в таблицу
	_, err = db.Exec("INSERT INTO migrations (migration_name) VALUES ($1)", file)
	if err != nil {
		return fmt.Errorf("failed to record migration in database: %w", err)
	}

	log.Printf("Migration %s applied successfully", file)
	return nil
}
