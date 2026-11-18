package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/errutil"
)

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
