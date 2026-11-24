package auth

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

type UpdateFields struct {
	Username *string
	Email    *string
	Password *string
}

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt string
}

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) create(ctx context.Context, username, email, hashedPassword string) (*User, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, created_at
	`, username, email, hashedPassword)

	var u User
	var createdAt time.Time
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &createdAt); err != nil {
		return nil, err
	}
	u.CreatedAt = createdAt.Format(time.RFC3339)
	return &u, nil
}

func (r *repository) update(ctx context.Context, userID string, fields UpdateFields) (*User, error) {
	clauses := []string{}
	args := []any{}
	i := 1

	if fields.Username != nil {
		clauses = append(clauses, fmt.Sprintf("username = $%d", i))
		args = append(args, *fields.Username)
		i++
	}
	if fields.Email != nil {
		clauses = append(clauses, fmt.Sprintf("email = $%d", i))
		args = append(args, *fields.Email)
		i++
	}
	if fields.Password != nil {
		clauses = append(clauses, fmt.Sprintf("password = $%d", i))
		args = append(args, *fields.Password)
		i++
	}

	if len(clauses) == 0 {
		return nil, errutil.BadRequest("no fields to update")
	}

	args = append(args, userID)
	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d
		RETURNING id, username, email, created_at
	`, strings.Join(clauses, ", "), i)

	row := r.db.QueryRow(ctx, query, args...)

	var u User
	var createdAt time.Time
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &createdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errutil.NotFound("user not found")
		}
		return nil, err
	}
	u.CreatedAt = createdAt.Format(time.RFC3339)
	return &u, nil
}

func (r *repository) delete(ctx context.Context, userID string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errutil.NotFound("user not found")
	}
	return nil
}

func (r *repository) findByUsername(ctx context.Context, username string) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE username = $1
	`, username)

	var u User
	var createdAt time.Time
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &createdAt); err != nil {
		return nil, errutil.NotFound("user not found")
	}
	u.CreatedAt = createdAt.Format(time.RFC3339)
	return &u, nil
}

func (r *repository) findByID(ctx context.Context, userID string) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE id = $1
	`, userID)

	var u User
	var createdAt time.Time
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &createdAt); err != nil {
		return nil, errutil.NotFound("user not found")
	}
	u.CreatedAt = createdAt.Format(time.RFC3339)
	return &u, nil
}
