CREATE TABLE coupons (
    coupon_code VARCHAR(50) PRIMARY KEY,
    expiry_date TIMESTAMP NOT NULL,
    usage_type VARCHAR(20) NOT NULL,
    applicable_medicine_ids TEXT,
    applicable_categories TEXT,
    min_order_value DECIMAL(10, 2) NOT NULL,
    valid_time_window VARCHAR(50),
    terms_and_conditions TEXT,
    discount_type VARCHAR(20) NOT NULL,
    discount_value DECIMAL(10, 2) NOT NULL,
    max_usage_per_user INT NOT NULL
);