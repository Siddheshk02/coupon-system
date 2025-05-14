package repository

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (u *UserRepository) Login(ctx context.Context, req User) error {
	var user User
	err := u.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE name = $1", req.Name).Scan(&user.ID, &user.Password)
	if err == sql.ErrNoRows {
		// User does not exist, create new user
		_, err := u.DB.ExecContext(ctx, "INSERT INTO users (name, password, created_at) VALUES ($1, $2, $3)", req.Name, req.Password, time.Now())
		if err != nil {
			return err
		}
		return nil // User created and logged in
	} else if err != nil {
		return err
	}

	// User exists, check password
	if user.Password != req.Password {
		return sql.ErrNoRows
	}

	return nil // Login successful
}
