package handler

import (
	"encoding/json"
	"github.com/omidMorovati/expenseTracker/internal/domain"
	"github.com/omidMorovati/expenseTracker/internal/service"
	"log/slog"
	"net/http"
	"time"
)

type ExpenseHandler struct {
	svc    *service.ExpenseService
	logger *slog.Logger
}

func NewExpenseHandler(svc *service.ExpenseService, logger *slog.Logger) *ExpenseHandler {
	return &ExpenseHandler{svc: svc, logger: logger}
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string  `json:"title"`
		Amount   float64 `json:"amount"`
		Category string  `json:"category"`
		Date     string  `json:"date"` // YYYY-MM-DD
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	exp := &domain.Expense{
		Title:    req.Title,
		Amount:   req.Amount,
		Category: req.Category,
		Date:     parsedDate,
	}

	if err := h.svc.Create(r.Context(), exp); err != nil {
		h.logger.Error("create expense", "error", err)
		http.Error(w, "failed to create expense", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ExpenseHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Render Go template + fetch reports
}
