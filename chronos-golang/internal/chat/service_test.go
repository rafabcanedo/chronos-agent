package chat

import (
	"context"
	"errors"
	"testing"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/llm"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/recurring"
	"github.com/rafabcanedo/obsidian-golang-toolkit/testutil"
)

type mockTxService struct {
	generateMonthlyFn func(ctx context.Context, userID string) error
	findAndMarkPaidFn func(ctx context.Context, userID, description, date string) error
	createExpenseFn   func(ctx context.Context, userID, description, category, paidAt string, amount float64) error
	queryMonthlyFn    func(ctx context.Context, userID, month string) (map[string]float64, error)
}

func (m *mockTxService) GenerateMonthly(ctx context.Context, userID string) error {
	return m.generateMonthlyFn(ctx, userID)
}

func (m *mockTxService) FindAndMarkPaid(ctx context.Context, userID, description, date string) error {
	return m.findAndMarkPaidFn(ctx, userID, description, date)
}

func (m *mockTxService) CreateExpense(ctx context.Context, userID, description, category, paidAt string, amount float64) error {
	return m.createExpenseFn(ctx, userID, description, category, paidAt, amount)
}

func (m *mockTxService) QueryMonthly(ctx context.Context, userID, month string) (map[string]float64, error) {
	return m.queryMonthlyFn(ctx, userID, month)
}

type mockRecurringService struct {
	createFn func(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (recurring.RecurringBill, error)
}

func (m *mockRecurringService) Create(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (recurring.RecurringBill, error) {
	return m.createFn(ctx, userID, name, category, defaultAmount, billingDay)
}

type mockLLMClient struct {
	parseFn func(ctx context.Context, userMessage, currentDate string) (llm.Action, error)
}

func (m *mockLLMClient) Parse(ctx context.Context, userMessage, currentDate string) (llm.Action, error) {
	return m.parseFn(ctx, userMessage, currentDate)
}

func okTxSvc() *mockTxService {
	return &mockTxService{
		generateMonthlyFn: func(_ context.Context, _ string) error { return nil },
	}
}

func TestHandle_generateMonthlyError(t *testing.T) {
	tx := &mockTxService{
		generateMonthlyFn: func(_ context.Context, _ string) error {
			return errors.New("db error")
		},
	}
	svc := &Service{txSvc: tx, recurringSvc: &mockRecurringService{}, llm: &mockLLMClient{}}
	_, err := svc.Handle(context.Background(), "u1", "qualquer coisa")
	testutil.Error(t, err)
}

func TestHandle_llmParseError(t *testing.T) {
	svc := &Service{
		txSvc:        okTxSvc(),
		recurringSvc: &mockRecurringService{},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{}, errors.New("llm unavailable")
			},
		},
	}
	_, err := svc.Handle(context.Background(), "u1", "xyz")
	testutil.Error(t, err)
}

func TestHandle_unknownAction(t *testing.T) {
	svc := &Service{
		txSvc:        okTxSvc(),
		recurringSvc: &mockRecurringService{},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{Action: "invalid"}, nil
			},
		},
	}
	_, err := svc.Handle(context.Background(), "u1", "xyz")
	testutil.Error(t, err)
}

func TestHandle_pay_success(t *testing.T) {
	tx := okTxSvc()
	tx.findAndMarkPaidFn = func(_ context.Context, _, _, _ string) error { return nil }
	svc := &Service{
		txSvc:        tx,
		recurringSvc: &mockRecurringService{},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{Action: "pay", Description: "internet"}, nil
			},
		},
	}
	msg, err := svc.Handle(context.Background(), "u1", "paguei a internet")
	testutil.NoError(t, err)
	testutil.Equal(t, msg, "Conta de internet marcada como paga.")
}

func TestHandle_pay_withDate(t *testing.T) {
	var capturedDate string
	tx := okTxSvc()
	tx.findAndMarkPaidFn = func(_ context.Context, _, _, date string) error {
		capturedDate = date
		return nil
	}
	date := "2026-11-01"
	svc := &Service{
		txSvc:        tx,
		recurringSvc: &mockRecurringService{},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{Action: "pay", Description: "luz", Date: &date}, nil
			},
		},
	}
	_, err := svc.Handle(context.Background(), "u1", "paguei a luz no dia 1")
	testutil.NoError(t, err)
	testutil.Equal(t, capturedDate, "2026-11-01")
}

