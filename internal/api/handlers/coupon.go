package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Siddheshk02/coupon-system/internal/repository"
)

type CouponHandler struct {
	Repo *repository.CouponRepository
}

func NewCouponHandler(repo *repository.CouponRepository) *CouponHandler {
	return &CouponHandler{Repo: repo}
}

func (h *CouponHandler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var req repository.Coupon
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Repo.CreateCoupon(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "coupon created"})
}

func (h *CouponHandler) GetAllCoupons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	res, err := h.Repo.GetAllCoupons(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"coupons": res,
	})
}

func (h *CouponHandler) GetApplicableCoupons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var req repository.CouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	switch {
	case req.CartItems == nil || len(req.CartItems) == 0:
		http.Error(w, "invalid request body: no items added", http.StatusBadRequest)
		return
	case req.OrderTotal == 0:
		http.Error(w, "invalid request body: order total is zero", http.StatusBadRequest)
		return
	case req.Timestamp == "":
		http.Error(w, "invalid request body: timestamp required", http.StatusBadRequest)
		return
	}

	res, err := h.Repo.GetCoupons(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"applicable_coupons": res,
	})
}

func (h *CouponHandler) ValidateCoupon(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var req repository.CouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	userID := query.Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	switch {
	case req.CartItems == nil || len(req.CartItems) == 0:
		http.Error(w, "invalid request body: no items added", http.StatusBadRequest)
		return
	case req.OrderTotal == 0:
		http.Error(w, "invalid request body: order total is zero", http.StatusBadRequest)
		return
	case req.Timestamp == "":
		http.Error(w, "invalid request body: timestamp required", http.StatusBadRequest)
		return
	case req.CouponCode == "":
		http.Error(w, "invalid request body: coupon code required", http.StatusBadRequest)
		return
	}

	items_discount, charges_discount, err := h.Repo.CheckCoupon(ctx, req, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if items_discount == 0 && charges_discount == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"is_valid": false,
			"reason":   "coupon expired or not applicable",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"is_valid": true,
		"discount": map[string]float64{
			"items_discount":   items_discount,
			"charges_discount": charges_discount,
		},
		"message": "coupon applied successfully",
	})
}
