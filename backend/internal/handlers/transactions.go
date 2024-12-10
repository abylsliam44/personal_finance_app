package handlers

import (
	"encoding/json"
	"finance_project/internal/models"
	"finance_project/internal/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// TransactionHandler представляет обработчики для транзакций.
type TransactionHandler struct {
	Service *services.TransactionService
}

// NewTransactionHandler создает новый экземпляр TransactionHandler.
func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: service}
}

// CreateTransactionHandler добавляет новую транзакцию.
// @Summary Создать транзакцию
// @Description Добавляет новую транзакцию
// @Tags Transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction body"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create transaction"
// @Router /transactions/create [post]
func (h *TransactionHandler) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateTransaction(transaction); err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetTransactionByIDHandler возвращает транзакцию по ID.
// @Summary Получение транзакции по ID
// @Description Возвращает данные конкретной транзакции
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id query int true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Failure 400 {string} string "Invalid transaction ID"
// @Failure 404 {string} string "Transaction not found"
// @Failure 500 {string} string "Failed to retrieve transaction"
// @Router /transactions/get [get]
func (h *TransactionHandler) GetTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметра ID из запроса
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	// Получение транзакции из сервиса
	transaction, err := h.Service.GetTransactionByID(id)
	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	// Отправка успешного ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// GetAllTransactionsHandler возвращает список всех транзакций пользователя.
// @Summary Список транзакций
// @Description Возвращает все транзакции текущего пользователя
// @Tags Transactions
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} models.Transaction
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Failed to retrieve transactions"
// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	transactions, err := h.Service.GetAllTransactions(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// DeleteTransactionHandler удаляет транзакцию по ID.
// @Summary Удалить транзакцию
// @Description Удаляет транзакцию по ID
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id query int true "Transaction ID"
// @Success 200 {string} string "Deleted"
// @Failure 400 {string} string "Invalid transaction ID"
// @Failure 500 {string} string "Failed to delete transaction"
// @Router /transactions/delete [delete]
func (h *TransactionHandler) DeleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteTransaction(id); err != nil {
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetTransactionsByCategoryHandler возвращает транзакции по категории.
// @Summary Получить транзакции по категории
// @Description Возвращает список всех транзакций, связанных с заданной категорией
// @Tags Transactions
// @Param id path int true "Category ID"
// @Success 200 {array} models.Transaction
// @Failure 400 {string} string "Invalid category ID"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Failed to retrieve transactions"
// @Router /categories/{id}/transactions [get]
func (h *TransactionHandler) GetTransactionsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID категории из параметра пути
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil || categoryID <= 0 {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Проверка существования категории
	exists, err := h.Service.CategoryExists(categoryID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// Получение транзакций
	transactions, err := h.Service.GetTransactionsByCategory(categoryID)
	if err != nil {
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}

	// Ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// GetTransactionsByAccountHandler возвращает транзакции по ID счета.
// @Summary Получить транзакции по счёту
// @Description Возвращает список транзакций, связанных с указанным счётом
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {array} models.Transaction
// @Failure 400 {string} string "Invalid account ID"
// @Failure 404 {string} string "Account not found"
// @Failure 500 {string} string "Internal server error"
// @Router /accounts/{id}/transactions [get]
func (h *TransactionHandler) GetTransactionsByAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID счёта из параметра пути
	vars := mux.Vars(r)
	accountID, err := strconv.Atoi(vars["id"])
	if err != nil || accountID <= 0 {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	// Проверка существования счета
	exists, err := h.Service.AccountExists(accountID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Получение транзакций
	transactions, err := h.Service.GetTransactionsByAccount(accountID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// CompareIncomeAndExpensesHandler возвращает сравнение доходов и расходов за указанный период.
// @Summary Сравнение доходов и расходов
// @Description Возвращает общие доходы, расходы и баланс за указанный временной период
// @Tags Transactions
// @Param id path int true "User ID"
// @Param from query string true "Начало периода (YYYY-MM-DD)"
// @Param to query string true "Конец периода (YYYY-MM-DD)"
// @Success 200 {object} map[string]float64
// @Failure 400 {string} string "Invalid parameters"
// @Failure 500 {string} string "Failed to calculate income and expenses"
// @Router /users/{id}/transactions/compare [get]
func (h *TransactionHandler) CompareIncomeAndExpensesHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметров
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from == "" || to == "" {
		http.Error(w, "Missing 'from' or 'to' parameters", http.StatusBadRequest)
		return
	}

	// Вызов метода сервиса
	result, err := h.Service.CompareIncomeAndExpenses(userID, from, to)
	if err != nil {
		http.Error(w, "Failed to calculate income and expenses", http.StatusInternalServerError)
		return
	}

	// Возврат результата
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
