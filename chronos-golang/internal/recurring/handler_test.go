package recurring_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/mocks"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/recurring"
	itestutil "github.com/rafabcanedo/chronos-agent/chronos-golang/internal/testutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
	"github.com/rafabcanedo/obsidian-golang-toolkit/testutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
)

func TestList_success(t *testing.T) {
	svc := &mocks.MockRecurringService{
		ListFn: func(_ context.Context, _, _, _, _ string, _ pagination.Params) ([]recurring.RecurringBill, int, error) {
			return []recurring.RecurringBill{{ID: "1", Name: "Netflix", DefaultAmount: 45.0, BillingDay: 10}}, 1, nil
		},
	}
	h := recurring.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodGet, "/api/recurring", "")
	rec := itestutil.CallWithAuth(h.List, req)

	testutil.Equal(t, rec.Code, http.StatusOK)
	var body map[string]any
	json.NewDecoder(rec.Body).Decode(&body)
	_, hasData := body["data"]
	_, hasMeta := body["meta"]
	testutil.Equal(t, hasData, true)
	testutil.Equal(t, hasMeta, true)
}

func TestList_noAuth(t *testing.T) {
	h := recurring.NewHandler(&mocks.MockRecurringService{})
	req := httptest.NewRequest(http.MethodGet, "/api/recurring", nil)
	rec := itestutil.CallDirect(h.List, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestCreate_success(t *testing.T) {
	svc := &mocks.MockRecurringService{
		CreateFn: func(_ context.Context, _, name, category string, amount float64, day int) (recurring.RecurringBill, error) {
			return recurring.RecurringBill{ID: "1", Name: name, Category: category, DefaultAmount: amount, BillingDay: day}, nil
		},
	}
	h := recurring.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodPost, "/api/recurring",
		`{"name":"Netflix","category":"leisure","default_amount":45.0,"billing_day":10}`)
	rec := itestutil.CallWithAuth(h.Create, req)
	testutil.Equal(t, rec.Code, http.StatusCreated)
}

func TestCreate_noAuth(t *testing.T) {
	h := recurring.NewHandler(&mocks.MockRecurringService{})
	req := httptest.NewRequest(http.MethodPost, "/api/recurring",
		strings.NewReader(`{"name":"Netflix","category":"leisure","default_amount":45.0,"billing_day":10}`))
	rec := itestutil.CallDirect(h.Create, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestCreate_validationError(t *testing.T) {
	h := recurring.NewHandler(&mocks.MockRecurringService{})
	req := itestutil.AuthedReq(http.MethodPost, "/api/recurring",
		`{"name":"Netflix","category":"leisure","default_amount":45.0,"billing_day":0}`)
	rec := itestutil.CallWithAuth(h.Create, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestUpdate_success(t *testing.T) {
	svc := &mocks.MockRecurringService{
		UpdateFn: func(_ context.Context, _, id string, _ recurring.UpdateFields) (*recurring.RecurringBill, error) {
			return &recurring.RecurringBill{ID: id, Name: "Netflix HD", DefaultAmount: 55.0, BillingDay: 10}, nil
		},
	}
	h := recurring.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodPut, "/api/recurring/1", `{"name":"Netflix HD"}`)
	rec := itestutil.CallWithAuth(h.Update, req)
	testutil.Equal(t, rec.Code, http.StatusOK)
}

func TestUpdate_noAuth(t *testing.T) {
	h := recurring.NewHandler(&mocks.MockRecurringService{})
	req := httptest.NewRequest(http.MethodPut, "/api/recurring/1", strings.NewReader(`{"name":"Netflix HD"}`))
	rec := itestutil.CallDirect(h.Update, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestDelete_success(t *testing.T) {
	svc := &mocks.MockRecurringService{
		DeleteFn: func(_ context.Context, _, _ string) error { return nil },
	}
	h := recurring.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodDelete, "/api/recurring/1", "")
	rec := itestutil.CallWithAuth(h.Delete, req)

	testutil.Equal(t, rec.Code, http.StatusOK)
	var body map[string]any
	json.NewDecoder(rec.Body).Decode(&body)
	data, _ := body["data"].(map[string]any)
	testutil.Equal(t, data["message"], "Recurring bill deleted successfully.")
}

func TestDelete_noAuth(t *testing.T) {
	h := recurring.NewHandler(&mocks.MockRecurringService{})
	req := httptest.NewRequest(http.MethodDelete, "/api/recurring/1", nil)
	rec := itestutil.CallDirect(h.Delete, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestDelete_notFound(t *testing.T) {
	svc := &mocks.MockRecurringService{
		DeleteFn: func(_ context.Context, _, _ string) error {
			return errutil.NotFound("recurring bill not found")
		},
	}
	h := recurring.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodDelete, "/api/recurring/1", "")
	rec := itestutil.CallWithAuth(h.Delete, req)
	testutil.Equal(t, rec.Code, http.StatusNotFound)
}
