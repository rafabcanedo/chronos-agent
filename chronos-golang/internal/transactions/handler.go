package transactions

import (
	"net/http"

	"github.com/rafabcanedo/obsidian-golang-toolkit/api/binding"
	"github.com/rafabcanedo/obsidian-golang-toolkit/api/jwtauth"
	"github.com/rafabcanedo/obsidian-golang-toolkit/api/respond"
	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type listParams struct {
	Month    string `form:"month"`
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := jwtauth.GetClaims(r)
	if !ok {
		respond.Error(w, errutil.Unauthorized("unauthorized"))
		return
	}

	var params listParams
	if err := binding.Query(r, &params); err != nil {
		respond.Error(w, errutil.BadRequest(err.Error()))
		return
	}

	p := pagination.New(params.Page, params.PageSize)

	txs, total, err := h.svc.List(r.Context(), claims.Subject, params.Month, params.Status, p)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.Paginated(w, txs, p.NewMeta(total))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := jwtauth.GetClaims(r)
	if !ok {
		respond.Error(w, errutil.Unauthorized("unauthorized"))
		return
	}

	id := r.PathValue("id")

	var fields UpdateFields
	if err := binding.JSON(r, &fields); err != nil {
		respond.Error(w, errutil.BadRequest(err.Error()))
		return
	}

	tx, err := h.svc.Update(r.Context(), claims.Subject, id, fields)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.OK(w, tx)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := jwtauth.GetClaims(r)
	if !ok {
		respond.Error(w, errutil.Unauthorized("unauthorized"))
		return
	}

	id := r.PathValue("id")

	if err := h.svc.Delete(r.Context(), claims.Subject, id); err != nil {
		respond.Error(w, err)
		return
	}

	respond.OK(w, map[string]any{"message": "Transação removida com sucesso."})
}