func TestHandle_createExpense_nilAmount(t *testing.T) {
	svc := &Service{
		txSvc:        okTxSvc(),
		recurringSvc: &mockRecurringService{},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{Action: "create_expense", Description: "sorvete"}, nil
			},
		},
	}
	_, err := svc.Handle(context.Background(), "u1", "gastei no sorvete")
	testutil.Error(t, err)
}

func TestHandle_createExpense_success(t *testing.T) {
	tx := okTxSvc()
	tx.createExpenseFn = func(_ context.Context, _, _, _, _ string, _ float64) error { return nil }
	amount := 20.0
	svc := &Service{
		txSvc:        tx,
		recurringSvc: &mockRecurringService{},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{Action: "create_expense", Description: "sorvete", Amount: &amount}, nil
			},
		},
	}
	msg, err := svc.Handle(context.Background(), "u1", "gastei 20 no sorvete")
	testutil.NoError(t, err)
	testutil.Equal(t, msg, "Despesa de R$ 20.00 em sorvete registrada com sucesso.")
}

func TestHandle_createRecurring_nilAmount(t *testing.T) {
	svc := &Service{
		txSvc:        okTxSvc(),
		recurringSvc: &mockRecurringService{},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{Action: "create_recurring", Description: "aluguel"}, nil
			},
		},
	}
	_, err := svc.Handle(context.Background(), "u1", "cadastrar aluguel")
	testutil.Error(t, err)
}

func TestHandle_createRecurring_nilBillingDay(t *testing.T) {
	tx := okTxSvc()
	amount := 1200.0
	svc := &Service{
		txSvc:        tx,
		recurringSvc: &mockRecurringService{},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{Action: "create_recurring", Description: "aluguel", Amount: &amount}, nil
			},
		},
	}
	_, err := svc.Handle(context.Background(), "u1", "cadastrar aluguel de 1200")
	testutil.Error(t, err)
}

func TestHandle_createRecurring_success(t *testing.T) {
	tx := okTxSvc()
	amount := 1200.0
	day := 5
	svc := &Service{
		txSvc: tx,
		recurringSvc: &mockRecurringService{
			createFn: func(_ context.Context, _, _, _ string, _ float64, _ int) (recurring.RecurringBill, error) {
				return recurring.RecurringBill{}, nil
			},
		},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{Action: "create_recurring", Description: "aluguel", Amount: &amount, BillingDay: &day}, nil
			},
		},
	}
	msg, err := svc.Handle(context.Background(), "u1", "cadastrar aluguel de 1200 todo dia 5")
	testutil.NoError(t, err)
	testutil.Equal(t, msg, "Conta recorrente aluguel de R$ 1200.00 cadastrada com sucesso.")
}

func TestHandle_query_noTransactions(t *testing.T) {
	tx := okTxSvc()
	tx.queryMonthlyFn = func(_ context.Context, _, _ string) (map[string]float64, error) {
		return map[string]float64{}, nil
	}
	svc := &Service{
		txSvc:        tx,
		recurringSvc: &mockRecurringService{},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{Action: "query"}, nil
			},
		},
	}
	msg, err := svc.Handle(context.Background(), "u1", "quanto gastei?")
	testutil.NoError(t, err)
	testutil.Equal(t, msg, "Nenhuma transação encontrada para este mês.")
}

func TestHandle_query_success(t *testing.T) {
	tx := okTxSvc()
	tx.queryMonthlyFn = func(_ context.Context, _, _ string) (map[string]float64, error) {
		return map[string]float64{"food": 150.0}, nil
	}
	svc := &Service{
		txSvc:        tx,
		recurringSvc: &mockRecurringService{},
		llm: &mockLLMClient{
			parseFn: func(_ context.Context, _, _ string) (llm.Action, error) {
				return llm.Action{Action: "query"}, nil
			},
		},
	}
	msg, err := svc.Handle(context.Background(), "u1", "quanto gastei?")
	testutil.NoError(t, err)
	if msg == "" {
		t.Fatal("expected summary message, got empty string")
	}
}
