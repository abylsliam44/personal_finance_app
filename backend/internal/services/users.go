package services

import (
	"database/sql"
	"errors"
	"finance_project/internal/models"
	"log"
)

// UserService предоставляет методы для работы с пользователями.
type UserService struct {
	DB *sql.DB
}

// RegisterUser регистрирует нового пользователя.
func (s *UserService) RegisterUser(user models.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, preferred_currency, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := s.DB.Exec(query, user.Name, user.Email, user.PasswordHash, user.PreferredCurrency)
	return err
}

// Authenticate аутентифицирует пользователя.
func (s *UserService) Authenticate(email, password string) (int, error) {
	var userID int
	var storedPasswordHash string

	query := `SELECT id, password_hash FROM users WHERE email = $1`
	err := s.DB.QueryRow(query, email).Scan(&userID, &storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}
		return 0, err
	}

	if storedPasswordHash != password { // Добавьте реальное хеширование
		return 0, errors.New("invalid credentials")
	}

	return userID, nil
}

// NewUserService создает новый сервис пользователей.
func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

// CreateUser добавляет нового пользователя в базу данных.
func (s *UserService) CreateUser(user models.User) error {
	query := `INSERT INTO users (name, email, password_hash, preferred_currency, created_at)
			  VALUES ($1, $2, $3, $4, NOW())`
	_, err := s.DB.Exec(query, user.Name, user.Email, user.PasswordHash, user.PreferredCurrency)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}

// GetAllUsers возвращает список всех пользователей.
func (s *UserService) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, name, email, preferred_currency, created_at FROM users`

	rows, err := s.DB.Query(query)
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PreferredCurrency, &user.CreatedAt); err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUserByID возвращает пользователя по ID.
func (s *UserService) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, name, email, preferred_currency, created_at FROM users WHERE id = $1`

	var user models.User
	err := s.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.PreferredCurrency, &user.CreatedAt)
	if err != nil {
		log.Printf("Error retrieving user by ID: %v", err)
		return nil, err
	}
	return &user, nil
}

// UpdateUser обновляет информацию о пользователе.
func (s *UserService) UpdateUser(user models.User) error {
	query := `UPDATE users SET name = $1, email = $2, preferred_currency = $3 WHERE id = $4`
	_, err := s.DB.Exec(query, user.Name, user.Email, user.PreferredCurrency, user.ID)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}
	return nil
}

// DeleteUser удаляет пользователя.
func (s *UserService) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}
	return nil
}
