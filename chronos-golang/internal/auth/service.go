package auth

import (
	"context"
	"time"

	"github.com/rafabcanedo/obsidian-golang-toolkit/api/jwtauth"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Servicer interface {
	Register(ctx context.Context, username, email, password string) (*User, error)
	Login(ctx context.Context, username, password string) (string, error)
	Update(ctx context.Context, userID string, fields UpdateFields) (*User, error)
	Delete(ctx context.Context, userID string) error
}

type Service struct {
	repo      *repository
	jwtSecret string
}

func NewService(db *pgxpool.Pool, jwtSecret string) *Service {
	return &Service{
		repo:      newRepository(db),
		jwtSecret: jwtSecret,
	}
}

func (s *Service) Register(ctx context.Context, username, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.repo.create(ctx, username, email, string(hash))
}

func (s *Service) Update(ctx context.Context, userID string, fields UpdateFields) (*User, error) {
	if fields.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*fields.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashed := string(hash)
		fields.Password = &hashed
	}
	return s.repo.update(ctx, userID, fields)
}

func (s *Service) Delete(ctx context.Context, userID string) error {
	return s.repo.delete(ctx, userID)
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.findByUsername(ctx, username)
	if err != nil {
		return "", errutil.Unauthorized("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errutil.Unauthorized("invalid credentials")
	}

	token, err := jwtauth.Generate(s.jwtSecret, jwtauth.Claims{
		Subject: user.ID,
		Expiry:  time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
