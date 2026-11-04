package transactions

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
)

const (
	StatusPending = "pending"
	StatusPaid    = "paid"
	StatusOverdue = "overdue"
)

const (
	CategoryMercado      = "mercado"
	CategoryRestaurantes = "restaurantes"
	CategoryMoradia      = "moradia"
	CategoryTransporte   = "transporte"
	CategoryLazer        = "lazer"
	CategoryOutros       = "outros"
)

type Transaction struct {
	ID              string  `json:"id"`
	RecurringBillID *string `json:"recurring_bill_id"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	Category        string  `json:"category"`
	DueDate         *string `json:"due_date"`
	Status          string  `json:"status"`
	PaidAt          *string `json:"paid_at"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type UpdateFields struct {
	Description *string  `json:"description"`
	Amount      *float64 `json:"amount"`
	Category    *string  `json:"category"`
	DueDate     *string  `json:"due_date"`
	Status      *string  `json:"status"`
	PaidAt      *string  `json:"paid_at"`
}

type scanner interface {
	Scan(dest ...any) error
}

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func scanTransaction(s scanner) (Transaction, error) {
	var t Transaction
	var recurringBillID *string
	var dueDate, paidAt *time.Time
	var createdAt, updatedAt time.Time

	err := s.Scan(
		&t.ID, &recurringBillID, &t.Description, &t.Amount, &t.Category,
		&dueDate, &t.Status, &paidAt, &createdAt, &updatedAt,
	)
	if err != nil {
		return Transaction{}, err
	}

	t.RecurringBillID = recurringBillID

	if dueDate != nil {
		s := dueDate.Format("2006-01-02")
		t.DueDate = &s
	}
	if paidAt != nil {
		s := paidAt.Format("2006-01-02")
		t.PaidAt = &s
	}

	t.CreatedAt = createdAt.Format(time.RFC3339)
	t.UpdatedAt = updatedAt.Format(time.RFC3339)

	return t, nil
}

func (r *repository) generateMonthly(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO transactions (user_id, recurring_bill_id, description, amount, category, due_date, status)
		SELECT
			rb.user_id,
			rb.id,
			rb.name,
			rb.default_amount,
			rb.category,
			MAKE_DATE(
				EXTRACT(YEAR FROM CURRENT_DATE)::int,
				EXTRACT(MONTH FROM CURRENT_DATE)::int,
				LEAST(
					rb.billing_day,
					EXTRACT(DAY FROM (DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month - 1 day'))::int
				)
			),
			'pending'
		FROM recurring_bills rb
		WHERE rb.user_id = $1 AND rb.active = TRUE
		ON CONFLICT DO NOTHING
	`, userID)
	return err
}

func (r *repository) markOverdue(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE transactions
		SET status = 'overdue', updated_at = NOW()
		WHERE user_id = $1
		  AND status = 'pending'
		  AND due_date < CURRENT_DATE
	`, userID)
	return err
}

func (r *repository) list(ctx context.Context, userID, month, status string) ([]Transaction, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, recurring_bill_id, description, amount, category,
		       due_date, status, paid_at, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		  AND DATE_TRUNC('month', COALESCE(due_date, paid_at)) = DATE_TRUNC('month', $2::date)
		  AND ($3 = '' OR status = $3)
		ORDER BY due_date ASC NULLS LAST
	`, userID, month+"-01", status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := []Transaction{}
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, rows.Err()
}

func (r *repository) update(ctx context.Context, userID, id string, fields UpdateFields) (*Transaction, error) {
	clauses := []string{}
	args := []any{}
	i := 1

	if fields.Description != nil {
		clauses = append(clauses, fmt.Sprintf("description = $%d", i))
		args = append(args, *fields.Description)
		i++
	}
	if fields.Amount != nil {
		clauses = append(clauses, fmt.Sprintf("amount = $%d", i))
		args = append(args, *fields.Amount)
		i++
	}
	if fields.Category != nil {
		clauses = append(clauses, fmt.Sprintf("category = $%d", i))
		args = append(args, *fields.Category)
		i++
	}
	if fields.DueDate != nil {
		clauses = append(clauses, fmt.Sprintf("due_date = $%d", i))
		args = append(args, *fields.DueDate)
		i++
	}
	if fields.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", i))
		args = append(args, *fields.Status)
		i++
	}
	if fields.PaidAt != nil {
		clauses = append(clauses, fmt.Sprintf("paid_at = $%d", i))
		args = append(args, *fields.PaidAt)
		i++
	}

	if len(clauses) == 0 {
		return nil, errutil.BadRequest("no fields to update")
	}

	args = append(args, id, userID)
	query := fmt.Sprintf(`
		UPDATE transactions
		SET %s
		WHERE id = $%d AND user_id = $%d
		RETURNING id, recurring_bill_id, description, amount, category,
		          due_date, status, paid_at, created_at, updated_at
	`, strings.Join(clauses, ", "), i, i+1)

	row := r.db.QueryRow(ctx, query, args...)
	t, err := scanTransaction(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errutil.NotFound("transaction not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *repository) delete(ctx context.Context, userID, id string) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM transactions
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errutil.NotFound("transaction not found")
	}
	return nil
}

func (r *repository) findAndMarkPaid(ctx context.Context, userID, description, date string) error {
	result, err := r.db.Exec(ctx, `
		UPDATE transactions
		SET status = 'paid', paid_at = $3::date, updated_at = NOW()
		WHERE id = (
			SELECT id FROM transactions
			WHERE user_id = $1
			  AND status IN ('pending', 'overdue')
			  AND description ILIKE '%' || $2 || '%'
			  AND DATE_TRUNC('month', due_date) = DATE_TRUNC('month', $3::date)
			LIMIT 1
		)
	`, userID, description, date)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errutil.NotFound("Conta não encontrada para este mês.")
	}
	return nil
}

func (r *repository) createExpense(ctx context.Context, userID, description, category, paidAt string, amount float64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO transactions (user_id, description, amount, category, status, paid_at)
		VALUES ($1, $2, $3, $4, 'paid', $5::date)
	`, userID, description, amount, category, paidAt)
	return err
}

func (r *repository) queryMonthly(ctx context.Context, userID, month string) (map[string]float64, error) {
	rows, err := r.db.Query(ctx, `
		SELECT category, SUM(amount)
		FROM transactions
		WHERE user_id = $1
		  AND DATE_TRUNC('month', COALESCE(due_date, paid_at)) = DATE_TRUNC('month', $2::date)
		GROUP BY category
		ORDER BY SUM(amount) DESC
	`, userID, month+"-01")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]float64{}
	for rows.Next() {
		var category string
		var total float64
		if err := rows.Scan(&category, &total); err != nil {
			return nil, err
		}
		result[category] = total
	}
	return result, rows.Err()
}
