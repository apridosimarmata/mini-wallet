package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"mini-wallet/domain"
	"mini-wallet/domain/wallet"
	"mini-wallet/domain/worker"
	"mini-wallet/infrastructure"
	"time"

	"github.com/go-redis/redis"
)

type workerUsecase struct {
	cacheClient      redis.Client
	walletRepository wallet.WalletRepository
	config           infrastructure.Config
}

func NewWorkerUsecase(cacheClient redis.Client, repositories domain.Repositories, config infrastructure.Config) worker.WorkerUsecase {
	return &workerUsecase{
		cacheClient:      cacheClient,
		walletRepository: repositories.WalletRepository,
		config:           config,
	}
}

func (workerUsecase *workerUsecase) SubscribeWalletTransaction(ctx context.Context) (err error) {
	subscriber := workerUsecase.cacheClient.Subscribe(workerUsecase.config.WALLET_TRANSACTION_CHANNEL)

	transaction := wallet.WalletTransaction{}
	for {
		time.Sleep(time.Second * 2)
		msg, err := subscriber.ReceiveMessage()
		if msg == nil {
			continue
		}

		infrastructure.Log(fmt.Sprintf("message received - SubscribeWalletTransaction : %v", msg))

		if err != nil {
			infrastructure.Log("got error on usecase.walletRepository.GetWalletById() - SubscribeWalletTransaction")
			continue
		}

		if err := json.Unmarshal([]byte(msg.Payload), &transaction); err != nil {
			infrastructure.Log("got error on usecase.walletRepository.GetWalletById() - SubscribeWalletTransaction")
			continue
		}

		// ...
	}

}
