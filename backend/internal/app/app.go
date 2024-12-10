package app

import (
	"fmt"
	"log"
	"net/http"

	"finance_project/internal/config"
	"finance_project/internal/database"
	"finance_project/internal/handlers"
	"finance_project/internal/middleware"
	"finance_project/internal/services"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Run(cfg *config.Config) {
	// Подключение к базе данных
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Выполнение миграций
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Инициализация сервисов
	userService := services.NewUserService(db)
	accountService := services.NewAccountService(db)
	transactionService := services.NewTransactionService(db)
	categoryService := services.NewCategoryService(db)
	financialGoalsService := services.NewFinancialGoalsService(db)
	reportsService := services.NewReportsService(db)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	financialGoalsHandler := handlers.NewFinancialGoalsHandler(financialGoalsService)
	reportsHandler := handlers.NewReportsHandler(reportsService)

	// Создание маршрутизатора
	r := mux.NewRouter()

	// Роуты для авторизации (не защищенные JWT)
	r.HandleFunc("/auth/login", userHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/register", userHandler.RegisterHandler).Methods("POST")

	// Группировка роутов, защищенных JWT
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware) // Добавляем middleware для авторизации

	// Роуты для пользователей
	api.HandleFunc("/users", userHandler.GetAllUsersHandler).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUserByIDHandler).Methods("GET")
	api.HandleFunc("/users/update", userHandler.UpdateUserHandler).Methods("PUT")
	api.HandleFunc("/users/delete", userHandler.DeleteUserHandler).Methods("DELETE")

	// Роуты для счетов
	api.HandleFunc("/accounts", accountHandler.GetAccountsHandler).Methods("GET")
	api.HandleFunc("/accounts/create", accountHandler.CreateAccountHandler).Methods("POST")
	api.HandleFunc("/accounts/{id}", accountHandler.GetAccountByIDHandler).Methods("GET")
	api.HandleFunc("/accounts/update", accountHandler.UpdateAccountHandler).Methods("PUT")
	api.HandleFunc("/accounts/delete", accountHandler.DeleteAccountHandler).Methods("DELETE")

	// Роуты для транзакций
	api.HandleFunc("/transactions", transactionHandler.GetAllTransactionsHandler).Methods("GET")
	api.HandleFunc("/transactions/create", transactionHandler.CreateTransactionHandler).Methods("POST")
	api.HandleFunc("/transactions/{id}", transactionHandler.GetTransactionByIDHandler).Methods("GET")
	api.HandleFunc("/transactions/delete", transactionHandler.DeleteTransactionHandler).Methods("DELETE")
	api.HandleFunc("/users/{id}/transactions/compare", transactionHandler.CompareIncomeAndExpensesHandler).Methods("GET")

	// Роуты для категорий
	api.HandleFunc("/categories", categoryHandler.GetAllCategoriesHandler).Methods("GET")
	api.HandleFunc("/categories/create", categoryHandler.CreateCategoryHandler).Methods("POST")
	api.HandleFunc("/categories/{id}", categoryHandler.GetCategoryByIDHandler).Methods("GET")
	api.HandleFunc("/categories/update", categoryHandler.UpdateCategoryHandler).Methods("PUT")
	api.HandleFunc("/categories/delete", categoryHandler.DeleteCategoryHandler).Methods("DELETE")
	api.HandleFunc("/categories/{id}/transactions", transactionHandler.GetTransactionsByCategoryHandler).Methods("GET")
	api.HandleFunc("/accounts/{id}/transactions", transactionHandler.GetTransactionsByAccountHandler).Methods("GET")

	// Роуты для финансовых целей
	api.HandleFunc("/financial-goals", financialGoalsHandler.GetFinancialGoalsHandler).Methods(http.MethodGet)
	api.HandleFunc("/financial-goals/create", financialGoalsHandler.CreateFinancialGoalHandler).Methods(http.MethodPost)
	api.HandleFunc("/financial-goals/update", financialGoalsHandler.UpdateFinancialGoalHandler).Methods(http.MethodPut)
	api.HandleFunc("/financial-goals/delete", financialGoalsHandler.DeleteFinancialGoalHandler).Methods(http.MethodDelete)
	api.HandleFunc("/users/{id}/goals/progress", financialGoalsHandler.GetGoalProgressHandler).Methods("GET")

	// Роуты для отчетов
	api.HandleFunc("/reports/summary", reportsHandler.GetSummaryHandler).Methods(http.MethodGet)
	api.HandleFunc("/reports/by-category", reportsHandler.GetExpensesByCategoryHandler).Methods(http.MethodGet)

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	fmt.Println("Swagger docs available at http://localhost:8080/swagger/")
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
