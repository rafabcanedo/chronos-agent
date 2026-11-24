package transactions

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
)

type Servicer interface {
	Create(ctx context.Context, userID, description, category string, amount float64, dueDate, paidAt *string, status string) (*Transaction, error)
	List(ctx context.Context, userID, month, status, category string, p pagination.Params) ([]Transaction, int, error)
	Summary(ctx context.Context, userID, month string) (map[string]StatusSummary, error)
	Update(ctx context.Context, userID, id string, fields UpdateFields) (*Transaction, error)
	Delete(ctx context.Context, userID, id string) error
}

type repositorier interface {
	create(ctx context.Context, userID, description, category string, amount float64, dueDate, paidAt *string, status string) (*Transaction, error)
	list(ctx context.Context, userID, month, status, category string, p pagination.Params) ([]Transaction, int, error)
	update(ctx context.Context, userID, id string, fields UpdateFields) (*Transaction, error)
	delete(ctx context.Context, userID, id string) error
	generateMonthly(ctx context.Context, userID string) error
	markOverdue(ctx context.Context, userID string) error
	findAndMarkPaid(ctx context.Context, userID, description, date string) error
	createExpense(ctx context.Context, userID, description, category, paidAt string, amount float64) error
	summary(ctx context.Context, userID, month string) (map[string]StatusSummary, error)
	queryMonthly(ctx context.Context, userID, month string) (map[string]float64, error)
}

type Service struct {
	repo repositorier
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) Create(ctx context.Context, userID, description, category string, amount float64, dueDate, paidAt *string, status string) (*Transaction, error) {
	if status == "" {
		status = StatusPending
	}
	if status == StatusPaid && paidAt == nil {
		today := time.Now().Format("2006-01-02")
		paidAt = &today
	}
	return s.repo.create(ctx, userID, description, category, amount, dueDate, paidAt, status)
}

func (s *Service) List(ctx context.Context, userID, month, status, category string, p pagination.Params) ([]Transaction, int, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	if err := s.repo.generateMonthly(ctx, userID); err != nil {
		return nil, 0, err
	}

	if err := s.repo.markOverdue(ctx, userID); err != nil {
		return nil, 0, err
	}

	return s.repo.list(ctx, userID, month, status, category, p)
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

func (s *Service) GenerateMonthly(ctx context.Context, userID string) error {
	return s.repo.generateMonthly(ctx, userID)
}

func (s *Service) FindAndMarkPaid(ctx context.Context, userID, description, date string) error {
	return s.repo.findAndMarkPaid(ctx, userID, description, date)
}

func (s *Service) CreateExpense(ctx context.Context, userID, description, category, paidAt string, amount float64) error {
	return s.repo.createExpense(ctx, userID, description, category, paidAt, amount)
}

func (s *Service) Summary(ctx context.Context, userID, month string) (map[string]StatusSummary, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	if err := s.repo.generateMonthly(ctx, userID); err != nil {
		return nil, err
	}

	if err := s.repo.markOverdue(ctx, userID); err != nil {
		return nil, err
	}

	return s.repo.summary(ctx, userID, month)
}

func (s *Service) QueryMonthly(ctx context.Context, userID, month string) (map[string]float64, error) {
	return s.repo.queryMonthly(ctx, userID, month)
}
