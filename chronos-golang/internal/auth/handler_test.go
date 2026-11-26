package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/auth"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/mocks"
	itestutil "github.com/rafabcanedo/chronos-agent/chronos-golang/internal/testutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/testutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
)

func TestRegister_success(t *testing.T) {
	svc := &mocks.MockAuthService{
		RegisterFn: func(_ context.Context, username, email, _ string) (*auth.User, error) {
			return &auth.User{ID: "1", Username: username, Email: email, CreatedAt: "2024-01-01T00:00:00Z"}, nil
		},
	}
	h := auth.NewHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/user",
		strings.NewReader(`{"username":"rafa","email":"rafa@email.com","password":"12345678"}`))
	rec := itestutil.CallDirect(h.Register, req)

	testutil.Equal(t, rec.Code, http.StatusCreated)
	var body map[string]any
	json.NewDecoder(rec.Body).Decode(&body)
	_, hasData := body["data"]
	testutil.Equal(t, hasData, true)
}

func TestRegister_invalidJSON(t *testing.T) {
	h := auth.NewHandler(&mocks.MockAuthService{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/user", strings.NewReader(`{invalid}`))
	rec := itestutil.CallDirect(h.Register, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestRegister_missingField(t *testing.T) {
	h := auth.NewHandler(&mocks.MockAuthService{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/user",
		strings.NewReader(`{"email":"rafa@email.com","password":"12345678"}`))
	rec := itestutil.CallDirect(h.Register, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestLogin_success(t *testing.T) {
	svc := &mocks.MockAuthService{
		LoginFn: func(_ context.Context, _, _ string) (string, string, error) {
			return "jwt-token-here", "refresh-token-here", nil
		},
	}
	h := auth.NewHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"rafa","password":"12345678"}`))
	rec := itestutil.CallDirect(h.Login, req)

	testutil.Equal(t, rec.Code, http.StatusOK)
	var body map[string]any
	json.NewDecoder(rec.Body).Decode(&body)
	data, _ := body["data"].(map[string]any)
	_, hasToken := data["token"]
	testutil.Equal(t, hasToken, true)
}

func TestLogin_invalidJSON(t *testing.T) {
	h := auth.NewHandler(&mocks.MockAuthService{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{invalid}`))
	rec := itestutil.CallDirect(h.Login, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestLogin_missingField(t *testing.T) {
	h := auth.NewHandler(&mocks.MockAuthService{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"rafa"}`))
	rec := itestutil.CallDirect(h.Login, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestLogin_unauthorized(t *testing.T) {
	svc := &mocks.MockAuthService{
		LoginFn: func(_ context.Context, _, _ string) (string, string, error) {
			return "", "", errutil.Unauthorized("invalid credentials")
		},
	}
	h := auth.NewHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"rafa","password":"wrong"}`))
	rec := itestutil.CallDirect(h.Login, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestUpdate_success(t *testing.T) {
	svc := &mocks.MockAuthService{
		UpdateFn: func(_ context.Context, userID string, _ auth.UpdateFields) (*auth.User, error) {
			return &auth.User{ID: userID, Username: "newname", Email: "rafa@email.com", CreatedAt: "2024-01-01T00:00:00Z"}, nil
		},
	}
	h := auth.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodPatch, "/api/auth/user", `{"username":"newname"}`)
	rec := itestutil.CallWithAuth(h.Update, req)
	testutil.Equal(t, rec.Code, http.StatusOK)
}

func TestUpdate_noAuth(t *testing.T) {
	h := auth.NewHandler(&mocks.MockAuthService{})
	req := httptest.NewRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"username":"newname"}`))
	rec := itestutil.CallDirect(h.Update, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestUpdate_invalidJSON(t *testing.T) {
	h := auth.NewHandler(&mocks.MockAuthService{})
	req := itestutil.AuthedReq(http.MethodPatch, "/api/auth/user", `{invalid}`)
	rec := itestutil.CallWithAuth(h.Update, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestDelete_success(t *testing.T) {
	svc := &mocks.MockAuthService{
		DeleteFn: func(_ context.Context, _ string) error { return nil },
	}
	h := auth.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodDelete, "/api/auth/user", "")
	rec := itestutil.CallWithAuth(h.Delete, req)

	testutil.Equal(t, rec.Code, http.StatusOK)
	var body map[string]any
	json.NewDecoder(rec.Body).Decode(&body)
	data, _ := body["data"].(map[string]any)
	testutil.Equal(t, data["message"], "User deleted successfully.")
}

func TestDelete_noAuth(t *testing.T) {
	h := auth.NewHandler(&mocks.MockAuthService{})
	req := httptest.NewRequest(http.MethodDelete, "/api/auth/user", nil)
	rec := itestutil.CallDirect(h.Delete, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestMe_success(t *testing.T) {
	svc := &mocks.MockAuthService{
		MeFn: func(_ context.Context, userID string) (*auth.User, error) {
			return &auth.User{ID: userID, Username: "rafa", Email: "rafa@email.com", CreatedAt: "2024-01-01T00:00:00Z"}, nil
		},
	}
	h := auth.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodGet, "/api/auth/user", "")
	rec := itestutil.CallWithAuth(h.Me, req)
	testutil.Equal(t, rec.Code, http.StatusOK)
	var body map[string]any
	json.NewDecoder(rec.Body).Decode(&body)
	data, _ := body["data"].(map[string]any)
	testutil.Equal(t, data["username"], "rafa")
}

func TestMe_noAuth(t *testing.T) {
	h := auth.NewHandler(&mocks.MockAuthService{})
	req := httptest.NewRequest(http.MethodGet, "/api/auth/user", nil)
	rec := itestutil.CallDirect(h.Me, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestMe_notFound(t *testing.T) {
	svc := &mocks.MockAuthService{
		MeFn: func(_ context.Context, _ string) (*auth.User, error) {
			return nil, errutil.NotFound("user not found")
		},
	}
	h := auth.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodGet, "/api/auth/user", "")
	rec := itestutil.CallWithAuth(h.Me, req)
	testutil.Equal(t, rec.Code, http.StatusNotFound)
}

func TestDelete_notFound(t *testing.T) {
	svc := &mocks.MockAuthService{
		DeleteFn: func(_ context.Context, _ string) error {
			return errutil.NotFound("user not found")
		},
	}
	h := auth.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodDelete, "/api/auth/user", "")
	rec := itestutil.CallWithAuth(h.Delete, req)
	testutil.Equal(t, rec.Code, http.StatusNotFound)
}
