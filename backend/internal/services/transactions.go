package services

import (
	"database/sql"
	"errors"
	"finance_project/internal/models"
	"log"
)

type TransactionService struct {
	DB *sql.DB
}

// NewTransactionService создает новый сервис для работы с транзакциями.
func NewTransactionService(db *sql.DB) *TransactionService {
	return &TransactionService{DB: db}
}

// CreateTransaction добавляет новую транзакцию.
func (s *TransactionService) CreateTransaction(transaction models.Transaction) error {
	// Проверка существования пользователя
	if exists, err := s.userExists(transaction.UserID); err != nil || !exists {
		if err != nil {
			log.Printf("Error checking user existence: %v", err)
			return errors.New("failed to validate user existence")
		}
		return errors.New("user does not exist")
	}

	// Проверка существования аккаунта
	if exists, err := s.accountExists(transaction.AccountID); err != nil || !exists {
		if err != nil {
			log.Printf("Error checking account existence: %v", err)
			return errors.New("failed to validate account existence")
		}
		return errors.New("account does not exist")
	}

	// Проверка существования категории
	if exists, err := s.categoryExists(transaction.CategoryID); err != nil || !exists {
		if err != nil {
			log.Printf("Error checking category existence: %v", err)
			return errors.New("failed to validate category existence")
		}
		return errors.New("category does not exist")
	}

	// SQL-запрос на добавление транзакции
	query := `
		INSERT INTO transactions 
		(user_id, account_id, category_id, amount, currency, type, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`

	_, err := s.DB.Exec(query, transaction.UserID, transaction.AccountID, transaction.CategoryID,
		transaction.Amount, transaction.Currency, transaction.Type, transaction.Description)

	if err != nil {
		log.Printf("Error creating transaction: %v", err)
		return errors.New("failed to create transaction")
	}

	log.Println("Transaction created successfully")
	return nil
}

// userExists проверяет, существует ли пользователь.
func (s *TransactionService) userExists(userID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	err := s.DB.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return false, err
	}
	return exists, nil
}

// accountExists проверяет, существует ли аккаунт.
func (s *TransactionService) accountExists(accountID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM accounts WHERE id = $1)`
	var exists bool
	err := s.DB.QueryRow(query, accountID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking account existence: %v", err)
		return false, err
	}
	return exists, nil
}

// categoryExists проверяет, существует ли категория.
func (s *TransactionService) categoryExists(categoryID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1)`
	var exists bool
	err := s.DB.QueryRow(query, categoryID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking category existence: %v", err)
		return false, err
	}
	return exists, nil
}

// GetAllTransactions возвращает все транзакции для заданного пользователя.
func (s *TransactionService) GetAllTransactions(userID int) ([]models.Transaction, error) {
	query := `SELECT id, user_id, account_id, category_id, amount, currency, type, description, created_at 
			  FROM transactions WHERE user_id = $1`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error retrieving transactions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.AccountID, &t.CategoryID, &t.Amount, &t.Currency, &t.Type, &t.Description, &t.CreatedAt)
		if err != nil {
			log.Printf("Error scanning transaction row: %v", err)
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating transaction rows: %v", rows.Err())
		return nil, rows.Err()
	}

	return transactions, nil
}

// GetTransactionByID возвращает транзакцию по ID.
func (s *TransactionService) GetTransactionByID(id int) (*models.Transaction, error) {
	query := `SELECT id, user_id, account_id, category_id, amount, currency, type, description, created_at 
			  FROM transactions WHERE id = $1`
	row := s.DB.QueryRow(query, id)

	var t models.Transaction
	err := row.Scan(&t.ID, &t.UserID, &t.AccountID, &t.CategoryID, &t.Amount, &t.Currency, &t.Type, &t.Description, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Transaction with ID %d not found", id)
			return nil, nil
		}
		log.Printf("Error retrieving transaction: %v", err)
		return nil, err
	}
	return &t, nil
}

