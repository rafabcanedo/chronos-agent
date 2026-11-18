package testutil

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/rafabcanedo/obsidian-golang-toolkit/api/jwtauth"
)

const TestSecret = "test-secret"
const TestUserID = "user-123"

func AuthedReq(method, path, body string) *http.Request {
	token, _ := jwtauth.Generate(TestSecret, jwtauth.Claims{
		Subject: TestUserID,
		Expiry:  time.Now().Add(time.Hour),
	})
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func CallWithAuth(h http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	jwtauth.Middleware(TestSecret)(h).ServeHTTP(rec, req)
	return rec
}

func CallDirect(h http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}
