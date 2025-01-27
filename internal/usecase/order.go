// internal/usecase/order_usecase.go
package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/seyedmo30/order_management/internal/dto"
	"github.com/seyedmo30/order_management/internal/interfaces"
	"github.com/seyedmo30/order_management/pkg"
)

// orderUseCase is the concrete implementation of the OrderUseCase interface.
type orderUseCase struct {
	repo interfaces.OrderRepository
}

// NewOrderUseCase creates a new instance of orderUseCase.
func NewOrderUseCase(repo interfaces.OrderRepository) *orderUseCase {
	return &orderUseCase{repo: repo}
}

// CreateOrder processes the order and queues it.
func (u *orderUseCase) CreateOrder(ctx context.Context, req dto.CreateOrderUsecaseRequest) error {

	creatOrderRepositoryRequest := dto.CreatOrderRepositoryRequest{BaseOrder: dto.BaseOrder{
		ID:             uuid.NewString(),
		OrderID:        &req.OrderID,
		Priority:       &req.Priority,
		Status:         &pkg.StatusOrderManagementPending,
		ProcessingTime: &req.ProcessingTime,
	}}

	
	// Simulate queuing the order (for the sake of example, just call repository method)
	err := u.repo.CreateOrder(ctx, creatOrderRepositoryRequest)
	if err != nil {
		return err
	}

	// Return response (just echoing the OrderID and success message here)
	return nil
}
