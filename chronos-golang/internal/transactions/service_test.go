package transactions

import (
	"context"
	"testing"

	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
	"github.com/rafabcanedo/obsidian-golang-toolkit/testutil"
)

type mockRepository struct {
	createFn          func(ctx context.Context, userID, description, category string, amount float64, dueDate, paidAt *string, status string) (*Transaction, error)
	listFn            func(ctx context.Context, userID, month, status, category string, p pagination.Params) ([]Transaction, int, error)
	updateFn          func(ctx context.Context, userID, id string, fields UpdateFields) (*Transaction, error)
	deleteFn          func(ctx context.Context, userID, id string) error
	generateMonthlyFn func(ctx context.Context, userID string) error
	markOverdueFn     func(ctx context.Context, userID string) error
	findAndMarkPaidFn func(ctx context.Context, userID, description, date string) error
	createExpenseFn   func(ctx context.Context, userID, description, category, paidAt string, amount float64) error
	summaryFn         func(ctx context.Context, userID, month string) (map[string]StatusSummary, error)
	queryMonthlyFn    func(ctx context.Context, userID, month string) (map[string]float64, error)
}

func (m *mockRepository) create(ctx context.Context, userID, description, category string, amount float64, dueDate, paidAt *string, status string) (*Transaction, error) {
	return m.createFn(ctx, userID, description, category, amount, dueDate, paidAt, status)
}

func (m *mockRepository) list(ctx context.Context, userID, month, status, category string, p pagination.Params) ([]Transaction, int, error) {
	return m.listFn(ctx, userID, month, status, category, p)
}

func (m *mockRepository) update(ctx context.Context, userID, id string, fields UpdateFields) (*Transaction, error) {
	return m.updateFn(ctx, userID, id, fields)
}

func (m *mockRepository) delete(ctx context.Context, userID, id string) error {
	return m.deleteFn(ctx, userID, id)
}

func (m *mockRepository) generateMonthly(ctx context.Context, userID string) error {
	return m.generateMonthlyFn(ctx, userID)
}

func (m *mockRepository) markOverdue(ctx context.Context, userID string) error {
	return m.markOverdueFn(ctx, userID)
}

func (m *mockRepository) findAndMarkPaid(ctx context.Context, userID, description, date string) error {
	return m.findAndMarkPaidFn(ctx, userID, description, date)
}

func (m *mockRepository) createExpense(ctx context.Context, userID, description, category, paidAt string, amount float64) error {
	return m.createExpenseFn(ctx, userID, description, category, paidAt, amount)
}

func (m *mockRepository) summary(ctx context.Context, userID, month string) (map[string]StatusSummary, error) {
	return m.summaryFn(ctx, userID, month)
}

func (m *mockRepository) queryMonthly(ctx context.Context, userID, month string) (map[string]float64, error) {
	return m.queryMonthlyFn(ctx, userID, month)
}

func TestCreate_defaultsStatusToPending(t *testing.T) {
	var capturedStatus string
	repo := &mockRepository{
		createFn: func(_ context.Context, _, _, _ string, _ float64, _, _ *string, status string) (*Transaction, error) {
			capturedStatus = status
			return &Transaction{}, nil
		},
	}
	svc := &Service{repo: repo}
	_, err := svc.Create(context.Background(), "u1", "desc", "food", 50.0, nil, nil, "")
	testutil.NoError(t, err)
	testutil.Equal(t, capturedStatus, StatusPending)
}

func TestCreate_setsPaidAtWhenStatusPaid(t *testing.T) {
	var capturedPaidAt *string
	repo := &mockRepository{
		createFn: func(_ context.Context, _, _, _ string, _ float64, _, paidAt *string, _ string) (*Transaction, error) {
			capturedPaidAt = paidAt
			return &Transaction{}, nil
		},
	}
	svc := &Service{repo: repo}
	_, err := svc.Create(context.Background(), "u1", "desc", "food", 50.0, nil, nil, StatusPaid)
	testutil.NoError(t, err)
	if capturedPaidAt == nil {
		t.Fatal("expected paidAt to be set, got nil")
	}
}

func TestUpdate_setsPaidAtWhenMarkingPaid(t *testing.T) {
	var capturedFields UpdateFields
	repo := &mockRepository{
		updateFn: func(_ context.Context, _, _ string, fields UpdateFields) (*Transaction, error) {
			capturedFields = fields
			return &Transaction{}, nil
		},
	}
	svc := &Service{repo: repo}
	status := StatusPaid
	_, err := svc.Update(context.Background(), "u1", "tx1", UpdateFields{Status: &status})
	testutil.NoError(t, err)
	if capturedFields.PaidAt == nil {
		t.Fatal("expected PaidAt to be set")
	}
}
