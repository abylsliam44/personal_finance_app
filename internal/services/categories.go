package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	"finance_project/internal/models"

	"github.com/go-redis/redis/v8"
)

type CategoryService struct {
	DB          *sql.DB
	RedisClient *redis.Client // Убедитесь, что это поле существует
}

// NewCategoryService создает новый сервис для работы с категориями.
func NewCategoryService(db *sql.DB) *CategoryService {
	return &CategoryService{DB: db}
}

// GetAllCategories возвращает все категории.
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	query := `SELECT id, user_id, name, type, created_at FROM categories`
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Printf("Error retrieving categories: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.CreatedAt)
		if err != nil {
			log.Printf("Error scanning category row: %v", err)
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

// CreateCategory добавляет новую категорию.
func (s *CategoryService) CreateCategory(category models.Category) error {
	query := `INSERT INTO categories (user_id, name, type, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := s.DB.Exec(query, category.UserID, category.Name, category.Type)
	if err != nil {
		log.Printf("Error creating category: %v", err)
		return err
	}
	return nil
}

// GetCategoryByID возвращает категорию по ID.
func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	query := `SELECT id, user_id, name, type, created_at FROM categories WHERE id = $1`
	row := s.DB.QueryRow(query, id)

	var c models.Category
	err := row.Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Category with ID %d not found", id)
			return nil, nil
		}
		log.Printf("Error retrieving category: %v", err)
		return nil, err
	}
	return &c, nil
}

// UpdateCategory обновляет категорию.
func (s *CategoryService) UpdateCategory(category models.Category) error {
	query := `UPDATE categories SET name = $1, type = $2 WHERE id = $3`
	_, err := s.DB.Exec(query, category.Name, category.Type, category.ID)
	if err != nil {
		log.Printf("Error updating category: %v", err)
		return err
	}
	return nil
}

// DeleteCategory удаляет категорию по ID.
func (s *CategoryService) DeleteCategory(id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		return err
	}
	return nil
}

// GetTransactionsByCategory retrieves transactions by category, with optional caching
func (s *CategoryService) GetTransactionsByCategory(categoryID int) ([]models.Transaction, error) {
	ctx := context.Background()
	cacheKey := "transactions:category:" + strconv.Itoa(categoryID)

	// Check Redis cache
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var transactions []models.Transaction
		if json.Unmarshal([]byte(cachedData), &transactions) == nil {
			return transactions, nil
		}
	}

	// Fetch from DB if cache missed
	query := `SELECT id, user_id, account_id, amount, type, category_id, currency, description, created_at 
			  FROM transactions WHERE category_id = $1`
	rows, err := s.DB.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.AccountID, &t.Amount, &t.Type, &t.CategoryID, &t.Currency, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	// Cache results in Redis
	if data, err := json.Marshal(transactions); err == nil {
		s.RedisClient.Set(ctx, cacheKey, data, 0)
	}

	return transactions, nil
}

// GetTransactionsByAccount retrieves transactions by account, with optional caching
func (s *CategoryService) GetTransactionsByAccount(accountID int) ([]models.Transaction, error) {
	ctx := context.Background()
	cacheKey := "transactions:account:" + strconv.Itoa(accountID)

	// Check Redis cache
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var transactions []models.Transaction
		if json.Unmarshal([]byte(cachedData), &transactions) == nil {
			return transactions, nil
		}
	}

	// Fetch from DB if cache missed
	query := `SELECT id, user_id, account_id, amount, type, category_id, currency, description, created_at 
			  FROM transactions WHERE account_id = $1`
	rows, err := s.DB.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.AccountID, &t.Amount, &t.Type, &t.CategoryID, &t.Currency, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	// Cache results in Redis
	if data, err := json.Marshal(transactions); err == nil {
		s.RedisClient.Set(ctx, cacheKey, data, 0)
	}

	return transactions, nil
}
