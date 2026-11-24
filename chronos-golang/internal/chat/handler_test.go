package chat_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/chat"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/mocks"
	itestutil "github.com/rafabcanedo/chronos-agent/chronos-golang/internal/testutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/testutil"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
)

func TestHandle_success(t *testing.T) {
	svc := &mocks.MockChatService{
		HandleFn: func(_ context.Context, _, _ string) (string, error) {
			return "Conta de internet marcada como paga.", nil
		},
	}
	h := chat.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodPost, "/api/chat", `{"message":"paguei a internet"}`)
	rec := itestutil.CallWithAuth(h.Handle, req)

	testutil.Equal(t, rec.Code, http.StatusOK)
	var body map[string]any
	json.NewDecoder(rec.Body).Decode(&body)
	data, _ := body["data"].(map[string]any)
	testutil.Equal(t, data["message"], "Conta de internet marcada como paga.")
}

func TestHandle_noAuth(t *testing.T) {
	h := chat.NewHandler(&mocks.MockChatService{})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"message":"paguei a internet"}`))
	rec := itestutil.CallDirect(h.Handle, req)
	testutil.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestHandle_invalidJSON(t *testing.T) {
	h := chat.NewHandler(&mocks.MockChatService{})
	req := itestutil.AuthedReq(http.MethodPost, "/api/chat", `{invalid}`)
	rec := itestutil.CallWithAuth(h.Handle, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestHandle_missingMessage(t *testing.T) {
	h := chat.NewHandler(&mocks.MockChatService{})
	req := itestutil.AuthedReq(http.MethodPost, "/api/chat", `{}`)
	rec := itestutil.CallWithAuth(h.Handle, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestHandle_serviceError(t *testing.T) {
	svc := &mocks.MockChatService{
		HandleFn: func(_ context.Context, _, _ string) (string, error) {
			return "", errutil.BadRequest("Could not understand your request.")
		},
	}
	h := chat.NewHandler(svc)
	req := itestutil.AuthedReq(http.MethodPost, "/api/chat", `{"message":"xyz"}`)
	rec := itestutil.CallWithAuth(h.Handle, req)
	testutil.Equal(t, rec.Code, http.StatusBadRequest)
}
