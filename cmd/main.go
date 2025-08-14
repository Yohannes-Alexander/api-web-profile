package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"

	"apiprofile/config"
	httpdelivery "apiprofile/internal/delivery/http"
	gormrepo "apiprofile/internal/repository/gorm"
	usecase "apiprofile/internal/usecase"
	"apiprofile/route"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()

	db, err := config.InitGorm(cfg.DSN)
	if err != nil {
		log.Fatalf("init db: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// optionally auto-migrate (kehilangan native SQL requirement? optional)
	// db.AutoMigrate(&models.User{}, &models.RefreshToken{})

	// wiring DI
	userRepo := gormrepo.NewUserRepository(db)
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	userUC := usecase.NewUserUsecase(userRepo)

	authHandler := httpdelivery.NewAuthHandler(authUC)
	userHandler := httpdelivery.NewUserHandler(userUC)

	r := route.SetupRouter(authHandler, userHandler, cfg.JWTSecret)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("server running %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}

	// keep alive briefly to shutdown gracefully if needed
	time.Sleep(1 * time.Second)
}
