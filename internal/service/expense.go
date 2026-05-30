package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/omidMorovati/expenseTracker/internal/domain"
	"time"
)

type ExpenseService struct {
	repo domain.ExpenseRepository
	val  *validator.Validate
}

func NewExpenseService(repo domain.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo, val: validator.New()}
}

func (s *ExpenseService) Create(ctx context.Context, e *domain.Expense) error {
	if err := s.val.Struct(e); err != nil {
		return err
	}
	e.ID = uuid.New().String()
	return s.repo.Create(ctx, e)
}

func (s *ExpenseService) GetDailyReport(ctx context.Context, date time.Time) (float64, error) {
	return s.repo.DailyTotals(ctx, date)
}
