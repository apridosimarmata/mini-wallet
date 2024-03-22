package worker

import (
	"context"
)

type WorkerUsecase interface {
	SubscribeWalletTransaction(ctx context.Context) (err error)
}
