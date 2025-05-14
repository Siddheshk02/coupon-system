package repository

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type Coupon struct {
	CouponCode           string    `json:"coupon_code"`
	ExpiryDate           time.Time `json:"expiry_date"`
	UsageType            string    `json:"usage_type"` // one-time / multi-use
	ApplicableMedicines  []string  `json:"applicable_medicine_ids,omitempty"`
	ApplicableCategories []string  `json:"applicable_categories"`
	MinOrderValue        float64   `json:"min_order_value"`
	ValidTimeWindow      string    `json:"valid_time_window,omitempty"`
	TermsAndConditions   string    `json:"terms_and_conditions,omitempty"`
	DiscountType         string    `json:"discount_type"`  // fixed / percentage
	DiscountValue        float64   `json:"discount_value"` // amount / percentage
	MaxUsagePerUser      int       `json:"max_usage_per_user"`
}

type CouponRequest struct {
	CartItems  []CartItem `json:"cart_items"`
	OrderTotal float64    `json:"order_total"`
	Timestamp  string     `json:"timestamp"`
	CouponCode string     `json:"coupon_code"`
}

type CartItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type CouponResult struct {
	CouponCode    string `json:"coupon_code"`
	DiscountValue string `json:"discount_value"`
}

type CouponRepository struct {
	DB    *sql.DB
	Cache *cache.Cache
}

func NewCouponRepository(db *sql.DB) *CouponRepository {
	c := cache.New(20*time.Minute, 30*time.Minute) // 5 min TTL
	return &CouponRepository{DB: db, Cache: c}
}

