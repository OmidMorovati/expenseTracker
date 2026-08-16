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
		`INSERT INTO expenses (id, user_id, title, amount, category, date) VALUES ($1, $2, $3, $4, $5, $6)`,
		e.ID, e.UserID, e.Title, e.Amount, e.Category, e.Date,
	)
	return err
}

func (r *ExpenseRepo) ListByDateRange(ctx context.Context, userID string, start, end time.Time) ([]domain.Expense, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, title, amount, category, date, created_at 
         FROM expenses 
         WHERE user_id = $1 AND date BETWEEN $2 AND $3 
         ORDER BY date DESC`,
		userID, start, end,
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

func (r *ExpenseRepo) ListRecent(ctx context.Context, userID string, limit int) ([]domain.Expense, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, title, amount, category, date, created_at 
         FROM expenses 
         WHERE user_id = $1 
         ORDER BY date DESC, created_at DESC 
         LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []domain.Expense
	for rows.Next() {
		var e domain.Expense
		if err := rows.Scan(&e.ID, &e.UserID, &e.Title, &e.Amount, &e.Category, &e.Date, &e.CreatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, rows.Err()
}

func (r *ExpenseRepo) GetTotalByDateRange(ctx context.Context, userID string, start, end time.Time) (float64, error) {
	var total float64
	// COALESCE ensures we return 0 instead of NULL if no expenses exist
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM expenses 
         WHERE user_id = $1 AND date >= $2 AND date < $3`,
		userID, start, end,
	).Scan(&total)

	if err != nil {
		return 0, err
	}
	return total, nil
}
