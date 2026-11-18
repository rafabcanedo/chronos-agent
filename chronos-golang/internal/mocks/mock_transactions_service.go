package mocks

import (
	"context"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/transactions"
	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
)

type MockTransactionsService struct {
	CreateFn  func(ctx context.Context, userID, description, category string, amount float64, dueDate, paidAt *string, status string) (*transactions.Transaction, error)
	ListFn    func(ctx context.Context, userID, month, status, category string, p pagination.Params) ([]transactions.Transaction, int, error)
	SummaryFn func(ctx context.Context, userID, month string) (map[string]transactions.StatusSummary, error)
	UpdateFn  func(ctx context.Context, userID, id string, fields transactions.UpdateFields) (*transactions.Transaction, error)
	DeleteFn  func(ctx context.Context, userID, id string) error
}

func (m *MockTransactionsService) Create(ctx context.Context, userID, description, category string, amount float64, dueDate, paidAt *string, status string) (*transactions.Transaction, error) {
	return m.CreateFn(ctx, userID, description, category, amount, dueDate, paidAt, status)
}

func (m *MockTransactionsService) List(ctx context.Context, userID, month, status, category string, p pagination.Params) ([]transactions.Transaction, int, error) {
	return m.ListFn(ctx, userID, month, status, category, p)
}

func (m *MockTransactionsService) Summary(ctx context.Context, userID, month string) (map[string]transactions.StatusSummary, error) {
	return m.SummaryFn(ctx, userID, month)
}

func (m *MockTransactionsService) Update(ctx context.Context, userID, id string, fields transactions.UpdateFields) (*transactions.Transaction, error) {
	return m.UpdateFn(ctx, userID, id, fields)
}

func (m *MockTransactionsService) Delete(ctx context.Context, userID, id string) error {
	return m.DeleteFn(ctx, userID, id)
}
