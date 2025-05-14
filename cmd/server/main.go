package main

import (
	"log"
	"net/http"

	"github.com/Siddheshk02/coupon-system/internal/config"
	"github.com/Siddheshk02/coupon-system/internal/db"
	"github.com/Siddheshk02/coupon-system/internal/repository"
)

func main() {
	config.Load()
	db.Init(config.AppConfig)

	couponRepo := repository.NewCouponRepository(db.Conn)
	itemRepo := repository.NewItemRepository(db.Conn)
	orderRepo := repository.NewOrderRepository(db.Conn)
	userRepo := repository.NewUserRepository(db.Conn)

	r := RegisterRoutes(couponRepo, itemRepo, orderRepo, userRepo)

	log.Printf("Server starting on port %s...", config.AppConfig.Port)
	log.Fatal(http.ListenAndServe(":"+config.AppConfig.Port, r))
}
