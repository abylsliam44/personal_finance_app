package app

import (
	"fmt"
	"log"
	"net/http"

	"finance_project/internal/config"
	"finance_project/internal/database"
	"finance_project/internal/handlers"
	"finance_project/internal/redis_client"
	"finance_project/internal/services"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Run(cfg *config.Config) {
	// Подключение к PostgreSQL
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Подключение к Redis
	redisClient := redis_client.NewRedisClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	// Initialize services
	userService := services.NewUserService(db)
	accountService := services.NewAccountService(db)
	transactionService := services.NewTransactionService(db, redisClient)
	categoryService := services.NewCategoryService(db)
	financialGoalsService := services.NewFinancialGoalsService(db)
	reportsService := services.NewReportsService(db)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	financialGoalsHandler := handlers.NewFinancialGoalsHandler(financialGoalsService)
	reportsHandler := handlers.NewReportsHandler(reportsService)

	// Create router
	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/users", userHandler.GetAllUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.GetUserByIDHandler).Methods("GET")
	r.HandleFunc("/users/update", userHandler.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/users/delete", userHandler.DeleteUserHandler).Methods("DELETE")

	// Account routes
	r.HandleFunc("/accounts", accountHandler.GetAccountsHandler).Methods("GET")
	r.HandleFunc("/accounts/create", accountHandler.CreateAccountHandler).Methods("POST")
	r.HandleFunc("/accounts/{id}", accountHandler.GetAccountByIDHandler).Methods("GET")
	r.HandleFunc("/accounts/update", accountHandler.UpdateAccountHandler).Methods("PUT")
	r.HandleFunc("/accounts/delete", accountHandler.DeleteAccountHandler).Methods("DELETE")

	// Transaction routes
	r.HandleFunc("/transactions", transactionHandler.GetAllTransactionsHandler).Methods("GET")
	r.HandleFunc("/transactions/create", transactionHandler.CreateTransactionHandler).Methods("POST")
	r.HandleFunc("/transactions/{id}", transactionHandler.GetTransactionByIDHandler).Methods("GET")
	r.HandleFunc("/transactions/delete", transactionHandler.DeleteTransactionHandler).Methods("DELETE")
	r.HandleFunc("/transactions/{userID}/cache", transactionHandler.GetAllTransactionsWithCacheHandler).Methods("GET")
	r.HandleFunc("/users/{id}/transactions/compare", transactionHandler.CompareIncomeAndExpensesHandler).Methods("GET")

	// Category routes
	r.HandleFunc("/categories", categoryHandler.GetAllCategoriesHandler).Methods("GET")
	r.HandleFunc("/categories/create", categoryHandler.CreateCategoryHandler).Methods("POST")
	r.HandleFunc("/categories/{id}", categoryHandler.GetCategoryByIDHandler).Methods("GET")
	r.HandleFunc("/categories/update", categoryHandler.UpdateCategoryHandler).Methods("PUT")
	r.HandleFunc("/categories/delete", categoryHandler.DeleteCategoryHandler).Methods("DELETE")
	r.HandleFunc("/categories/{id}/transactions", transactionHandler.GetTransactionsByCategoryHandler).Methods("GET")
	r.HandleFunc("/accounts/{id}/transactions", transactionHandler.GetTransactionsByAccountHandler).Methods("GET")

	// Financial goals routes
	r.HandleFunc("/financial-goals", financialGoalsHandler.GetFinancialGoalsHandler).Methods(http.MethodGet)
	r.HandleFunc("/financial-goals/create", financialGoalsHandler.CreateFinancialGoalHandler).Methods(http.MethodPost)
	r.HandleFunc("/financial-goals/update", financialGoalsHandler.UpdateFinancialGoalHandler).Methods(http.MethodPut)
	r.HandleFunc("/financial-goals/delete", financialGoalsHandler.DeleteFinancialGoalHandler).Methods(http.MethodDelete)
	r.HandleFunc("/users/{id}/goals/progress", financialGoalsHandler.GetGoalProgressHandler).Methods("GET")

	// Reports routes
	r.HandleFunc("/reports/summary", reportsHandler.GetSummaryHandler).Methods(http.MethodGet)
	r.HandleFunc("/reports/by-category", reportsHandler.GetExpensesByCategoryHandler).Methods(http.MethodGet)

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start server
	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	fmt.Println("Swagger docs available at http://localhost:8080/swagger/")
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
