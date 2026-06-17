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
