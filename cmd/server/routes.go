package main

import (
	"github.com/Siddheshk02/coupon-system/internal/api/handlers"
	"github.com/Siddheshk02/coupon-system/internal/repository"
	"github.com/gorilla/mux"
)

func RegisterRoutes(couponRepo *repository.CouponRepository, itemRepo *repository.ItemRepository, orderRepo *repository.OrderRepository, userRepo *repository.UserRepository) *mux.Router {
	router := mux.NewRouter()

	couponHandler := handlers.NewCouponHandler(couponRepo)
	itemHandler := handlers.NewItemHandler(itemRepo)
	orderHandler := handlers.NewOrdersHandler(orderRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	router.HandleFunc("/admin/coupons", couponHandler.CreateCoupon).Methods("POST")
	router.HandleFunc("/coupons/applicable", couponHandler.GetApplicableCoupons).Methods("GET")
	router.HandleFunc("/coupons/validate", couponHandler.ValidateCoupon).Methods("POST")
	router.HandleFunc("/coupons", couponHandler.GetAllCoupons).Methods("GET")
	router.HandleFunc("/items", itemHandler.AddItem).Methods("POST")
	router.HandleFunc("/items", itemHandler.GetItems).Methods("GET")
	router.HandleFunc("/createorder", orderHandler.AddOrder).Methods("POST")
	router.HandleFunc("/users", userHandler.UserLogin).Methods("POST")

	return router
}
