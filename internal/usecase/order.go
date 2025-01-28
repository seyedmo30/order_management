// internal/usecase/order_usecase.go
package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/seyedmo30/order_management/internal/config"
	"github.com/seyedmo30/order_management/internal/dto"
	"github.com/seyedmo30/order_management/internal/interfaces"
	"github.com/seyedmo30/order_management/pkg"
)

// GetLogger initializes the global logger instance
var logger = pkg.GetLogger()

// orderUseCase is the concrete implementation of the OrderUseCase interface.
// It is responsible for creating orders and triggering the processing worker through a signal.
type orderUseCase struct {
	config             config.App
	repo               interfaces.OrderRepository
	process            interfaces.Process
	orderCreatedSignal chan struct{} // Signal to notify worker when an order is created
}

// NewOrderUseCase creates a new instance of orderUseCase.
// It initializes the order repository and the signal channel for worker triggering.
func NewOrderUseCase(config config.App, repo interfaces.OrderRepository, process interfaces.Process) *orderUseCase {
	// Create a buffered channel to send a signal when an order is created.
	// This acts as an "observer" to trigger the worker to process the newly created order.
	orderCreatedSignal := make(chan struct{}, 1)

	logger.Info("NewOrderUseCase instance created")
	return &orderUseCase{config: config, repo: repo, orderCreatedSignal: orderCreatedSignal, process: process}
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
		logger.Error("Failed to create order", "error", err)
		return err
	}

	// Log successful order creation
	logger.Info("Order created successfully", "orderID", req.OrderID)

	// Attempt to send a signal to notify the worker to process the order.
	// This operation is non-blocking; if the channel is full, it won't block the execution.
	select {
	case u.orderCreatedSignal <- struct{}{}:
		// Signal sent successfully, worker will process the order
		logger.Info("Signal sent to worker to process order", "orderID", req.OrderID)
	default:
		// If the channel is full, we simply don't send the signal to avoid deadlock
		// and continue with the execution.
		logger.Warn("Order creation signal channel full, no signal sent", "orderID", req.OrderID)
	}
	// Return response (just echoing the OrderID and success message here)
	return nil
}

// ProcessOrder handles worker creation and processing of orders.
func (u *orderUseCase) ProcessOrder(ctx context.Context) error {
	workerCount := u.config.WorkerCount
	var wg sync.WaitGroup
	wg.Add(workerCount)

	// Log the start of the worker pool
	logger.Info("Starting order processing workers", "workerCount", workerCount)

	// Start the workers
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					// Context canceled, stop worker
					logger.Info("Worker context canceled, stopping worker", "workerID", workerID)
					return
				case <-u.orderCreatedSignal:
					// Signal received, process an order
					logger.Info("Worker received signal to process order", "workerID", workerID)
					if err := u.processSingleOrder(ctx); err != nil {
						// Handle error if needed (e.g., log it)
						logger.Error("Error processing order", "workerID", workerID, "error", err)
					}
				}
			}
		}(i)
	}

	// Wait for all workers to finish (only ends if the context is canceled)
	wg.Wait()
	logger.Info("All workers finished processing orders")
	return nil
}

// ListAggregateOrderReport logs the order status counts every 2 seconds.
func (u *orderUseCase) ListAggregateOrderReport(ctx context.Context) error {
	reportInterval := u.config.ReportInterval
	ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer ticker.Stop()

	// Loop that runs every 2 seconds to log order status counts
	for {
		select {
		case <-ticker.C:
			counts, err := u.repo.ListAggregateOrderReport(ctx)
			if err != nil {
				log.Printf("Error fetching order status report: %v", err)
				continue
			}
			log.Printf("Order status report: \n%v \n", counts)
		case <-ctx.Done():
			// Graceful shutdown if context is canceled
			return ctx.Err()
		}
	}
}

func (u *orderUseCase) GetOrder(ctx context.Context, orderID string) (res dto.GetOrderUsecaseResponse, err error) {
	// Call the repository method to fetch the order by ID
	repoRes, err := u.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return dto.GetOrderUsecaseResponse{}, fmt.Errorf("failed to get order: %w", err)
	}

	// Map repository response to use case response (you can add more transformation logic here if needed)
	res = dto.GetOrderUsecaseResponse{
		BaseOrder: repoRes.BaseOrder,
	}

	return res, nil
}

// processSingleOrder handles fetching, locking, and processing a single order.
func (u *orderUseCase) processSingleOrder(ctx context.Context) error {
	// Get the next high-priority ready order
	getNextHighPriorityReadyOrder, err := u.repo.GetNextHighPriorityReadyOrder(ctx)
	if err != nil {
		logger.Error("Failed to fetch next high-priority order", "error", err)
		return err // Handle error (e.g., log or return)
	}

	// Log the order being processed
	logger.Info("Processing order", "orderID", getNextHighPriorityReadyOrder.OrderID)

	select {
	case u.orderCreatedSignal <- struct{}{}:
		// Signal sent successfully, worker will process the order
		logger.Info("Signal sent to worker to process order", "orderID", getNextHighPriorityReadyOrder.OrderID)
	default:
		// If the channel is full, we simply don't send the signal to avoid deadlock
		// and continue with the execution.
		logger.Warn("Order creation signal channel full, no signal sent", "orderID", getNextHighPriorityReadyOrder.OrderID)
	}

	// Lock the order optimistically
	lockOrderOptimisticRepositoryRequest := dto.LockOrderOptimisticRepositoryRequest{
		BaseOrder: getNextHighPriorityReadyOrder.BaseOrder,
	}
	if err := u.repo.LockOrderOptimistic(ctx, lockOrderOptimisticRepositoryRequest); err != nil {
		logger.Error("Failed to lock order optimistically", "orderID", getNextHighPriorityReadyOrder.OrderID, "error", err)
		return err // Handle error (e.g., log or return)
	}

	// Process the order using the process service
	status, err := u.process.ProcessOrder(ctx, *getNextHighPriorityReadyOrder.ProcessingTime)
	if err != nil {
		// Handle process failure
		logger.Error("Failed to process order", "orderID", getNextHighPriorityReadyOrder.OrderID, "error", err)
		return err
	}

	// Log order status after processing
	logger.Info("Order processed successfully", "orderID", getNextHighPriorityReadyOrder.OrderID, "status", status)

	// Update the order status based on the processing result
	updateRequest := dto.UpdateOrderByIDRepositoryRequest{
		BaseOrder: dto.BaseOrder{
			OrderID: getNextHighPriorityReadyOrder.OrderID,
			Status:  &status,
			Lock:    new(bool),
		},
	}
	if err := u.repo.UpdateOrderByID(ctx, updateRequest); err != nil {
		logger.Error("Failed to update order status", "orderID", getNextHighPriorityReadyOrder.OrderID, "error", err)
		return err // Handle error (e.g., log or return)
	}

	return nil
}
