package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Siddheshk02/coupon-system/internal/repository"
)

type ItemHandler struct {
	Repo *repository.ItemRepository
}

func NewItemHandler(repo *repository.ItemRepository) *ItemHandler {
	return &ItemHandler{Repo: repo}
}

func (h *ItemHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var req repository.Item
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Repo.CreateItem(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added successfully"})
}

func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	query := r.URL.Query()
	idStr := query.Get("id")
	var id int
	if idStr != "" {
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id parameter", http.StatusBadRequest)
			return
		}
	}
	category := query.Get("category")

	items, err := h.Repo.GetItems(ctx, id, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items": items,
	})
}
