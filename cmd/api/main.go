package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prasetyoary48/wms/internal/config"
	"github.com/prasetyoary48/wms/internal/delivery/http/handler"
	"github.com/prasetyoary48/wms/internal/delivery/http/router"
	"github.com/prasetyoary48/wms/internal/repository"
	"github.com/prasetyoary48/wms/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("gagal konek ke database: %v", err)
	}

	// Repository layer
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	warehouseRepo := repository.NewWarehouseRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	stockRepo := repository.NewStockRepository(db)
	moveRepo := repository.NewStockMovementRepository(db)
	adjRepo := repository.NewAdjustmentRepository(db)

	// Usecase layer
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret, cfg.JWTExpireHours)
	stockUC := usecase.NewStockUsecase(stockRepo, moveRepo)
	transferUC := usecase.NewTransferUsecase(stockRepo, moveRepo)
	adjustmentUC := usecase.NewAdjustmentUsecase(adjRepo, stockRepo)

	// Handler layer
	handlers := &router.Handlers{
		Auth:       handler.NewAuthHandler(authUC),
		Product:    handler.NewProductHandler(productRepo),
		Warehouse:  handler.NewWarehouseHandler(warehouseRepo, locationRepo),
		Stock:      handler.NewStockHandler(stockUC),
		Transfer:   handler.NewTransferHandler(transferUC),
		Adjustment: handler.NewAdjustmentHandler(adjustmentUC),
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	router.Setup(r, handlers, cfg.JWTSecret)

	log.Printf("server berjalan di port %s (env: %s)", cfg.AppPort, cfg.AppEnv)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}