// internal/usecase/order_usecase.go
package usecase

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/seyedmo30/order_management/internal/dto"
	"github.com/seyedmo30/order_management/internal/interfaces"
	"github.com/seyedmo30/order_management/pkg"
)

// orderUseCase is the concrete implementation of the OrderUseCase interface.
// It is responsible for creating orders and triggering the processing worker through a signal.
type orderUseCase struct {
	repo               interfaces.OrderRepository
	process            interfaces.Process
	orderCreatedSignal chan struct{} // Signal to notify worker when an order is created
}

// NewOrderUseCase creates a new instance of orderUseCase.
// It initializes the order repository and the signal channel for worker triggering.
func NewOrderUseCase(repo interfaces.OrderRepository, process interfaces.Process) *orderUseCase {
	// Create a buffered channel to send a signal when an order is created.
	// This acts as an "observer" to trigger the worker to process the newly created order.

	orderCreatedSignal := make(chan struct{}, 1)

	return &orderUseCase{repo: repo, orderCreatedSignal: orderCreatedSignal, process: process}
}

// CreateOrder processes the order and queues it.
// After creating the order in the repository, it triggers the worker to process it by sending a signal.
func (u *orderUseCase) CreateOrder(ctx context.Context, req dto.CreateOrderUsecaseRequest) error {

	// Construct the repository request for creating the order
	creatOrderRepositoryRequest := dto.CreatOrderRepositoryRequest{BaseOrder: dto.BaseOrder{
		ID:             uuid.NewString(),
		OrderID:        &req.OrderID,
		Priority:       &req.Priority,
		Status:         &pkg.StatusOrderManagementPending,
		ProcessingTime: &req.ProcessingTime,
	}}

	// Call the repository method to persist the order
	err := u.repo.CreateOrder(ctx, creatOrderRepositoryRequest)
	if err != nil {
		return err
	}

	// Attempt to send a signal to notify the worker to process the order.
	// This operation is non-blocking; if the channel is full, it won't block the execution.
	select {
	case u.orderCreatedSignal <- struct{}{}:
		// Signal sent successfully, worker will process the order
	default:
		// If the channel is full, we simply don't send the signal to avoid deadlock
		// and continue with the execution.
	}
	// Return response (just echoing the OrderID and success message here)
	return nil
}

// ProcessOrder handles worker creation and processing of orders.
func (u *orderUseCase) ProcessOrder(ctx context.Context) error {
	workerCount := 5 // Number of workers
	var wg sync.WaitGroup
	wg.Add(workerCount)

	// Start the workers
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					// Context canceled, stop worker
					return
				case <-u.orderCreatedSignal:
					// Signal received, process an order
					if err := u.processSingleOrder(ctx); err != nil {
						// Handle error if needed (e.g., log it)
					}
				}
			}false
		}(i)
	}

	// Wait for all workers to finish (only ends if the context is canceled)
	wg.Wait()
	return nil
}

// processSingleOrder handles fetching, locking, and processing a single order.
func (u *orderUseCase) processSingleOrder(ctx context.Context) error {
	// Get the next high-priority ready order
	getNextHighPriorityReadyOrder, err := u.repo.GetNextHighPriorityReadyOrder(ctx)
	if err != nil {
		return err // Handle error (e.g., log or return)
	}

	// Lock the order optimistically
	updateOrderByIDRepositoryRequest := dto.UpdateOrderByIDRepositoryRequest{
		BaseOrder: getNextHighPriorityReadyOrder.BaseOrder,
	}
	if err := u.repo.LockOrderOptimistic(ctx, updateOrderByIDRepositoryRequest); err != nil {
		return err // Handle error (e.g., log or return)
	}
	// Process the order using the process service
	status, err := u.process.ProcessOrder(ctx, *getNextHighPriorityReadyOrder.ProcessingTime)
	if err != nil {
		// Handle process failure
		return err
	}

	// Update the order status based on the processing result
	updateRequest := dto.UpdateOrderByIDRepositoryRequest{
		BaseOrder: dto.BaseOrder{
			OrderID: getNextHighPriorityReadyOrder.OrderID,
			Status:  &status,
			Lock:    new(bool),
		},
	}
	if err := u.repo.UpdateOrderByID(ctx, updateRequest); err != nil {
		return err // Handle error (e.g., log or return)
	}

	return nil
}
