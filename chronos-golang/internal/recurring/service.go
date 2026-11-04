package recurring

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo *repository
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) List(ctx context.Context, userID, active string) ([]RecurringBill, error) {
	a := true
	if active == "false" {
		a = false
	}
	return s.repo.list(ctx, userID, a)
}

func (s *Service) Create(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (RecurringBill, error) {
	return s.repo.create(ctx, userID, name, category, defaultAmount, billingDay)
}

func (s *Service) Update(ctx context.Context, userID, id string, fields UpdateFields) (*RecurringBill, error) {
	return s.repo.update(ctx, userID, id, fields)
}
