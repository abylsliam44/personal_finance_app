package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"finance_project/internal/services"
)

type ReportsHandler struct {
	Service *services.ReportsService
}

// NewReportsHandler создает новый обработчик для отчетов.
func NewReportsHandler(service *services.ReportsService) *ReportsHandler {
	return &ReportsHandler{Service: service}
}

// GetSummaryHandler возвращает или создает сводный отчет.
// @Summary Сводный отчет
// @Description Возвращает общий отчет или создает новый и сохраняет его в таблицу reports
// @Tags Reports
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Failed to retrieve or generate summary report"
// @Router /reports/summary [get]
func (h *ReportsHandler) GetSummaryHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из параметров запроса
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Получаем или создаем сводный отчет
	report, err := h.Service.GetOrCreateSummaryReport(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve or generate summary report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetExpensesByCategoryHandler возвращает расходы, сгруппированные по категориям за период.
// @Summary Расходы по категориям
// @Description Возвращает расходы, сгруппированные по категориям за указанный период
// @Tags Reports
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {array} map[string]float64
// @Failure 400 {string} string "Invalid parameters"
// @Failure 500 {string} string "Failed to retrieve expenses"
// @Router /reports/by-category [get]
func (h *ReportsHandler) GetExpensesByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из запроса
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	if startDate == "" || endDate == "" {
		http.Error(w, "Invalid date range", http.StatusBadRequest)
		return
	}

	// Проверка формата дат
	_, err = time.Parse("2006-01-02", startDate)
	if err != nil {
		http.Error(w, "Invalid start_date format", http.StatusBadRequest)
		return
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		http.Error(w, "Invalid end_date format", http.StatusBadRequest)
		return
	}

	// Получаем данные из сервиса
	expenses, err := h.Service.GetExpensesByCategory(userID, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to retrieve expenses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}
