package recurring

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
)

type Servicer interface {
	List(ctx context.Context, userID, active, category, name string, p pagination.Params) ([]RecurringBill, int, error)
	Create(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (RecurringBill, error)
	Update(ctx context.Context, userID, id string, fields UpdateFields) (*RecurringBill, error)
	Delete(ctx context.Context, userID, id string) error
}

type repositorier interface {
	list(ctx context.Context, userID string, active *bool, category, name string, p pagination.Params) ([]RecurringBill, int, error)
	create(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (RecurringBill, error)
	update(ctx context.Context, userID, id string, fields UpdateFields) (*RecurringBill, error)
	delete(ctx context.Context, userID, id string) error
}

type Service struct {
	repo repositorier
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) List(ctx context.Context, userID, active, category, name string, p pagination.Params) ([]RecurringBill, int, error) {
	var a *bool
	if active == "true" {
		v := true
		a = &v
	} else if active == "false" {
		v := false
		a = &v
	}

	return s.repo.list(ctx, userID, a, category, name, p)
}

func (s *Service) Create(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (RecurringBill, error) {
	return s.repo.create(ctx, userID, name, category, defaultAmount, billingDay)
}

func (s *Service) Update(ctx context.Context, userID, id string, fields UpdateFields) (*RecurringBill, error) {
	return s.repo.update(ctx, userID, id, fields)
}

func (s *Service) Delete(ctx context.Context, userID, id string) error {
	return s.repo.delete(ctx, userID, id)
}
