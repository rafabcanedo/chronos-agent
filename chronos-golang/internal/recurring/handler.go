package recurring

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
	svc Servicer
}

func NewHandler(svc Servicer) *Handler {
	return &Handler{svc: svc}
}

type listParams struct {
	Active   string `form:"active"`
	Category string `form:"category"`
	Name     string `form:"name"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type createRequest struct {
	Name          string  `json:"name"`
	DefaultAmount float64 `json:"default_amount"`
	BillingDay    int     `json:"billing_day"`
	Category      string  `json:"category"`
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

	bills, total, err := h.svc.List(r.Context(), claims.Subject, params.Active, params.Category, params.Name, p)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.Paginated(w, bills, p.NewMeta(total))
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
		Required("name", req.Name).
		MaxLen("name", req.Name, 100).
		Required("category", req.Category).
		Min("default_amount", req.DefaultAmount, 0).
		Between("billing_day", float64(req.BillingDay), 1, 31).
		Validate(); errs != nil {
		respond.Error(w, errutil.BadRequest(errs.Error()))
		return
	}

	bill, err := h.svc.Create(r.Context(), claims.Subject, req.Name, req.Category, req.DefaultAmount, req.BillingDay)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.Created(w, bill)
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

	respond.OK(w, map[string]any{"message": "Recurring bill deleted successfully."})
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

	bill, err := h.svc.Update(r.Context(), claims.Subject, id, fields)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.OK(w, bill)
}
