package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Siddheshk02/coupon-system/internal/repository"
)

type OrderHandler struct {
	Repo *repository.OrderRepository
}

func NewOrdersHandler(repo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{Repo: repo}
}

func (h *OrderHandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var req repository.Order
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Repo.CreateOrder(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order placed successfully"})
}
