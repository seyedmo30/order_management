package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/seyedmo30/order_management/internal/config"
	"github.com/seyedmo30/order_management/internal/delivery/http"
	"github.com/seyedmo30/order_management/internal/process"
	"github.com/seyedmo30/order_management/internal/repository"
	"github.com/seyedmo30/order_management/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize repository
	repo := repository.NewOrderManagementRepository(cfg.DatabaseConfig)

	process := process.NewProcessUseCase()

	usecase := usecase.NewOrderUseCase(repo,process)
	// Set up Echo instance
	e := echo.New()

	// Initialize handlers
	orderHandler := http.NewOrderHandler(usecase)

	// Register routes
	http.RegisterRoutes(e, orderHandler)

	// Start the server
	log.Fatal(e.Start(":" + (cfg.ServiceConfig.Port)))
}
