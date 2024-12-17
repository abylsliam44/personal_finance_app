package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"finance_project/internal/models"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type TransactionService struct {
	DB          *sql.DB
	RedisClient *redis.Client
}

func NewTransactionService(db *sql.DB, redisClient *redis.Client) *TransactionService {
	return &TransactionService{
		DB:          db,
		RedisClient: redisClient,
	}
}

// GetAllTransactions retrieves all transactions for a specific user
func (s *TransactionService) GetAllTransactions(userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := `SELECT id, user_id, account_id, amount, type, category_id, currency, description, created_at 
			FROM transactions WHERE user_id = $1`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.AccountID, &t.Amount, &t.Type, &t.CategoryID, &t.Currency, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// GetAllTransactionsWithCache retrieves all transactions for a user, with caching
func (s *TransactionService) GetAllTransactionsWithCache(userID int) ([]models.Transaction, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("transactions:user:%d", userID) // Формируем ключ для кэша

	// Попытка получить данные из Redis
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var transactions []models.Transaction
		if err := json.Unmarshal([]byte(cachedData), &transactions); err == nil {
			log.Printf("Cache hit for key: %s", cacheKey) // Логируем, что кэш был использован
			return transactions, nil
		}
	}

	// Если данных в кэше нет, получаем их из базы данных
	log.Printf("Cache miss for key: %s", cacheKey) // Логируем, что кэш отсутствует
	transactions, err := s.GetAllTransactions(userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем данные в Redis
	data, err := json.Marshal(transactions)
	if err == nil {
		s.RedisClient.Set(ctx, cacheKey, data, 0)
		log.Printf("Data cached for key: %s", cacheKey) // Логируем, что данные закэшированы
	}

	return transactions, nil
}

// CreateTransaction adds a new transaction to the database
func (s *TransactionService) CreateTransaction(transaction models.Transaction) error {
	query := `INSERT INTO transactions (user_id, account_id, amount, type, category_id, currency, description, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.DB.Exec(query, transaction.UserID, transaction.AccountID, transaction.Amount, transaction.Type,
		transaction.CategoryID, transaction.Currency, transaction.Description, transaction.CreatedAt)

	return err
}

// GetTransactionByID retrieves a transaction by its ID
func (s *TransactionService) GetTransactionByID(id int) (*models.Transaction, error) {
	var transaction models.Transaction

	query := `SELECT id, user_id, account_id, amount, type, category_id, currency, description, created_at 
			FROM transactions WHERE id = $1`
	err := s.DB.QueryRow(query, id).Scan(&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.Amount,
		&transaction.Type, &transaction.CategoryID, &transaction.Currency, &transaction.Description, &transaction.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("transaction not found")
	} else if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// DeleteTransaction deletes a transaction by its ID
func (s *TransactionService) DeleteTransaction(id int) error {
	query := `DELETE FROM transactions WHERE id = $1`
	result, err := s.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return nil
}

// CompareIncomeAndExpenses compares income and expenses for a user
func (s *TransactionService) CompareIncomeAndExpenses(userID int) (map[string]float64, error) {
	query := `SELECT type, SUM(amount) FROM transactions WHERE user_id = $1 GROUP BY type`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]float64{"income": 0, "expense": 0}
	for rows.Next() {
		var tType string
		var total float64
		if err := rows.Scan(&tType, &total); err != nil {
			return nil, err
		}
		result[tType] = total
	}

	return result, nil
}

// GetTransactionsByCategory retrieves transactions by category ID
func (s *TransactionService) GetTransactionsByCategory(categoryID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := `SELECT id, user_id, account_id, amount, type, category_id, currency, description, created_at 
			FROM transactions WHERE category_id = $1`
	rows, err := s.DB.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.AccountID, &t.Amount, &t.Type, &t.CategoryID, &t.Currency, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// GetTransactionsByAccount retrieves transactions by account ID
func (s *TransactionService) GetTransactionsByAccount(accountID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := `SELECT id, user_id, account_id, amount, type, category_id, currency, description, created_at 
			FROM transactions WHERE account_id = $1`
	rows, err := s.DB.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.AccountID, &t.Amount, &t.Type, &t.CategoryID, &t.Currency, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
