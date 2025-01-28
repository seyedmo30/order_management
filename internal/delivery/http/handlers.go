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

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"status": "Failed", "error": err.Error()})
	}
	err := h.usecase.CreateOrder(c.Request().Context(), dto.CreateOrderUsecaseRequest{BaseCreateOrderRequest: req.BaseCreateOrderRequest})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusAccepted, echo.Map{"message": "order created successfully"})
}

func (h *OrderHandler) GetOrders(c echo.Context) error {
	// Retrieve the order ID from the URL parameter
	orderID := c.Param("order_id")
	if orderID == "" {
		// Handle case where order_id is missing
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "order_id is required"})
	}

	// Call the use case to fetch the order details
	order, err := h.usecase.GetOrder(c.Request().Context(), orderID)
	if err != nil {
		// Handle errors (e.g., order not found, internal error)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Return the order details in the response
	return c.JSON(http.StatusOK, echo.Map{
		"order_id":        order.OrderID,
		"priority":        order.Priority,
		"processing_time": order.ProcessingTime,
		"status":          order.Status,
	})
}
