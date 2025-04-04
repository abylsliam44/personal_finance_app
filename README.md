# Finance Management Backend Project

## Overview
This project is a backend implementation for a **personal finance management system**. Designed with scalability, clean architecture principles, and high performance in mind, it leverages **Go (Golang)** to manage user transactions, financial goals, and insightful reports while integrating Redis caching for improved speed and efficiency. The API is well-documented with Swagger, making it easy to explore and test.

The system is ideal for individuals and small businesses who want to monitor their finances, track progress toward goals, and make data-driven decisions.

## Project Structure
```plaintext
.
├── cmd                 # Entry point for the application
│   ├── docs            # Swagger documentation
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   └── main.go         # Main application file
├── configs             # Configuration files
│   └── config.yaml     # Database and Redis configurations
├── docs                # Additional Swagger docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal            # Internal application components
│   ├── app             # Application initialization
│   │   └── app.go
│   ├── config          # Configuration handling
│   │   └── config.go
│   ├── database        # Database connections and migrations
│   │   ├── connection.go
│   │   └── migrations.go
│   ├── handlers        # HTTP request handlers
│   │   ├── accounts.go
│   │   ├── categories.go
│   │   ├── financial_goals.go
│   │   ├── reports.go
│   │   ├── transactions.go
│   │   └── users.go
│   ├── middleware      # Middleware components (JWT, Logging)
│   ├── models          # Data models for the application
│   │   ├── accounts.go
│   │   ├── transactions.go
│   │   ├── financial_goals.go
│   │   ├── currency_rates.go
│   │   └── users.go
│   ├── redis_client    # Redis client for caching
│   │   └── redis_client.go
│   ├── services        # Business logic implementation
│   │   ├── transactions.go
│   │   ├── reports.go
│   │   └── users.go
│   └── transport       # Transport layer (HTTP routes, not fully implemented)
├── migrations          # SQL migration files
│   ├── 001_add_columns_to_accounts.sql
│   ├── 002_add_columns_to_accounts.sql
│   ├── 003_add_columns_to_accounts.sql
│   └── ...
├── go.mod              # Go modules
├── go.sum              # Go module checksums
└── README.md           # Project documentation
```

## Key Features

1. **User Authentication**
   - Secure registration and login with JWT tokens.
   - Token-based middleware to protect API endpoints.

2. **Transaction Management**
   - Add, delete, and view transactions by linking them to accounts and categories.
   - Scheduled transactions support for recurring payments.

3. **Financial Goal Tracking**
   - Create, track, and monitor savings goals.
   - Automatically validate goal progress against targets.

4. **Dynamic Reports**
   - Generate insightful financial reports.
   - Export options for financial data.

5. **Caching with Redis**
   - Accelerates API responses for frequently requested data.
   - Reduces database load with efficient caching strategies.

6. **Database Migrations**
   - SQL-based migrations ensure smooth schema updates.
   - Migration Files: `migrations/001_add_columns_to_accounts.sql` ... `migrations/008_delete_test_table.sql`.

7. **Swagger Integration**
   - Fully documented REST API accessible through Swagger UI.

8. **Clean Architecture**
   - Separation of concerns into layers: Handlers, Services, and Models.

## Configuration File (`configs/config.yaml`)
```yaml
database:
  host: "localhost"
  port: 5432
  user: "finance_user"
  password: "secure_password"
  dbname: "finance_db"
  sslmode: "disable"

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
```

## Getting Started

### Prerequisites
- Go version 1.18+
- Redis server
- PostgreSQL database
- JWT Middleware

### Installation Steps
1. **Clone the Repository**
   ```bash
   git clone https://github.com/abylsliam44/personal_finance_project.git
   cd personal_finance_project/backend
   ```
2. **Install Dependencies**
   ```bash
   go mod tidy
   ```
3. **Configure the Environment**
   Update `configs/config.yaml` to match your local database and Redis configurations.

4. **Run Migrations**
   ```bash
   go run internal/database/migrations.go
   ```
5. **Start the Application**
   ```bash
   go run cmd/main.go
   ```

6. **Access Swagger API**
   Open your browser at:
   ```
   http://localhost:8080/swagger/
   ```

## Contributions
Contributions are welcome! Follow these steps:
1. Fork the repository.
2. Create a new feature branch.
3. Make your changes and write tests.
4. Submit a pull request.

## Future Enhancements
- Add support for expense analysis dashboards.
- Implement notifications for upcoming bill payments.
- Introduce support for multiple currencies.

## Contact
For suggestions, feedback, or issues:
- **Email**: abylajslamzanov@gmail.com
- **GitHub**: [abylay](https://github.com/abylsliam44/)

---
Developed with ❤️ by Abylay Slamzhanov
