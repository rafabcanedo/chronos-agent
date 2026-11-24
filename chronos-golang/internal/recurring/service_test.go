package recurring

import (
	"context"
	"testing"

	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
	"github.com/rafabcanedo/obsidian-golang-toolkit/testutil"
)

type mockRepository struct {
	listFn   func(ctx context.Context, userID string, active *bool, category, name string, p pagination.Params) ([]RecurringBill, int, error)
	createFn func(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (RecurringBill, error)
	updateFn func(ctx context.Context, userID, id string, fields UpdateFields) (*RecurringBill, error)
	deleteFn func(ctx context.Context, userID, id string) error
}

func (m *mockRepository) list(ctx context.Context, userID string, active *bool, category, name string, p pagination.Params) ([]RecurringBill, int, error) {
	return m.listFn(ctx, userID, active, category, name, p)
}

func (m *mockRepository) create(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (RecurringBill, error) {
	return m.createFn(ctx, userID, name, category, defaultAmount, billingDay)
}

func (m *mockRepository) update(ctx context.Context, userID, id string, fields UpdateFields) (*RecurringBill, error) {
	return m.updateFn(ctx, userID, id, fields)
}

func (m *mockRepository) delete(ctx context.Context, userID, id string) error {
	return m.deleteFn(ctx, userID, id)
}

func TestList_activeTruePassesBoolTrue(t *testing.T) {
	var captured *bool
	repo := &mockRepository{
		listFn: func(_ context.Context, _ string, active *bool, _, _ string, _ pagination.Params) ([]RecurringBill, int, error) {
			captured = active
			return []RecurringBill{}, 0, nil
		},
	}
	svc := &Service{repo: repo}
	svc.List(context.Background(), "u1", "true", "", "", pagination.New(1, 10))
	if captured == nil || *captured != true {
		t.Fatal("expected active=true, got nil or false")
	}
}

func TestList_activeFalsePassesBoolFalse(t *testing.T) {
	var captured *bool
	repo := &mockRepository{
		listFn: func(_ context.Context, _ string, active *bool, _, _ string, _ pagination.Params) ([]RecurringBill, int, error) {
			captured = active
			return []RecurringBill{}, 0, nil
		},
	}
	svc := &Service{repo: repo}
	svc.List(context.Background(), "u1", "false", "", "", pagination.New(1, 10))
	if captured == nil || *captured != false {
		t.Fatal("expected active=false, got nil or true")
	}
}

func TestList_activeEmptyPassesNil(t *testing.T) {
	var captured *bool
	repo := &mockRepository{
		listFn: func(_ context.Context, _ string, active *bool, _, _ string, _ pagination.Params) ([]RecurringBill, int, error) {
			captured = active
			return []RecurringBill{}, 0, nil
		},
	}
	svc := &Service{repo: repo}
	svc.List(context.Background(), "u1", "", "", "", pagination.New(1, 10))
	testutil.Equal(t, captured == nil, true)
}
