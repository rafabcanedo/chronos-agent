package chat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/llm"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/recurring"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/transactions"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
)

type Service struct {
	txSvc        *transactions.Service
	recurringSvc *recurring.Service
	llm          *llm.Client
}

func NewService(txSvc *transactions.Service, recurringSvc *recurring.Service, llmClient *llm.Client) *Service {
	return &Service{
		txSvc:        txSvc,
		recurringSvc: recurringSvc,
		llm:          llmClient,
	}
}

func (s *Service) Handle(ctx context.Context, userID, message string) (string, error) {
	today := time.Now().Format("2006-01-02")

	if err := s.txSvc.GenerateMonthly(ctx, userID); err != nil {
		return "", err
	}

	action, err := s.llm.Parse(ctx, message, today)
	if err != nil {
		return "", errutil.BadRequest("Não entendi o que você quis dizer.")
	}

	switch action.Action {
	case "pay":
		return s.handlePay(ctx, userID, action, today)
	case "create_expense":
		return s.handleCreateExpense(ctx, userID, action, today)
	case "create_recurring":
		return s.handleCreateRecurring(ctx, userID, action)
	case "query":
		return s.handleQuery(ctx, userID, action, today)
	default:
		return "", errutil.BadRequest("Não entendi o que você quis dizer.")
	}
}

func (s *Service) handlePay(ctx context.Context, userID string, action llm.Action, today string) (string, error) {
	date := today
	if action.Date != nil {
		date = *action.Date
	}

	if err := s.txSvc.FindAndMarkPaid(ctx, userID, action.Description, date); err != nil {
		return "", err
	}

	return fmt.Sprintf("Conta de %s marcada como paga.", action.Description), nil
}

func (s *Service) handleCreateExpense(ctx context.Context, userID string, action llm.Action, today string) (string, error) {
	if action.Amount == nil {
		return "", errutil.BadRequest("Não entendi o valor da despesa.")
	}

	date := today
	if action.Date != nil {
		date = *action.Date
	}

	category := "outros"
	if action.Category != nil {
		category = *action.Category
	}

	if err := s.txSvc.CreateExpense(ctx, userID, action.Description, category, date, *action.Amount); err != nil {
		return "", err
	}

	return fmt.Sprintf("Despesa de R$ %.2f em %s registrada com sucesso.", *action.Amount, action.Description), nil
}

func (s *Service) handleCreateRecurring(ctx context.Context, userID string, action llm.Action) (string, error) {
	if action.Amount == nil {
		return "", errutil.BadRequest("Não entendi o valor da conta recorrente.")
	}
	if action.BillingDay == nil {
		return "", errutil.BadRequest("Não entendi o dia de vencimento.")
	}

	category := "outros"
	if action.Category != nil {
		category = *action.Category
	}

	if _, err := s.recurringSvc.Create(ctx, userID, action.Description, category, *action.Amount, *action.BillingDay); err != nil {
		return "", err
	}

	return fmt.Sprintf("Conta recorrente %s de R$ %.2f cadastrada com sucesso.", action.Description, *action.Amount), nil
}

func (s *Service) handleQuery(ctx context.Context, userID string, action llm.Action, today string) (string, error) {
	month := today[:7]
	if action.Date != nil {
		month = (*action.Date)[:7]
	}

	totals, err := s.txSvc.QueryMonthly(ctx, userID, month)
	if err != nil {
		return "", err
	}

	if len(totals) == 0 {
		return "Nenhuma transação encontrada para este mês.", nil
	}

	var total float64
	var parts []string
	for category, amount := range totals {
		parts = append(parts, fmt.Sprintf("%s: R$ %.2f", category, amount))
		total += amount
	}

	return fmt.Sprintf("Resumo do mês: %s. Total: R$ %.2f.", strings.Join(parts, ", "), total), nil
}
