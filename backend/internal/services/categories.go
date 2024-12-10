package services

import (
	"database/sql"
	"log"

	"finance_project/internal/models"
)

type CategoryService struct {
	DB *sql.DB
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
