package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rafabcanedo/obsidian-golang-toolkit/testutil"
	"golang.org/x/crypto/bcrypt"
)

type mockRepository struct {
	createFn                  func(ctx context.Context, username, email, password string) (*User, error)
	findByUsernameFn          func(ctx context.Context, username string) (*User, error)
	findByIDFn                func(ctx context.Context, userID string) (*User, error)
	updateFn                  func(ctx context.Context, userID string, fields UpdateFields) (*User, error)
	deleteFn                  func(ctx context.Context, userID string) error
	saveRefreshTokenFn        func(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	findRefreshTokenByHashFn  func(ctx context.Context, hash string) (*RefreshToken, error)
	deleteRefreshTokenByHashFn func(ctx context.Context, hash string) error
}

func (m *mockRepository) create(ctx context.Context, username, email, password string) (*User, error) {
	return m.createFn(ctx, username, email, password)
}

func (m *mockRepository) findByUsername(ctx context.Context, username string) (*User, error) {
	return m.findByUsernameFn(ctx, username)
}

func (m *mockRepository) findByID(ctx context.Context, userID string) (*User, error) {
	return m.findByIDFn(ctx, userID)
}

func (m *mockRepository) update(ctx context.Context, userID string, fields UpdateFields) (*User, error) {
	return m.updateFn(ctx, userID, fields)
}

func (m *mockRepository) delete(ctx context.Context, userID string) error {
	return m.deleteFn(ctx, userID)
}

func (m *mockRepository) saveRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	if m.saveRefreshTokenFn != nil {
		return m.saveRefreshTokenFn(ctx, userID, tokenHash, expiresAt)
	}
	return nil
}

func (m *mockRepository) findRefreshTokenByHash(ctx context.Context, hash string) (*RefreshToken, error) {
	return m.findRefreshTokenByHashFn(ctx, hash)
}

func (m *mockRepository) deleteRefreshTokenByHash(ctx context.Context, hash string) error {
	if m.deleteRefreshTokenByHashFn != nil {
		return m.deleteRefreshTokenByHashFn(ctx, hash)
	}
	return nil
}

func TestRegister_hashesPassword(t *testing.T) {
	var stored string
	repo := &mockRepository{
		createFn: func(_ context.Context, _, _, password string) (*User, error) {
			stored = password
			return &User{ID: "1"}, nil
		},
	}
	svc := &Service{repo: repo, jwtSecret: "secret"}
	_, err := svc.Register(context.Background(), "john", "john@test.com", "plaintext")
	testutil.NoError(t, err)
	testutil.Equal(t, bcrypt.CompareHashAndPassword([]byte(stored), []byte("plaintext")), nil)
}

func TestLogin_userNotFound_returnsUnauthorized(t *testing.T) {
	repo := &mockRepository{
		findByUsernameFn: func(_ context.Context, _ string) (*User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := &Service{repo: repo, jwtSecret: "secret"}
	_, _, err := svc.Login(context.Background(), "john", "any")
	testutil.Error(t, err)
}

func TestLogin_wrongPassword_returnsUnauthorized(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	repo := &mockRepository{
		findByUsernameFn: func(_ context.Context, _ string) (*User, error) {
			return &User{ID: "1", Password: string(hash)}, nil
		},
	}
	svc := &Service{repo: repo, jwtSecret: "secret"}
	_, _, err := svc.Login(context.Background(), "john", "wrong")
	testutil.Error(t, err)
}

func TestLogin_correctPassword_returnsTokens(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	repo := &mockRepository{
		findByUsernameFn: func(_ context.Context, _ string) (*User, error) {
			return &User{ID: "1", Password: string(hash)}, nil
		},
	}
	svc := &Service{repo: repo, jwtSecret: "secret"}
	accessToken, refreshToken, err := svc.Login(context.Background(), "john", "correct")
	testutil.NoError(t, err)
	if accessToken == "" {
		t.Fatal("expected access token, got empty string")
	}
	if refreshToken == "" {
		t.Fatal("expected refresh token, got empty string")
	}
}

func TestRefresh_invalidToken_returnsUnauthorized(t *testing.T) {
	repo := &mockRepository{
		findRefreshTokenByHashFn: func(_ context.Context, _ string) (*RefreshToken, error) {
			return nil, errors.New("not found")
		},
	}
	svc := &Service{repo: repo, jwtSecret: "secret"}
	_, _, err := svc.Refresh(context.Background(), "invalid-token")
	testutil.Error(t, err)
}

func TestRefresh_expiredToken_returnsUnauthorized(t *testing.T) {
	repo := &mockRepository{
		findRefreshTokenByHashFn: func(_ context.Context, _ string) (*RefreshToken, error) {
			return &RefreshToken{UserID: "1", ExpiresAt: time.Now().Add(-1 * time.Hour)}, nil
		},
	}
	svc := &Service{repo: repo, jwtSecret: "secret"}
	_, _, err := svc.Refresh(context.Background(), "expired-token")
	testutil.Error(t, err)
}

func TestRefresh_validToken_returnsNewTokens(t *testing.T) {
	repo := &mockRepository{
		findRefreshTokenByHashFn: func(_ context.Context, _ string) (*RefreshToken, error) {
			return &RefreshToken{UserID: "1", ExpiresAt: time.Now().Add(1 * time.Hour)}, nil
		},
	}
	svc := &Service{repo: repo, jwtSecret: "secret"}
	newAccessToken, newRefreshToken, err := svc.Refresh(context.Background(), "valid-token")
	testutil.NoError(t, err)
	if newAccessToken == "" {
		t.Fatal("expected new access token, got empty string")
	}
	if newRefreshToken == "" {
		t.Fatal("expected new refresh token, got empty string")
	}
}

func TestUpdate_hashesPasswordWhenProvided(t *testing.T) {
	var stored string
	repo := &mockRepository{
		updateFn: func(_ context.Context, _ string, fields UpdateFields) (*User, error) {
			stored = *fields.Password
			return &User{ID: "1"}, nil
		},
	}
	svc := &Service{repo: repo, jwtSecret: "secret"}
	plain := "newpass"
	_, err := svc.Update(context.Background(), "1", UpdateFields{Password: &plain})
	testutil.NoError(t, err)
	testutil.Equal(t, bcrypt.CompareHashAndPassword([]byte(stored), []byte("newpass")), nil)
}

func TestMe_returnsUser(t *testing.T) {
	repo := &mockRepository{
		findByIDFn: func(_ context.Context, userID string) (*User, error) {
			return &User{ID: userID, Username: "rafa", Email: "rafa@test.com"}, nil
		},
	}
	svc := &Service{repo: repo, jwtSecret: "secret"}
	user, err := svc.Me(context.Background(), "user-123")
	testutil.NoError(t, err)
	testutil.Equal(t, user.ID, "user-123")
}

func TestUpdate_nilPasswordSkipsHash(t *testing.T) {
	repo := &mockRepository{
		updateFn: func(_ context.Context, _ string, fields UpdateFields) (*User, error) {
			if fields.Password != nil {
				t.Fatal("expected nil password, got value")
			}
			return &User{ID: "1"}, nil
		},
	}
	svc := &Service{repo: repo, jwtSecret: "secret"}
	name := "john"
	_, err := svc.Update(context.Background(), "1", UpdateFields{Username: &name})
	testutil.NoError(t, err)
}