func (r *CouponRepository) CreateCoupon(ctx context.Context, coupon Coupon) error {
	query := `INSERT INTO coupons (coupon_code, expiry_date, usage_type, applicable_categories, min_order_value, discount_type, discount_value, max_usage_per_user) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	applicableCategories := strings.Join(coupon.ApplicableCategories, ",")

	_, err := r.DB.ExecContext(ctx, query, coupon.CouponCode, coupon.ExpiryDate, coupon.UsageType, applicableCategories, coupon.MinOrderValue, coupon.DiscountType, coupon.DiscountValue, coupon.MaxUsagePerUser)
	if err == nil {
		r.Cache.Delete("all_coupons") // Invalidate the cache
	}
	return err
}

func (r *CouponRepository) GetCoupons(ctx context.Context, couponReq CouponRequest) ([]CouponResult, error) {
	query := `SELECT coupon_code, expiry_date, usage_type, applicable_categories, min_order_value, discount_type, discount_value, max_usage_per_user 
			  FROM coupons
			  WHERE expiry_date > $1 AND min_order_value <= $2`

	// Calculate the total price of all cart items
	var totalPrice float64
	for _, item := range couponReq.CartItems {
		totalPrice += item.Price
	}
	rows, err := r.DB.QueryContext(ctx, query, couponReq.Timestamp, totalPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applicableCoupons []CouponResult

	for rows.Next() {
		var coupon Coupon
		var applicableCategories string
		err := rows.Scan(&coupon.CouponCode, &coupon.ExpiryDate, &coupon.UsageType, &applicableCategories, &coupon.MinOrderValue, &coupon.DiscountType, &coupon.DiscountValue, &coupon.MaxUsagePerUser)
		if err != nil {
			return nil, err
		}

		// Check if the coupon is applicable to the cart items
		coupon.ApplicableCategories = splitCommaSeparatedString(applicableCategories)
		isApplicable := false
		for _, item := range couponReq.CartItems {
			if containsString(coupon.ApplicableCategories, item.Category) {
				isApplicable = true
				break
			}
		}
		if !isApplicable {
			continue
		}

		// Calculate the discount
		var discountValue float64
		if coupon.DiscountType == "percentage" {
			discountValue = (couponReq.OrderTotal * coupon.DiscountValue) / 100
		} else if coupon.DiscountType == "fixed" {
			discountValue = coupon.DiscountValue
		}

		discountStr := strconv.FormatFloat(discountValue, 'f', 2, 64)

		applicableCoupons = append(applicableCoupons, CouponResult{
			CouponCode:    coupon.CouponCode,
			DiscountValue: discountStr,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return applicableCoupons, nil
}

func (r *CouponRepository) CheckCoupon(ctx context.Context, couponReq CouponRequest, userID string) (float64, float64, error) {
	query := `SELECT coupon_code, expiry_date, usage_type, applicable_categories, min_order_value, discount_type, discount_value, max_usage_per_user 
			  FROM coupons
			  WHERE coupon_code = $1 AND min_order_value <= $2 AND expiry_date > $3`

	// Calculate the total price of all cart items
	var totalPrice float64
	for _, item := range couponReq.CartItems {
		totalPrice += item.Price
	}
	var coupon Coupon
	var applicableCategories string
	err := r.DB.QueryRowContext(ctx, query, couponReq.CouponCode, totalPrice, couponReq.Timestamp).
		Scan(&coupon.CouponCode, &coupon.ExpiryDate, &coupon.UsageType, &applicableCategories, &coupon.MinOrderValue, &coupon.DiscountType, &coupon.DiscountValue, &coupon.MaxUsagePerUser)
	if err != nil {
		return 0, 0, err
	}

	// Check if the coupon is applicable to the cart items
	coupon.ApplicableCategories = splitCommaSeparatedString(applicableCategories)
	isApplicable := false
	for _, item := range couponReq.CartItems {
		if containsString(coupon.ApplicableCategories, item.Category) {
			isApplicable = true
			break
		}
	}
	if !isApplicable {
		return 0, 0, nil
	}

	// Check if the coupon is a one-time usage or multi-use usage type
	usageCheckQuery := `SELECT COUNT(*) FROM orders WHERE coupon_code_used = $1 AND user_id = $2`
	var usageCount int
	err = r.DB.QueryRowContext(ctx, usageCheckQuery, couponReq.CouponCode, userID).Scan(&usageCount)
	if err != nil {
		return 0, 0, err
	}

	if coupon.UsageType == "one-time" {
		if usageCount >= 1 {
			return 0, 0, nil
		}
	}

	var itemsDiscount, chargesDiscount float64
	if coupon.DiscountType == "percentage" {
		discountValue := coupon.DiscountValue

		// Items discount
		itemsDiscount = totalPrice * (discountValue / 100)

		// Charges discount
		totalCharges := couponReq.OrderTotal - totalPrice
		chargesDiscount = totalCharges * (discountValue / 100)

		return itemsDiscount, chargesDiscount, nil
	} else if coupon.DiscountType == "fixed" {
		fixedDiscount := coupon.DiscountValue

		totalCharges := couponReq.OrderTotal - totalPrice
		if fixedDiscount >= totalCharges {
			// Apply the entire fixed discount to the charges and set the items discount to 0
			chargesDiscount = fixedDiscount
			itemsDiscount = 0
		} else {
			// Apply the fixed discount to both the items and the charges
			itemsDiscount = fixedDiscount
			chargesDiscount = fixedDiscount
		}

		return itemsDiscount, chargesDiscount, nil
	}

	return itemsDiscount, chargesDiscount, nil
}

func (r *CouponRepository) GetAllCoupons(ctx context.Context) ([]Coupon, error) {
	if cached, found := r.Cache.Get("all_coupons"); found {
		return cached.([]Coupon), nil
	}

	rows, err := r.DB.QueryContext(ctx, `SELECT coupon_code, expiry_date, usage_type, applicable_categories, min_order_value, discount_type, discount_value, max_usage_per_user FROM coupons`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []Coupon
	for rows.Next() {
		var coupon Coupon
		var applicableCategories string
		if err := rows.Scan(&coupon.CouponCode, &coupon.ExpiryDate, &coupon.UsageType, &applicableCategories, &coupon.MinOrderValue, &coupon.DiscountType, &coupon.DiscountValue, &coupon.MaxUsagePerUser); err != nil {
			return nil, err
		}

		coupon.ApplicableCategories = splitCommaSeparatedString(applicableCategories)
		coupons = append(coupons, coupon)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	r.Cache.Set("all_coupons", coupons, cache.DefaultExpiration)
	return coupons, nil
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func splitCommaSeparatedString(input string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(input, ",")
}
