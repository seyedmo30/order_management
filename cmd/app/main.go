package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/seyedmo30/order_management/internal/config"
	"github.com/seyedmo30/order_management/internal/delivery/http"
	"github.com/seyedmo30/order_management/internal/process"
	"github.com/seyedmo30/order_management/internal/repository"
	"github.com/seyedmo30/order_management/internal/usecase"
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Load configuration
	cfg := config.Load()
	ctx := context.Background()

	// Initialize repository
	repo := repository.NewOrderManagementRepository(cfg.DatabaseConfig)

	process := process.NewProcessUseCase(cfg)

	usecase := usecase.NewOrderUseCase(cfg, repo, process)

	// Start the order processing in a background worker (goroutine)
	go func() {

		if err := usecase.ListAggregateOrderReport(ctx); err != nil {
			log.Printf("Error in ListAggregateOrderReport: %v", err)
		}

	}()
	go func() {

		if err := usecase.ProcessOrder(ctx); err != nil {
			log.Fatalf("Error while processing orders: %v", err)
		}

	}()

	// Set up Echo instance
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Initialize handlers
	orderHandler := http.NewOrderHandler(usecase)

	// Register routes
	http.RegisterRoutes(e, orderHandler)
	// Start the server
	log.Fatal(e.Start(":8099"))
}
