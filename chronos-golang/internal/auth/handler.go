package auth

import (
	"net/http"
	"os"

	"github.com/rafabcanedo/obsidian-golang-toolkit/api/binding"
	"github.com/rafabcanedo/obsidian-golang-toolkit/api/jwtauth"
	"github.com/rafabcanedo/obsidian-golang-toolkit/api/respond"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/validator"
)

const refreshTokenCookie = "refresh_token"

type Handler struct {
	svc Servicer
}

func NewHandler(svc Servicer) *Handler {
	return &Handler{svc: svc}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := binding.JSON(r, &req); err != nil {
		respond.Error(w, errutil.BadRequest(err.Error()))
		return
	}

	if errs := validator.New().
		Required("username", req.Username).
		MaxLen("username", req.Username, 50).
		Required("email", req.Email).
		Email("email", req.Email).
		Required("password", req.Password).
		MinLen("password", req.Password, 8).
		Validate(); errs != nil {
		respond.Error(w, errutil.BadRequest(errs.Error()))
		return
	}

	user, err := h.svc.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.Created(w, map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := binding.JSON(r, &req); err != nil {
		respond.Error(w, errutil.BadRequest(err.Error()))
		return
	}

	if errs := validator.New().
		Required("username", req.Username).
		Required("password", req.Password).
		Validate(); errs != nil {
		respond.Error(w, errutil.BadRequest(errs.Error()))
		return
	}

	accessToken, refreshToken, err := h.svc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		respond.Error(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    refreshToken,
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	respond.OK(w, map[string]any{"token": accessToken})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshTokenCookie)
	if err != nil {
		respond.Error(w, errutil.Unauthorized("missing refresh token"))
		return
	}

	newAccessToken, newRefreshToken, err := h.svc.Refresh(r.Context(), cookie.Value)
	if err != nil {
		respond.Error(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    newRefreshToken,
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	respond.OK(w, map[string]any{"token": newAccessToken})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshTokenCookie)
	if err == nil {
		h.svc.Logout(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookie,
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		MaxAge:   -1,
	})

	respond.OK(w, map[string]any{"message": "logged out successfully"})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := jwtauth.GetClaims(r)
	if !ok {
		respond.Error(w, errutil.Unauthorized("unauthorized"))
		return
	}

	var fields UpdateFields
	if err := binding.JSON(r, &fields); err != nil {
		respond.Error(w, errutil.BadRequest(err.Error()))
		return
	}

	user, err := h.svc.Update(r.Context(), claims.Subject, fields)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.OK(w, map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := jwtauth.GetClaims(r)
	if !ok {
		respond.Error(w, errutil.Unauthorized("unauthorized"))
		return
	}

	if err := h.svc.Delete(r.Context(), claims.Subject); err != nil {
		respond.Error(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookie,
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		MaxAge:   -1,
	})

	respond.OK(w, map[string]any{"message": "User deleted successfully."})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := jwtauth.GetClaims(r)
	if !ok {
		respond.Error(w, errutil.Unauthorized("unauthorized"))
		return
	}

	user, err := h.svc.Me(r.Context(), claims.Subject)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.OK(w, map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}
