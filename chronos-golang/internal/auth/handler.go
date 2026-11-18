package auth

import (
	"net/http"

	"github.com/rafabcanedo/obsidian-golang-toolkit/api/binding"
	"github.com/rafabcanedo/obsidian-golang-toolkit/api/jwtauth"
	"github.com/rafabcanedo/obsidian-golang-toolkit/api/respond"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/validator"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
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

	respond.OK(w, map[string]any{"message": "User deleted successfully."})
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

	token, err := h.svc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.OK(w, map[string]any{"token": token})
}
