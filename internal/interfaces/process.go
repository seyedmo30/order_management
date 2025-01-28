package interfaces

import "context"

type Process interface {
	ProcessOrder(ctx context.Context, ProcessingTime int ) (status string , err error)
}
