# 💰 Expense Tracker - Go + HTMX

A modern, production-ready expense tracking application built with Go, PostgreSQL, and HTMX. Features JWT authentication, server-rendered HTML, and dynamic UI updates without heavy JavaScript frameworks.

![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## ✨ Features

- 🔐 **JWT Authentication** - Secure login/register with bcrypt password hashing
- 📊 **Interactive Dashboard** - Real-time expense tracking with HTMX-powered updates
- 📅 **Date Range Filtering** - View expenses by Today, This Week, or This Month
- 💳 **Expense Management** - Create, view, and track daily expenses by category
- 📱 **Responsive Design** - Mobile-friendly UI with modern CSS
- 🔒 **Multi-Tenant Security** - User isolation with UUID-based foreign keys
- 🚀 **Production-Ready** - Graceful shutdown, structured logging, connection pooling

## 🛠 Tech Stack

### Backend
- **Go 1.21+** - Clean architecture with domain-driven design
- **Chi Router** - Lightweight, idiomatic HTTP router
- **pgx** - High-performance PostgreSQL driver
- **Goose** - Database migrations
- **golang-jwt/jwt** - JWT authentication
- **bcrypt** - Secure password hashing
- **slog** - Structured logging (stdlib)

### Frontend
- **HTMX 1.9** - Dynamic HTML updates without JavaScript frameworks
- **Go Templates** - Server-side rendering with auto-escaping
- **Vanilla CSS** - Modern, responsive design (no frameworks)

### Infrastructure
- **PostgreSQL 14+** - Relational database with UUID support
- **Docker** - Containerized development environment

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Docker (optional, for easy PostgreSQL setup)

## 🚀 Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/expense-tracker.git
cd expense-tracker
```

### 2. Start PostgreSQL (Docker)
```bash
docker run -d \
  --name expense-tracker-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=expenses \
  -p 5432:5432 \
  postgres:16
```

### 3. Configure Environment
Create a .env file in the project root:
```bash
PORT=:8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=expenses
DB_SSLMODE=disable
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRATION=24h
```

### 4. Install Dependencies
```bash
go mod download
```

### 5. Run Database Migrations
```bash
go run github.com/pressly/goose/v3/cmd/goose@latest -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/expenses?sslmode=disable" up
```

### 6. Start the Server
```bash
go run cmd/server/main.go
```

## 📁 Project Structure

```azure
expense-tracker/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Environment configuration
│   ├── domain/
│   │   ├── auth.go              # User domain & interfaces
│   │   └── expense.go           # Expense domain & interfaces
│   ├── handler/
│   │   ├── auth.go              # Authentication handlers
│   │   └── expense.go           # Expense handlers
│   ├── middleware/
│   │   └── auth.go              # JWT authentication middleware
│   ├── repository/
│   │   ├── expense_pgx.go       # PostgreSQL expense repository
│   │   └── user_pgx.go          # PostgreSQL user repository
│   └── service/
│       ├── auth.go              # Authentication business logic
│       └── expense.go           # Expense business logic
├── migrations/
│   └── *.sql                    # Database migrations
├── web/
│   ├── static/
│   │   └── style.css            # Application styles
│   └── templates/
│       ├── dashboard.html       # Main dashboard view
│       ├── create.html          # Expense creation form
│       ├── login.html           # Login page
│       ├── expenses_table_rows.html
│       └── stats_cards.html
├── .env                         # Environment variables (not in git)
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 🔑 API Endpoints

### Authentication
- `POST /auth/register` - Register new user (no auth)
- `POST /auth/login` - Login and get JWT token (no auth)
- `GET /login` - Login page (no auth)

### Expenses
- `POST /expenses` - Create new expense (auth required)
- `GET /expenses/new` - Expense creation form (auth required)
- `GET /dashboard` - Main dashboard (auth required)
- `GET /api/expenses/recent?limit=10` - Get recent expenses (auth required)
- `GET /api/expenses/stats?period=today` - Get expense totals by period (auth required)

### Query Parameters
- `/api/expenses/recent` accepts `limit` parameter (5, 10, 50)
- `/api/expenses/stats` accepts `period` parameter (today, week, month)

## 🔧 Configuration

### Environment Variables

**PORT**
- Default: `:8080`
- Description: Server port

**DB_HOST**
- Default: `localhost`
- Description: PostgreSQL host

**DB_PORT**
- Default: `5432`
- Description: PostgreSQL port

**DB_USER**
- Default: `postgres`
- Description: Database user

**DB_PASSWORD**
- Default: `postgres`
- Description: Database password

**DB_NAME**
- Default: `expenses`
- Description: Database name

**DB_SSLMODE**
- Default: `disable`
- Description: SSL mode (disable/require)

**JWT_SECRET**
- Default: (required)
- Description: Secret key for JWT signing

**JWT_EXPIRATION**
- Default: `24h`
- Description: Token expiration duration

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.
