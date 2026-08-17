package service

import (
	"context"
	"fmt"
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
	userID, ok := ctx.Value(domain.UserIDKey).(string)
	if !ok {
		return fmt.Errorf("unauthorized: missing user ID in context")
	}
	e.UserID = userID
	return s.repo.Create(ctx, e)
}

func (s *ExpenseService) GetDailyReport(ctx context.Context, date time.Time) (float64, error) {
	return s.repo.DailyTotals(ctx, date)
}

func (s *ExpenseService) ListRecent(ctx context.Context, limit int) ([]domain.Expense, error) {
	userID, ok := ctx.Value(domain.UserIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("missing user ID in context")
	}
	return s.repo.ListRecent(ctx, userID, limit)
}
func (s *ExpenseService) GetTotalForPeriod(ctx context.Context, period string) (float64, error) {
	userID, ok := ctx.Value(domain.UserIDKey).(string)
	if !ok {
		return 0, fmt.Errorf("missing user ID in context")
	}

	now := time.Now()
	var start, end time.Time

	// Calculate start and end boundaries based on the period
	switch period {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1) // Tomorrow at 00:00
	case "week":
		// Calculate start of the week (Monday)
		offset := int(now.Weekday())
		if offset == 0 {
			offset = 7
		} // Go's Sunday is 0, we want Monday to be 0
		start = time.Date(now.Year(), now.Month(), now.Day()-offset+1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 7) // Next Monday at 00:00
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0) // 1st of next month at 00:00
	default:
		return 0, fmt.Errorf("invalid period: %s", period)
	}

	return s.repo.GetTotalByDateRange(ctx, userID, start, end)
}

func (s *ExpenseService) GetExpense(ctx context.Context, expenseId string) (domain.Expense, error) {
	userID, ok := ctx.Value(domain.UserIDKey).(string)
	if !ok {
		return domain.Expense{}, fmt.Errorf("missing user ID in context")
	}
	return s.repo.GetExpense(ctx, userID, expenseId)
}
