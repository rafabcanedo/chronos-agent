package transactions

import (
	"net/http"

	"github.com/rafabcanedo/obsidian-golang-toolkit/api/binding"
	"github.com/rafabcanedo/obsidian-golang-toolkit/api/jwtauth"
	"github.com/rafabcanedo/obsidian-golang-toolkit/api/respond"
	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/validator"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type createRequest struct {
	Description string   `json:"description"`
	Amount      float64  `json:"amount"`
	Category    string   `json:"category"`
	DueDate     *string  `json:"due_date"`
	Status      string   `json:"status"`
	PaidAt      *string  `json:"paid_at"`
}

type listParams struct {
	Month    string `form:"month"`
	Status   string `form:"status"`
	Category string `form:"category"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type summaryParams struct {
	Month string `form:"month"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := jwtauth.GetClaims(r)
	if !ok {
		respond.Error(w, errutil.Unauthorized("unauthorized"))
		return
	}

	var req createRequest
	if err := binding.JSON(r, &req); err != nil {
		respond.Error(w, errutil.BadRequest(err.Error()))
		return
	}

	if errs := validator.New().
		Required("description", req.Description).
		MaxLen("description", req.Description, 150).
		Required("category", req.Category).
		Min("amount", req.Amount, 0).
		Validate(); errs != nil {
		respond.Error(w, errutil.BadRequest(errs.Error()))
		return
	}

	tx, err := h.svc.Create(r.Context(), claims.Subject, req.Description, req.Category, req.Amount, req.DueDate, req.PaidAt, req.Status)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.Created(w, tx)
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	claims, ok := jwtauth.GetClaims(r)
	if !ok {
		respond.Error(w, errutil.Unauthorized("unauthorized"))
		return
	}

	var params summaryParams
	if err := binding.Query(r, &params); err != nil {
		respond.Error(w, errutil.BadRequest(err.Error()))
		return
	}

	result, err := h.svc.Summary(r.Context(), claims.Subject, params.Month)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.OK(w, result)
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

	txs, total, err := h.svc.List(r.Context(), claims.Subject, params.Month, params.Status, params.Category, p)
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

	respond.OK(w, map[string]any{"message": "Transaction deleted successfully."})
}
