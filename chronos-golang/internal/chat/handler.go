package chat

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

type chatRequest struct {
	Message string `json:"message"`
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := jwtauth.GetClaims(r)
	if !ok {
		respond.Error(w, errutil.Unauthorized("unauthorized"))
		return
	}

	var req chatRequest
	if err := binding.JSON(r, &req); err != nil {
		respond.Error(w, errutil.BadRequest(err.Error()))
		return
	}

	if errs := validator.New().Required("message", req.Message).Validate(); errs != nil {
		respond.Error(w, errutil.BadRequest(errs.Error()))
		return
	}

	msg, err := h.svc.Handle(r.Context(), claims.Subject, req.Message)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.OK(w, map[string]any{"message": msg})
}
