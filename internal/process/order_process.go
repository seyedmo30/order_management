package process

import (
	"context"
	"time"

	"github.com/seyedmo30/order_management/internal/config"
	"github.com/seyedmo30/order_management/pkg"
)

type processUseCase struct {
	config config.App
}

func NewProcessUseCase(config config.App) *processUseCase {
	return &processUseCase{config: config}
}

// ProcessOrder simulates a process that uses a context timeout and processes for a given time.
func (u *processUseCase) ProcessOrder(ctx context.Context, processingTime int) (status string, err error) {
	timeout := u.config.OrderProcessTimeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel() // Ensure the context is canceled to release resources

	// Create a timer for the processing time
	processingTimer := time.NewTimer(time.Duration(processingTime) * time.Second)
	defer processingTimer.Stop() // Ensure the timer is stopped to release resources

	select {
	case <-ctx.Done():
		// Context canceled or timed out
		return pkg.StatusOrderManagementFailed, nil
	case <-processingTimer.C:
		// Processing finished within the context deadline
		if processingTime > timeout {
			return pkg.StatusOrderManagementFailed, nil
		}
		return pkg.StatusOrderManagementProcessed, nil
	}
}
