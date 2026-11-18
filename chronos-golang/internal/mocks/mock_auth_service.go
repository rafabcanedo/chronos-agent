package mocks

import (
	"context"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/auth"
)

type MockAuthService struct {
	RegisterFn func(ctx context.Context, username, email, password string) (*auth.User, error)
	LoginFn    func(ctx context.Context, username, password string) (string, error)
	UpdateFn   func(ctx context.Context, userID string, fields auth.UpdateFields) (*auth.User, error)
	DeleteFn   func(ctx context.Context, userID string) error
}

func (m *MockAuthService) Register(ctx context.Context, username, email, password string) (*auth.User, error) {
	return m.RegisterFn(ctx, username, email, password)
}

func (m *MockAuthService) Login(ctx context.Context, username, password string) (string, error) {
	return m.LoginFn(ctx, username, password)
}

func (m *MockAuthService) Update(ctx context.Context, userID string, fields auth.UpdateFields) (*auth.User, error) {
	return m.UpdateFn(ctx, userID, fields)
}

func (m *MockAuthService) Delete(ctx context.Context, userID string) error {
	return m.DeleteFn(ctx, userID)
}
