package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/seyedmo30/order_management/internal/dto"
	"github.com/seyedmo30/order_management/internal/interfaces"
)

type OrderHandler struct {
	usecase interfaces.OrderUseCase
}

func NewOrderHandler(usecase interfaces.OrderUseCase) *OrderHandler {
	return &OrderHandler{usecase: usecase}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var req dto.CreateOrderHttpHandlerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	err := h.usecase.CreateOrder(c.Request().Context(), dto.CreateOrderUsecaseRequest{BaseCreateOrderRequest: req.BaseCreateOrderRequest})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusAccepted, echo.Map{"message": "order created successfully"})
}

func (h *OrderHandler) GetOrders(c echo.Context) error {
	// Implement fetching orders
	return c.JSON(http.StatusOK, echo.Map{"orders": "[]dto.Order{}"})
}
