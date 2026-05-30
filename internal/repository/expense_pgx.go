package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/omidMorovati/expenseTracker/internal/domain"
	"time"
)

type ExpenseRepo struct {
	pool *pgxpool.Pool
}

func (r *ExpenseRepo) DailyTotals(ctx context.Context, date time.Time) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ExpenseRepo) MonthlyTotals(ctx context.Context, year, month int) ([]domain.MonthTotal, error) {
	//TODO implement me
	panic("implement me")
}

func NewExpenseRepo(pool *pgxpool.Pool) *ExpenseRepo {
	return &ExpenseRepo{pool: pool}
}

func (r *ExpenseRepo) Create(ctx context.Context, e *domain.Expense) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO expenses (id, title, amount, category, date) VALUES ($1, $2, $3, $4, $5)`,
		e.ID, e.Title, e.Amount, e.Category, e.Date,
	)
	return err
}

func (r *ExpenseRepo) ListByDateRange(ctx context.Context, start, end time.Time) ([]domain.Expense, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, title, amount, category, date, created_at FROM expenses WHERE date BETWEEN $1 AND $2 ORDER BY date DESC`,
		start, end,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []domain.Expense
	for rows.Next() {
		var e domain.Expense
		if err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Category, &e.Date, &e.CreatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, rows.Err()
}

// Implement DailyTotals & MonthlyTotals using GROUP BY date/extract
