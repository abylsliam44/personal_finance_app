package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"finance_project/internal/models"
	"finance_project/internal/services"

	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	Service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: service}
}

// GetAllTransactionsHandler godoc
// @Summary Retrieve all transactions
// @Description Retrieves all transactions for a specific user
// @Tags Transactions
// @Param userID query int true "User ID"
// @Success 200 {array} models.Transaction
// @Failure 400 {string} string "Invalid User ID"
// @Failure 500 {string} string "Failed to retrieve transactions"
// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("userID"))
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
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

// GetAllTransactionsWithCacheHandler godoc
// @Summary Retrieve all transactions with cache
// @Description Retrieves all transactions for a user, using Redis caching
// @Tags Transactions
// @Param userID path int true "User ID"
// @Success 200 {array} models.Transaction
// @Failure 400 {string} string "Invalid User ID"
// @Failure 500 {string} string "Failed to retrieve transactions"
// @Router /transactions/{userID}/cache [get]
func (h *TransactionHandler) GetAllTransactionsWithCacheHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	transactions, err := h.Service.GetAllTransactionsWithCache(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}

	// Определяем источник данных
	source := "database"
	ctx := context.Background()
	cacheKey := fmt.Sprintf("transactions:user:%d", userID)
	if _, err := h.Service.RedisClient.Get(ctx, cacheKey).Result(); err == nil {
		source = "cache (Redis)"
	}

	// Формируем ответ
	response := map[string]interface{}{
		"source":       source,
		"transactions": transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateTransactionHandler godoc
// @Summary Create a transaction
// @Description Creates a new transaction
// @Tags Transactions
// @Param transaction body models.Transaction true "Transaction Data"
// @Success 201 {string} string "Transaction created successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Failed to create transaction"
// @Router /transactions/create [post]
func (h *TransactionHandler) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateTransaction(transaction); err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetTransactionByIDHandler godoc
// @Summary Retrieve transaction by ID
// @Description Retrieves a transaction using its ID
// @Tags Transactions
// @Param id path int true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Failure 400 {string} string "Invalid transaction ID"
// @Failure 404 {string} string "Transaction not found"
// @Failure 500 {string} string "Internal server error"
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.Service.GetTransactionByID(id)
	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// DeleteTransactionHandler godoc
// @Summary Delete a transaction
// @Description Deletes a transaction using its ID
// @Tags Transactions
// @Param id path int true "Transaction ID"
// @Success 204 {string} string "Transaction deleted successfully"
// @Failure 400 {string} string "Invalid transaction ID"
// @Failure 500 {string} string "Failed to delete transaction"
// @Router /transactions/delete [delete]
func (h *TransactionHandler) DeleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteTransaction(id); err != nil {
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CompareIncomeAndExpensesHandler godoc
// @Summary Compare income and expenses
// @Description Compares income and expenses for a user
// @Tags Transactions
// @Param id path int true "User ID"
// @Success 200 {object} map[string]float64 "Comparison of income and expenses"
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id}/transactions/compare [get]
func (h *TransactionHandler) CompareIncomeAndExpensesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	result, err := h.Service.CompareIncomeAndExpenses(userID)
	if err != nil {
		http.Error(w, "Failed to compare income and expenses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetTransactionsByCategoryHandler godoc
// @Summary Get transactions by category
// @Description Retrieves all transactions for a specific category, with optional caching
// @Tags Categories
// @Param categoryID path int true "Category ID"
// @Success 200 {array} services.Transaction "List of transactions"
// @Failure 400 {string} string "Invalid Category ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories/{id}/transactions [get]
func (h *CategoryHandler) GetTransactionsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	transactions, err := h.Service.GetTransactionsByCategory(categoryID)
	if err != nil {
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// GetTransactionsByAccountHandler godoc
// @Summary Get transactions by account
// @Description Retrieves all transactions for a specific account, with optional caching
// @Tags Categories
// @Param accountID path int true "Account ID"
// @Success 200 {array} services.Transaction "List of transactions"
// @Failure 400 {string} string "Invalid Account ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /accounts/{id}/transactions [get]
func (h *CategoryHandler) GetTransactionsByAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Account ID", http.StatusBadRequest)
		return
	}

	transactions, err := h.Service.GetTransactionsByAccount(accountID)
	if err != nil {
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
