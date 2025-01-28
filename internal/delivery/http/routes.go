package http

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, orderHandler *OrderHandler) {
	api := e.Group("/api/v1/orders")
	api.POST("", orderHandler.CreateOrder)
	api.GET("/:order_id", orderHandler.GetOrders)

}
