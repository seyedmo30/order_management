package interfaces

import (
	"context"

	"github.com/seyedmo30/order_management/internal/dto"
)

// OrderUseCase defines the methods that the use case layer will implement.
type OrderUseCase interface {
	CreateOrder(ctx context.Context, params dto.CreateOrderUsecaseRequest) error
	ProcessOrder(ctx context.Context) error
	ListAggregateOrderReport(ctx context.Context) error
	GetOrder(ctx context.Context, orderID string) (res dto.GetOrderUsecaseResponse, err error)
}
