CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    order_status VARCHAR(50) NOT NULL,
    ordered_at TIMESTAMP NOT NULL,
    coupon_code_used VARCHAR(50) REFERENCES coupons(coupon_code),
    amount_paid DECIMAL(10, 2) NOT NULL
);