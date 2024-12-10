package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"finance_project/internal/models"
	"finance_project/internal/services"
)

// AccountHandler представляет обработчики для счетов
type AccountHandler struct {
	Service *services.AccountService
}

// NewAccountHandler создаёт новый обработчик
func NewAccountHandler(service *services.AccountService) *AccountHandler {
	return &AccountHandler{Service: service}
}

// CreateAccountHandler добавляет новый счёт
// @Summary Создание счёта
// @Description Добавляет новый счёт
// @Tags Accounts
// @Accept json
// @Produce json
// @Param account body models.Account true "Account body"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create account"
// @Router /accounts/create [post]
func (h *AccountHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateAccount(account); err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAccountByIDHandler возвращает информацию о счёте по ID
// @Summary Информация о счёте
// @Description Возвращает информацию о счёте по его ID
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id query int true "Account ID"
// @Success 200 {object} models.Account
// @Failure 400 {string} string "Invalid account ID"
// @Failure 404 {string} string "Account not found"
// @Router /accounts/get [get]
func (h *AccountHandler) GetAccountByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	account, err := h.Service.GetAccountByID(id)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// GetAllAccountsHandler возвращает все счета пользователя
// @Summary Список счетов
// @Description Возвращает список всех счетов пользователя
// @Tags Accounts
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} models.Account
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Failed to retrieve accounts"
// @Router /accounts [get]
func (h *AccountHandler) GetAccountsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	accounts, err := h.Service.GetAllAccounts(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// UpdateAccountHandler обновляет данные счёта
// @Summary Обновление счёта
// @Description Обновляет данные существующего счёта
// @Tags Accounts
// @Accept json
// @Produce json
// @Param account body models.Account true "Account body"
// @Success 200 {string} string "Updated"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to update account"
// @Router /accounts/update [put]
func (h *AccountHandler) UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateAccount(account); err != nil {
		http.Error(w, "Failed to update account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteAccountHandler удаляет счёт
// @Summary Удаление счёта
// @Description Удаляет существующий счёт по его ID
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id query int true "Account ID"
// @Success 200 {string} string "Deleted"
// @Failure 400 {string} string "Invalid account ID"
// @Failure 500 {string} string "Failed to delete account"
// @Router /accounts/delete [delete]
func (h *AccountHandler) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteAccount(id); err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
