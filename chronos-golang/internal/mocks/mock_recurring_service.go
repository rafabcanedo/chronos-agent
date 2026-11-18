package mocks

import (
	"context"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/recurring"
	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
)

type MockRecurringService struct {
	ListFn   func(ctx context.Context, userID, active, category, name string, p pagination.Params) ([]recurring.RecurringBill, int, error)
	CreateFn func(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (recurring.RecurringBill, error)
	UpdateFn func(ctx context.Context, userID, id string, fields recurring.UpdateFields) (*recurring.RecurringBill, error)
	DeleteFn func(ctx context.Context, userID, id string) error
}

func (m *MockRecurringService) List(ctx context.Context, userID, active, category, name string, p pagination.Params) ([]recurring.RecurringBill, int, error) {
	return m.ListFn(ctx, userID, active, category, name, p)
}

func (m *MockRecurringService) Create(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (recurring.RecurringBill, error) {
	return m.CreateFn(ctx, userID, name, category, defaultAmount, billingDay)
}

func (m *MockRecurringService) Update(ctx context.Context, userID, id string, fields recurring.UpdateFields) (*recurring.RecurringBill, error) {
	return m.UpdateFn(ctx, userID, id, fields)
}

func (m *MockRecurringService) Delete(ctx context.Context, userID, id string) error {
	return m.DeleteFn(ctx, userID, id)
}
