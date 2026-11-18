package transactions_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/mocks"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/transactions"
	itestutil "github.com/rafabcanedo/chronos-agent/chronos-golang/internal/testutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
	"github.com/rafabcanedo/obsidian-golang-toolkit/testutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
)

func TestCreate_success(t *testing.T) {
	svc := &mocks.MockTransactionsService{
		CreateFn: func(_ context.Context, _, description, category string, _ float64, _, _ *string, _ string) (*transactions.Transaction, error) {
			return &transactions.Transaction{ID: "1", Description: description, Category: category, Amount: 50.0, Status: transactions.StatusPending}, nil
		},
	}
	h := transactions.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodPost, "/api/transactions",
		`{"description":"Almoço","category":"restaurants","amount":50.0}`)
	rec := itestutil.CallWithAuth(h.Create, req)
	testutil.Equal(t, rec.Code, http.StatusCreated)
}

func TestCreate_noAuth(t *testing.T) {
	h := transactions.NewHandler(&mocks.MockTransactionsService{})
	req := httptest.NewRequest(http.MethodPost, "/api/transactions",
		strings.NewReader(`{"description":"Almoço","category":"restaurants","amount":50.0}`))
	rec := itestutil.CallDirect(h.Create, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestCreate_invalidJSON(t *testing.T) {
	h := transactions.NewHandler(&mocks.MockTransactionsService{})
	req := itestutil.AuthedReq(http.MethodPost, "/api/transactions", `{invalid}`)
	rec := itestutil.CallWithAuth(h.Create, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestCreate_validationError(t *testing.T) {
	h := transactions.NewHandler(&mocks.MockTransactionsService{})
	req := itestutil.AuthedReq(http.MethodPost, "/api/transactions",
		`{"category":"restaurants","amount":50.0}`)
	rec := itestutil.CallWithAuth(h.Create, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestList_success(t *testing.T) {
	svc := &mocks.MockTransactionsService{
		ListFn: func(_ context.Context, _, _, _, _ string, _ pagination.Params) ([]transactions.Transaction, int, error) {
			return []transactions.Transaction{{ID: "1", Description: "Almoço", Amount: 50.0}}, 1, nil
		},
	}
	h := transactions.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodGet, "/api/transactions", "")
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
	h := transactions.NewHandler(&mocks.MockTransactionsService{})
	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	rec := itestutil.CallDirect(h.List, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestSummary_success(t *testing.T) {
	svc := &mocks.MockTransactionsService{
		SummaryFn: func(_ context.Context, _, _ string) (map[string]transactions.StatusSummary, error) {
			return map[string]transactions.StatusSummary{"pending": {Count: 1, Total: 50.0}}, nil
		},
	}
	h := transactions.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodGet, "/api/transactions/summary", "")
	rec := itestutil.CallWithAuth(h.Summary, req)
	testutil.Equal(t, rec.Code, http.StatusOK)
}

func TestSummary_noAuth(t *testing.T) {
	h := transactions.NewHandler(&mocks.MockTransactionsService{})
	req := httptest.NewRequest(http.MethodGet, "/api/transactions/summary", nil)
	rec := itestutil.CallDirect(h.Summary, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestUpdate_success(t *testing.T) {
	svc := &mocks.MockTransactionsService{
		UpdateFn: func(_ context.Context, _, id string, _ transactions.UpdateFields) (*transactions.Transaction, error) {
			return &transactions.Transaction{ID: id, Description: "Updated", Amount: 100.0, Status: transactions.StatusPaid}, nil
		},
	}
	h := transactions.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodPatch, "/api/transactions/1", `{"status":"paid"}`)
	rec := itestutil.CallWithAuth(h.Update, req)
	testutil.Equal(t, rec.Code, http.StatusOK)
}

func TestUpdate_noAuth(t *testing.T) {
	h := transactions.NewHandler(&mocks.MockTransactionsService{})
	req := httptest.NewRequest(http.MethodPatch, "/api/transactions/1", strings.NewReader(`{"status":"paid"}`))
	rec := itestutil.CallDirect(h.Update, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestDelete_success(t *testing.T) {
	svc := &mocks.MockTransactionsService{
		DeleteFn: func(_ context.Context, _, _ string) error { return nil },
	}
	h := transactions.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodDelete, "/api/transactions/1", "")
	rec := itestutil.CallWithAuth(h.Delete, req)

	testutil.Equal(t, rec.Code, http.StatusOK)
	var body map[string]any
	json.NewDecoder(rec.Body).Decode(&body)
	data, _ := body["data"].(map[string]any)
	testutil.Equal(t, data["message"], "Transaction deleted successfully.")
}

func TestDelete_noAuth(t *testing.T) {
	h := transactions.NewHandler(&mocks.MockTransactionsService{})
	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/1", nil)
	rec := itestutil.CallDirect(h.Delete, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestDelete_notFound(t *testing.T) {
	svc := &mocks.MockTransactionsService{
		DeleteFn: func(_ context.Context, _, _ string) error {
			return errutil.NotFound("transaction not found")
		},
	}
	h := transactions.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodDelete, "/api/transactions/1", "")
	rec := itestutil.CallWithAuth(h.Delete, req)
	testutil.Equal(t, rec.Code, http.StatusNotFound)
}
