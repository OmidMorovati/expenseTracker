package domain

import (
	"context"
	"time"
)

type Expense struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title" validate:"required,max=100"`
	Amount    float64   `json:"amount" validate:"required,gt=0"`
	Category  string    `json:"category" validate:"required,oneof=food transport entertainment bills other"`
	Date      time.Time `json:"date" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type ExpenseRepository interface {
	Create(ctx context.Context, e *Expense) error
	ListByDateRange(ctx context.Context, userID string, start, end time.Time) ([]Expense, error)
	ListRecent(ctx context.Context, userID string, limit int) ([]Expense, error)
	DailyTotals(ctx context.Context, date time.Time) (float64, error)
	MonthlyTotals(ctx context.Context, year, month int) ([]MonthTotal, error)
}

type MonthTotal struct {
	Day   int     `json:"day"`
	Total float64 `json:"total"`
}
