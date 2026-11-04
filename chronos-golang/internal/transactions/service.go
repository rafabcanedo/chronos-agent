package transactions

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo *repository
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) List(ctx context.Context, userID, month, status string) ([]Transaction, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	if err := s.repo.generateMonthly(ctx, userID); err != nil {
		return nil, err
	}

	return s.repo.list(ctx, userID, month, status)
}

func (s *Service) Update(ctx context.Context, userID, id string, fields UpdateFields) (*Transaction, error) {
	if fields.Status != nil && *fields.Status == "paid" && fields.PaidAt == nil {
		today := time.Now().Format("2006-01-02")
		fields.PaidAt = &today
	}

	return s.repo.update(ctx, userID, id, fields)
}

func (s *Service) Delete(ctx context.Context, userID, id string) error {
	return s.repo.delete(ctx, userID, id)
}
