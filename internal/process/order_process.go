package process

import (
	"context"
	"time"

	"github.com/seyedmo30/order_management/pkg"
)

type processUseCase struct {
}

func NewProcessUseCase() *processUseCase {

	return &processUseCase{}
}

// ProcessOrder simulates a process that uses a context timeout and processes for a given time.
func (u *processUseCase) ProcessOrder(ctx context.Context, processingTime int) (status string, err error) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
		if processingTime > 5 {
			return pkg.StatusOrderManagementFailed, nil
		}
		return pkg.StatusOrderManagementProcessed, nil
	}
}
