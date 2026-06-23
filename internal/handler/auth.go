package handler

import (
	"encoding/json"
	"github.com/omidMorovati/expenseTracker/internal/service"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

type AuthHandler struct {
	svc       *service.AuthService
	logger    *slog.Logger
	templates *template.Template
}

func NewAuthHandler(svc *service.AuthService, logger *slog.Logger, templates *template.Template) *AuthHandler {
	return &AuthHandler{svc: svc, logger: logger, templates: templates}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.svc.Register(r.Context(), req.Email, req.Password); err != nil {
		h.logger.Error("register failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to register user"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "user created successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	token, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("login failed", "error", err)
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		return
	}

	// Set secure cookie for page navigation
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(24 * time.Hour.Seconds()),
		HttpOnly: true,  // Prevent XSS theft
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "login", nil)
}

// writeJSON is a small helper to keep handlers DRY and consistent
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
