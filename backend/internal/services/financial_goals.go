package services

import (
	"database/sql"
	"log"

	"finance_project/internal/models"
)

type FinancialGoalsService struct {
	DB *sql.DB
}

// NewFinancialGoalsService создает новый сервис для работы с финансовыми целями.
func NewFinancialGoalsService(db *sql.DB) *FinancialGoalsService {
	return &FinancialGoalsService{DB: db}
}

// GetFinancialGoalsByUserID возвращает финансовые цели пользователя по user_id.
func (s *FinancialGoalsService) GetFinancialGoalsByUserID(userID int) ([]models.FinancialGoal, error) {
	query := `SELECT id, user_id, name, target_amount, saved_amount, deadline, priority, description, created_at 
	          FROM financial_goals WHERE user_id = $1`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error retrieving financial goals: %v", err)
		return nil, err
	}
	defer rows.Close()

	var goals []models.FinancialGoal
	for rows.Next() {
		var g models.FinancialGoal
		err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.TargetAmount, &g.CurrentAmount, &g.Deadline, &g.Priority, &g.Description, &g.CreatedAt)
		if err != nil {
			log.Printf("Error scanning financial goal row: %v", err)
			return nil, err
		}
		goals = append(goals, g)
	}

	return goals, nil
}

// CreateFinancialGoal добавляет новую финансовую цель в базу данных.
func (s *FinancialGoalsService) CreateFinancialGoal(goal models.FinancialGoal) error {
	query := `INSERT INTO financial_goals (user_id, name, target_amount, saved_amount, deadline, priority, description, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`
	_, err := s.DB.Exec(query, goal.UserID, goal.Name, goal.TargetAmount, goal.CurrentAmount, goal.Deadline, goal.Priority, goal.Description)
	if err != nil {
		log.Printf("Error creating financial goal: %v", err)
		return err
	}
	return nil
}

// UpdateFinancialGoal обновляет данные финансовой цели в базе данных.
func (s *FinancialGoalsService) UpdateFinancialGoal(goal models.FinancialGoal) error {
	query := `UPDATE financial_goals 
	          SET name = $1, target_amount = $2, saved_amount = $3, deadline = $4, priority = $5, description = $6 
	          WHERE id = $7`
	_, err := s.DB.Exec(query, goal.Name, goal.TargetAmount, goal.CurrentAmount, goal.Deadline, goal.Priority, goal.Description, goal.ID)
	if err != nil {
		log.Printf("Error updating financial goal: %v", err)
		return err
	}
	return nil
}

// DeleteFinancialGoal удаляет финансовую цель из базы данных по ID.
func (s *FinancialGoalsService) DeleteFinancialGoal(id int) error {
	query := `DELETE FROM financial_goals WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting financial goal: %v", err)
		return err
	}
	return nil
}

// GetGoalProgress возвращает прогресс выполнения финансовых целей.
func (s *FinancialGoalsService) GetGoalProgress(userID int) ([]models.GoalProgress, error) {
	query := `
		SELECT id, name, target_amount, saved_amount, 
		       (saved_amount * 100.0 / target_amount) AS progress
		FROM financial_goals
		WHERE user_id = $1;
	`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error retrieving goal progress: %v", err)
		return nil, err
	}
	defer rows.Close()

	var goals []models.GoalProgress
	for rows.Next() {
		var goal models.GoalProgress
		err := rows.Scan(&goal.ID, &goal.Name, &goal.TargetAmount, &goal.SavedAmount, &goal.Progress)
		if err != nil {
			log.Printf("Error scanning goal row: %v", err)
			return nil, err
		}
		goals = append(goals, goal)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating goal rows: %v", rows.Err())
		return nil, rows.Err()
	}

	return goals, nil
}
