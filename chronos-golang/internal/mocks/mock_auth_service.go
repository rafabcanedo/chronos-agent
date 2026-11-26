package mocks

import (
	"context"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/auth"
)

type MockAuthService struct {
	RegisterFn func(ctx context.Context, username, email, password string) (*auth.User, error)
	LoginFn    func(ctx context.Context, username, password string) (string, string, error)
	RefreshFn  func(ctx context.Context, refreshToken string) (string, string, error)
	LogoutFn   func(ctx context.Context, refreshToken string) error
	UpdateFn   func(ctx context.Context, userID string, fields auth.UpdateFields) (*auth.User, error)
	DeleteFn   func(ctx context.Context, userID string) error
	MeFn       func(ctx context.Context, userID string) (*auth.User, error)
}

func (m *MockAuthService) Register(ctx context.Context, username, email, password string) (*auth.User, error) {
	return m.RegisterFn(ctx, username, email, password)
}

func (m *MockAuthService) Login(ctx context.Context, username, password string) (string, string, error) {
	return m.LoginFn(ctx, username, password)
}

func (m *MockAuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	return m.RefreshFn(ctx, refreshToken)
}

func (m *MockAuthService) Logout(ctx context.Context, refreshToken string) error {
	return m.LogoutFn(ctx, refreshToken)
}

func (m *MockAuthService) Update(ctx context.Context, userID string, fields auth.UpdateFields) (*auth.User, error) {
	return m.UpdateFn(ctx, userID, fields)
}

func (m *MockAuthService) Delete(ctx context.Context, userID string) error {
	return m.DeleteFn(ctx, userID)
}

func (m *MockAuthService) Me(ctx context.Context, userID string) (*auth.User, error) {
	return m.MeFn(ctx, userID)
}
