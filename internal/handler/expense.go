package handler

import (
	"encoding/json"
	"fmt"
	"github.com/omidMorovati/expenseTracker/internal/domain"
	"github.com/omidMorovati/expenseTracker/internal/service"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ExpenseHandler struct {
	svc       *service.ExpenseService
	logger    *slog.Logger
	templates *template.Template
}

func NewExpenseHandler(svc *service.ExpenseService, logger *slog.Logger, templates *template.Template) *ExpenseHandler {
	return &ExpenseHandler{svc: svc, logger: logger, templates: templates}
}

func (h *ExpenseHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "dashboard", nil)
	if err != nil {
		return
	}
}

func (h *ExpenseHandler) CreatePage(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "create", nil)
	if err != nil {
		return
	}
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var title, category, dateStr string
	var amount float64

	// Handle both HTMX form submission & JSON API
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		title = r.FormValue("title")
		amount, _ = strconv.ParseFloat(r.FormValue("amount"), 64)
		category = r.FormValue("category")
		dateStr = r.FormValue("date")
	} else {
		var req struct {
			Title, Category, Date string
			Amount                float64
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.respond(w, http.StatusBadRequest, "invalid json")
			return
		}
		title, amount, category, dateStr = req.Title, req.Amount, req.Category, req.Date
	}

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.respond(w, http.StatusBadRequest, "invalid date format (YYYY-MM-DD)")
		return
	}

	exp := &domain.Expense{
		Title: title, Amount: amount, Category: category, Date: parsedDate,
	}

	if err := h.svc.Create(r.Context(), exp); err != nil {
		if r.Header.Get("HX-Request") == "true" {
			_, err := fmt.Fprintf(w, `<div class="error">❌ Failed: %v</div>`, err)
			if err != nil {
				return
			}
		} else {
			h.respond(w, http.StatusInternalServerError, "failed to create expense")
		}
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		_, err := fmt.Fprint(w, `<div class="success">✅ Expense added successfully!</div>`)
		if err != nil {
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *ExpenseHandler) respond(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(map[string]string{"error": msg})
	if err != nil {
		return
	}
}

func (h *ExpenseHandler) Recent(w http.ResponseWriter, r *http.Request) {
	// Default limit, allow override via query param
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	h.logger.Info(r.URL.Query().Encode())
	expenses, err := h.svc.ListRecent(r.Context(), limit)
	if err != nil {
		h.logger.Error("list recent expenses", "error", err)
		// For HTMX: return empty rows with error message
		fmt.Fprint(w, `<tr><td colspan="4" class="error">Failed to load expenses</td></tr>`)
		return
	}

	// Render ONLY the table rows partial
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "expenses_table_rows", expenses); err != nil {
		h.logger.Error("render template", "error", err)
	}
}

func (h *ExpenseHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "today" // Default fallback
	}

	total, err := h.svc.GetTotalForPeriod(r.Context(), period)
	if err != nil {
		h.logger.Error("get stats failed", "error", err)
		http.Error(w, "failed to load stats", http.StatusInternalServerError)
		return
	}

	// Pass data to the template
	data := struct {
		Total  float64
		Period string
	}{
		Total:  total,
		Period: period,
	}

	// Render ONLY the stats fragment
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "stats_cards", data); err != nil {
		h.logger.Error("render stats template", "error", err)
	}
}
