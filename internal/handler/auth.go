package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LinerGit/go-auth/internal/handler/dto"
	"github.com/LinerGit/go-auth/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, ErrorResponse{Error: message})
}

// Register godoc
//
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "register"
// @Success 201 {object} dto.AuthResponse
// @Failure 400
// @Failure 409
// @Failure 500
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.RegisterRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request")
		return
	}

	if req.Username == "" || req.Password == "" {
		writeError(w, r, http.StatusBadRequest, "username and password are required")
		return
	}

	tokens, err := h.auth.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserAlreadyExists):
			writeError(w, r, http.StatusConflict, err.Error())
		default:
			writeError(w, r, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Login godoc
//
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "login"
// @Success 200 {object} dto.AuthResponse
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		writeError(w, r, http.StatusBadRequest, "username and password are required")
		return
	}

	tokens, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			writeError(w, r, http.StatusUnauthorized, "invalid username or password")
		default:
			writeError(w, r, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Refresh godoc
//
// @Summary Refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "refresh"
// @Success 200 {object} dto.AuthResponse
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, r, http.StatusBadRequest, "refresh token is required")
		return
	}

	tokens, err := h.auth.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidToken):
			writeError(w, r, http.StatusUnauthorized, "invalid refresh token")
		default:
			writeError(w, r, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Logout godoc
//
// @Summary Logout user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutRequest true "logout"
// @Success 204
// @Failure 400
// @Failure 500
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.LogoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, r, http.StatusBadRequest, "refresh token is required")
		return
	}

	if err := h.auth.Logout(r.Context(), req.RefreshToken); err != nil {
		writeError(w, r, http.StatusInternalServerError, "failed to logout")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
		r.Post("/logout", h.Logout)
	})
}
