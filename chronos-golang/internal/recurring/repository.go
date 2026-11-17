package recurring

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafabcanedo/obsidian-golang-toolkit/pagination"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
)

type RecurringBill struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	DefaultAmount float64 `json:"default_amount"`
	BillingDay    int     `json:"billing_day"`
	Category      string  `json:"category"`
	Active        bool    `json:"active"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type UpdateFields struct {
	Name          *string  `json:"name"`
	DefaultAmount *float64 `json:"default_amount"`
	BillingDay    *int     `json:"billing_day"`
	Category      *string  `json:"category"`
	Active        *bool    `json:"active"`
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

func scanRecurringBill(s scanner) (RecurringBill, error) {
	var rb RecurringBill
	var createdAt, updatedAt time.Time

	err := s.Scan(
		&rb.ID, &rb.Name, &rb.DefaultAmount, &rb.BillingDay,
		&rb.Category, &rb.Active, &createdAt, &updatedAt,
	)
	if err != nil {
		return RecurringBill{}, err
	}

	rb.CreatedAt = createdAt.Format(time.RFC3339)
	rb.UpdatedAt = updatedAt.Format(time.RFC3339)

	return rb, nil
}

func (r *repository) list(ctx context.Context, userID string, active bool, p pagination.Params) ([]RecurringBill, int, error) {
	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM recurring_bills
		WHERE user_id = $1 AND active = $2
	`, userID, active).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, name, default_amount, billing_day, category, active, created_at, updated_at
		FROM recurring_bills
		WHERE user_id = $1 AND active = $2
		ORDER BY name ASC
		LIMIT $3 OFFSET $4
	`, userID, active, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	bills := []RecurringBill{}
	for rows.Next() {
		rb, err := scanRecurringBill(rows)
		if err != nil {
			return nil, 0, err
		}
		bills = append(bills, rb)
	}
	return bills, total, rows.Err()
}

func (r *repository) create(ctx context.Context, userID, name, category string, defaultAmount float64, billingDay int) (RecurringBill, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO recurring_bills (user_id, name, default_amount, billing_day, category)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, default_amount, billing_day, category, active, created_at, updated_at
	`, userID, name, defaultAmount, billingDay, category)

	return scanRecurringBill(row)
}

func (r *repository) update(ctx context.Context, userID, id string, fields UpdateFields) (*RecurringBill, error) {
	clauses := []string{}
	args := []any{}
	i := 1

	if fields.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = $%d", i))
		args = append(args, *fields.Name)
		i++
	}
	if fields.DefaultAmount != nil {
		clauses = append(clauses, fmt.Sprintf("default_amount = $%d", i))
		args = append(args, *fields.DefaultAmount)
		i++
	}
	if fields.BillingDay != nil {
		clauses = append(clauses, fmt.Sprintf("billing_day = $%d", i))
		args = append(args, *fields.BillingDay)
		i++
	}
	if fields.Category != nil {
		clauses = append(clauses, fmt.Sprintf("category = $%d", i))
		args = append(args, *fields.Category)
		i++
	}
	if fields.Active != nil {
		clauses = append(clauses, fmt.Sprintf("active = $%d", i))
		args = append(args, *fields.Active)
		i++
	}

	if len(clauses) == 0 {
		return nil, errutil.BadRequest("no fields to update")
	}

	args = append(args, id, userID)
	query := fmt.Sprintf(`
		UPDATE recurring_bills
		SET %s, updated_at = NOW()
		WHERE id = $%d AND user_id = $%d
		RETURNING id, name, default_amount, billing_day, category, active, created_at, updated_at
	`, strings.Join(clauses, ", "), i, i+1)

	row := r.db.QueryRow(ctx, query, args...)
	rb, err := scanRecurringBill(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errutil.NotFound("recurring bill not found")
		}
		return nil, err
	}
	return &rb, nil
}
