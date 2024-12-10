package services

import (
	"database/sql"
	"finance_project/internal/models"
	"log"
)

type AccountService struct {
	DB *sql.DB
}

// NewAccountService создаёт новый сервис для работы со счетами
func NewAccountService(db *sql.DB) *AccountService {
	return &AccountService{DB: db}
}

// CreateAccount добавляет новый счёт
func (s *AccountService) CreateAccount(account models.Account) error {
	query := `INSERT INTO accounts (user_id, name, balance, currency, type, created_at)
			  VALUES ($1, $2, $3, $4, $5, NOW())`
	_, err := s.DB.Exec(query, account.UserID, account.Name, account.Balance, account.Currency, account.Type)
	if err != nil {
		log.Printf("Error creating account: %v", err)
		return err
	}
	return nil
}

// GetAllAccounts возвращает все счета пользователя
func (s *AccountService) GetAllAccounts(userID int) ([]models.Account, error) {
	query := `SELECT id, user_id, name, balance, currency, type, created_at 
			  FROM accounts WHERE user_id = $1`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error retrieving accounts: %v", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Name, &account.Balance, &account.Currency, &account.Type, &account.CreatedAt); err != nil {
			log.Printf("Error scanning account: %v", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// GetAccountByID возвращает счёт по ID
func (s *AccountService) GetAccountByID(id int) (*models.Account, error) {
	query := `SELECT id, user_id, name, balance, currency, type, created_at 
			  FROM accounts WHERE id = $1`

	var account models.Account
	err := s.DB.QueryRow(query, id).Scan(&account.ID, &account.UserID, &account.Name, &account.Balance, &account.Currency, &account.Type, &account.CreatedAt)
	if err != nil {
		log.Printf("Error retrieving account by ID: %v", err)
		return nil, err
	}
	return &account, nil
}

// UpdateAccount обновляет данные счёта
func (s *AccountService) UpdateAccount(account models.Account) error {
	query := `UPDATE accounts SET name = $1, balance = $2, currency = $3, type = $4 WHERE id = $5`
	_, err := s.DB.Exec(query, account.Name, account.Balance, account.Currency, account.Type, account.ID)
	if err != nil {
		log.Printf("Error updating account: %v", err)
		return err
	}
	return nil
}

// DeleteAccount удаляет счёт
func (s *AccountService) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting account: %v", err)
		return err
	}
	return nil
}
