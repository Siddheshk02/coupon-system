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

func (u *UserRepository) Login(ctx context.Context, req User) (int, error) {
	var user User
	var id int
	err := u.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE name = $1", req.Name).Scan(&user.ID, &user.Password)
	if err == sql.ErrNoRows {
		// User does not exist, create new user
		err := u.DB.QueryRowContext(ctx, "INSERT INTO users (name, password, created_at) VALUES ($1, $2, $3) RETURNING id", req.Name, req.Password, time.Now()).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil // User created and logged in
	} else if err != nil {
		return 0, err
	}

	// User exists, check password
	if user.Password != req.Password {
		return 0, sql.ErrNoRows
	}

	return user.ID, nil // Login successful
}
