package interfaces

import (
	"context"

	"github.com/seyedmo30/order_management/internal/dto"
)

// Repository is an interface that defines the methods for interacting with the repository.
// OrderRepository is an interface that defines the methods for interacting with the order repository.
type OrderRepository interface {
	CreateOrder(ctx context.Context, params dto.CreatOrderRepositoryRequest) (err error)

	GetOrderByID(ctx context.Context, orderID string) (res dto.GetOrderByIDRepositoryResponse, err error)

	GetNextHighPriorityReadyOrder(ctx context.Context) (res dto.GetNextHighPriorityReadyOrderRepositoryResponse, err error)

	UpdateOrderByID(ctx context.Context, params dto.UpdateOrderByIDRepositoryRequest) (err error)

	ListAggregateOrderReport(ctx context.Context) (counts map[string]int, err error)
}
