package repository

import (
	"context"
	"database/sql"
	"time"
)

type Order struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	OrderStatus    string    `json:"order_status"`
	OrderedAt      time.Time `json:"ordered_at"`
	CouponCodeUsed string    `json:"coupon_code_used"`
	AmountPaid     float64   `json:"amount_paid"`
}

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (o *OrderRepository) CreateOrder(ctx context.Context, req Order) error {
	query := `INSERT INTO orders (user_id, order_status, ordered_at, coupon_code_used, amount_paid) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := o.DB.ExecContext(ctx, query, req.UserID, req.OrderStatus, req.OrderedAt, req.CouponCodeUsed, req.AmountPaid)
	return err
}
