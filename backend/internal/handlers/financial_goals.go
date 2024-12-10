package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"finance_project/internal/models"
	"finance_project/internal/services"

	"github.com/gorilla/mux"
)

type FinancialGoalsHandler struct {
	Service *services.FinancialGoalsService
}

// NewFinancialGoalsHandler создает новый обработчик для финансовых целей.
func NewFinancialGoalsHandler(service *services.FinancialGoalsService) *FinancialGoalsHandler {
	return &FinancialGoalsHandler{Service: service}
}

// GetFinancialGoalsHandler возвращает список всех финансовых целей пользователя.
// @Summary Список финансовых целей
// @Description Возвращает список всех финансовых целей для указанного пользователя
// @Tags Financial Goals
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} models.FinancialGoal "Список финансовых целей"
// @Failure 400 {string} string "Invalid user_id"
// @Failure 500 {string} string "Failed to retrieve financial goals"
// @Router /financial-goals [get]
func (h *FinancialGoalsHandler) GetFinancialGoalsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	goals, err := h.Service.GetFinancialGoalsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve financial goals", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(goals)
}

// CreateFinancialGoalHandler создает новую финансовую цель.
// @Summary Создание финансовой цели
// @Description Создает новую финансовую цель
// @Tags Financial Goals
// @Accept json
// @Produce json
// @Param goal body models.FinancialGoal true "Financial Goal body"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create financial goal"
// @Router /financial-goals/create [post]
func (h *FinancialGoalsHandler) CreateFinancialGoalHandler(w http.ResponseWriter, r *http.Request) {
	var goal models.FinancialGoal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateFinancialGoal(goal); err != nil {
		http.Error(w, "Failed to create financial goal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateFinancialGoalHandler обновляет данные финансовой цели.
// @Summary Обновление финансовой цели
// @Description Обновляет данные существующей финансовой цели
// @Tags Financial Goals
// @Accept json
// @Produce json
// @Param goal body models.FinancialGoal true "Financial Goal body"
// @Success 200 {string} string "Updated"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to update financial goal"
// @Router /financial-goals/update [put]
func (h *FinancialGoalsHandler) UpdateFinancialGoalHandler(w http.ResponseWriter, r *http.Request) {
	var goal models.FinancialGoal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateFinancialGoal(goal); err != nil {
		http.Error(w, "Failed to update financial goal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteFinancialGoalHandler удаляет финансовую цель.
// @Summary Удаление финансовой цели
// @Description Удаляет существующую финансовую цель по её ID
// @Tags Financial Goals
// @Accept json
// @Produce json
// @Param id query int true "Financial Goal ID"
// @Success 200 {string} string "Deleted"
// @Failure 400 {string} string "Invalid financial goal ID"
// @Failure 500 {string} string "Failed to delete financial goal"
// @Router /financial-goals/delete [delete]
func (h *FinancialGoalsHandler) DeleteFinancialGoalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid financial goal ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteFinancialGoal(id); err != nil {
		http.Error(w, "Failed to delete financial goal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetGoalProgressHandler возвращает прогресс финансовых целей.
// @Summary Получить прогресс выполнения финансовых целей
// @Description Возвращает процент выполнения для каждой финансовой цели пользователя
// @Tags Financial Goals
// @Param id path int true "User ID"
// @Success 200 {array} models.GoalProgress
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Failed to retrieve goal progress"
// @Router /users/{id}/goals/progress [get]
func (h *FinancialGoalsHandler) GetGoalProgressHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	goals, err := h.Service.GetGoalProgress(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve goal progress", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(goals)
}