// UpdateTransaction обновляет существующую транзакцию.
func (s *TransactionService) UpdateTransaction(transaction models.Transaction) error {
	// Проверка существования пользователя
	userExists, err := s.userExists(transaction.UserID)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return errors.New("failed to validate user existence")
	}
	if !userExists {
		return errors.New("user does not exist")
	}

	// Проверка существования аккаунта
	accountExists, err := s.accountExists(transaction.AccountID)
	if err != nil {
		log.Printf("Error checking account existence: %v", err)
		return errors.New("failed to validate account existence")
	}
	if !accountExists {
		return errors.New("account does not exist")
	}

	// Проверка существования категории
	categoryExists, err := s.categoryExists(transaction.CategoryID)
	if err != nil {
		log.Printf("Error checking category existence: %v", err)
		return errors.New("failed to validate category existence")
	}
	if !categoryExists {
		return errors.New("category does not exist")
	}

	// SQL-запрос на обновление транзакции
	query := `
		UPDATE transactions 
		SET user_id = $1, account_id = $2, category_id = $3, 
		    amount = $4, currency = $5, type = $6, description = $7, created_at = NOW()
		WHERE id = $8
	`

	_, err = s.DB.Exec(query, transaction.UserID, transaction.AccountID, transaction.CategoryID,
		transaction.Amount, transaction.Currency, transaction.Type, transaction.Description, transaction.ID)

	if err != nil {
		log.Printf("Error updating transaction: %v", err)
		return errors.New("failed to update transaction")
	}

	log.Println("Transaction updated successfully")
	return nil
}

// DeleteTransaction удаляет транзакцию по ID.
func (s *TransactionService) DeleteTransaction(id int) error {
	query := `DELETE FROM transactions WHERE id = $1`

	_, err := s.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting transaction: %v", err)
		return err
	}
	return nil
}

// GetTransactionsByCategory возвращает транзакции, связанные с категорией.
func (s *TransactionService) GetTransactionsByCategory(categoryID int) ([]models.Transaction, error) {
	query := `
		SELECT t.id, t.user_id, t.account_id, t.category_id, t.amount, t.currency, t.type, t.description, t.created_at
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		WHERE c.id = $1;
	`

	rows, err := s.DB.Query(query, categoryID)
	if err != nil {
		log.Printf("Error retrieving transactions by category: %v", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.AccountID,
			&transaction.CategoryID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Type,
			&transaction.Description,
			&transaction.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning transaction row: %v", err)
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating transaction rows: %v", rows.Err())
		return nil, rows.Err()
	}

	return transactions, nil
}

// GetTransactionsByAccount возвращает список транзакций по счёту.
func (s *TransactionService) GetTransactionsByAccount(accountID int) ([]models.Transaction, error) {
	// SQL-запрос для получения транзакций по счёту
	query := `
		SELECT id, user_id, account_id, category_id, amount, currency, type, description, created_at
		FROM transactions
		WHERE account_id = $1
	`

	// Выполняем запрос
	rows, err := s.DB.Query(query, accountID)
	if err != nil {
		log.Printf("Error retrieving transactions by account ID %d: %v", accountID, err)
		return nil, err
	}
	defer rows.Close()

	// Сканируем результаты
	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.AccountID,
			&transaction.CategoryID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Type,
			&transaction.Description,
			&transaction.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning transaction row for account ID %d: %v", accountID, err)
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	// Проверяем на ошибки при итерации
	if rows.Err() != nil {
		log.Printf("Error iterating transaction rows for account ID %d: %v", accountID, rows.Err())
		return nil, rows.Err()
	}

	// Если транзакций нет, возвращаем пустой список
	if len(transactions) == 0 {
		log.Printf("No transactions found for account ID %d", accountID)
	}

	return transactions, nil
}

// AccountExists проверяет, существует ли счёт в базе данных.
func (s *TransactionService) AccountExists(accountID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM accounts WHERE id = $1)`
	var exists bool
	err := s.DB.QueryRow(query, accountID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking account existence: %v", err)
		return false, err
	}
	return exists, nil
}

// CategoryExists проверяет, существует ли категория в базе данных.
func (s *TransactionService) CategoryExists(categoryID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1)`
	var exists bool
	err := s.DB.QueryRow(query, categoryID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking category existence: %v", err)
		return false, err
	}
	return exists, nil
}

// CompareIncomeAndExpenses рассчитывает доходы, расходы и баланс за указанный период.
func (s *TransactionService) CompareIncomeAndExpenses(userID int, from, to string) (map[string]float64, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expense
		FROM transactions
		WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
	`
	var totalIncome, totalExpense float64

	err := s.DB.QueryRow(query, userID, from, to).Scan(&totalIncome, &totalExpense)
	if err != nil {
		log.Printf("Error comparing income and expenses: %v", err)
		return nil, err
	}

	result := map[string]float64{
		"total_income":  totalIncome,
		"total_expense": totalExpense,
		"balance":       totalIncome - totalExpense,
	}

	return result, nil
}
