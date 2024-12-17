package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

type ReportsService struct {
	DB *sql.DB
}

// NewReportsService создает новый сервис для отчетов.
func NewReportsService(db *sql.DB) *ReportsService {
	return &ReportsService{DB: db}
}

// GetOrCreateSummaryReport возвращает существующий отчет или создает новый.
func (s *ReportsService) GetOrCreateSummaryReport(userID int) (map[string]interface{}, error) {
	// Проверяем, есть ли уже созданный отчет
	var reportData string
	query := `SELECT data FROM reports WHERE user_id = $1 AND report_name = 'summary' ORDER BY generated_at DESC LIMIT 1`
	err := s.DB.QueryRow(query, userID).Scan(&reportData)
	if err == nil {
		// Если отчет найден, возвращаем его
		var report map[string]interface{}
		if err := json.Unmarshal([]byte(reportData), &report); err != nil {
			log.Printf("Error unmarshalling report data: %v", err)
			return nil, err
		}
		return report, nil
	}

	// Если отчет не найден, создаем новый
	report, err := s.GenerateSummaryReport(userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем новый отчет в таблицу
	reportJSON, err := json.Marshal(report)
	if err != nil {
		log.Printf("Error marshalling report data: %v", err)
		return nil, err
	}

	insertQuery := `INSERT INTO reports (user_id, report_name, generated_at, data) VALUES ($1, $2, $3, $4)`
	_, err = s.DB.Exec(insertQuery, userID, "summary", time.Now(), reportJSON)
	if err != nil {
		log.Printf("Error inserting report: %v", err)
		return nil, err
	}

	return report, nil
}

// GenerateSummaryReport создает сводный отчет.
func (s *ReportsService) GenerateSummaryReport(userID int) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	// Общий баланс по счетам
	var totalBalance float64
	err := s.DB.QueryRow(`SELECT COALESCE(SUM(balance), 0) FROM accounts WHERE user_id = $1`, userID).Scan(&totalBalance)
	if err != nil {
		log.Printf("Error fetching total balance: %v", err)
		return nil, err
	}
	summary["total_balance"] = totalBalance

	// Расходы за текущий месяц
	var totalExpenses float64
	err = s.DB.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE user_id = $1 AND type = 'expense' AND EXTRACT(MONTH FROM created_at) = EXTRACT(MONTH FROM CURRENT_DATE)
	`, userID).Scan(&totalExpenses)
	if err != nil {
		log.Printf("Error fetching total expenses: %v", err)
		return nil, err
	}
	summary["total_expenses"] = totalExpenses

	// Выполненные финансовые цели
	var completedGoals int
	query := `
	SELECT COUNT(*) 
	FROM financial_goals 
	WHERE user_id = $1 AND target_amount <= saved_amount AND deadline >= CURRENT_DATE
`

	err = s.DB.QueryRow(query, userID).Scan(&completedGoals)
	if err != nil {
		log.Printf("Error fetching completed goals: %v", err)
		return nil, err
	}

	// Добавляем выполненные цели в отчёт
	summary["completed_goals"] = completedGoals
	return summary, nil
}

// GetExpensesByCategory возвращает расходы, сгруппированные по категориям.
func (s *ReportsService) GetExpensesByCategory(userID int, startDate, endDate string) (map[string]float64, error) {
	query := `
		SELECT c.name AS category, COALESCE(SUM(t.amount), 0) AS total_expenses
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = $1 AND t.created_at BETWEEN $2 AND $3 AND t.type = 'expense'
		GROUP BY c.name
	`
	rows, err := s.DB.Query(query, userID, startDate, endDate)
	if err != nil {
		log.Printf("Error fetching expenses by category: %v", err)
		return nil, err
	}
	defer rows.Close()

	expenses := make(map[string]float64)
	for rows.Next() {
		var category string
		var total float64
		if err := rows.Scan(&category, &total); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		expenses[category] = total
	}

	return expenses, nil
}
