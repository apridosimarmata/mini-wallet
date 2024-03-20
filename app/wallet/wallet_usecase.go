package wallet

import (
	"context"
	"fmt"
	"mini-wallet/domain"
	"mini-wallet/domain/wallet"
	"mini-wallet/infrastructure"

	"github.com/go-redsync/redsync/v4"
)

const (
	mutexKey = "wallet-%s"
)

type walletUsecase struct {
	walletRepository wallet.WalletRepository
	cache            infrastructure.Cache
	mutexProvider    *redsync.Redsync
}

func NewWalletUsecase(repositories domain.Repositories, cache infrastructure.Cache, mutexProvider *redsync.Redsync) wallet.WalletUsecase {
	return &walletUsecase{
		walletRepository: repositories.WalletRepository,
		cache:            cache,
		mutexProvider:    mutexProvider,
	}
}

func (usecase *walletUsecase) CreateWallet(ctx context.Context, req wallet.WalletCreationRequest) (err error) {

	return
}

func (usecase *walletUsecase) EnableWallet(ctx context.Context, token string, walletId string) (err error) {

	return
}

func (usecase *walletUsecase) DisableWallet(ctx context.Context, token string, walletId string) (err error) {

	return
}

func (usecase *walletUsecase) AddWalletBalance(ctx context.Context, req wallet.WalletTransactionRequest) (err error) {

	if err := usecase.getWalletLock("anty"); err != nil {
		infrastructure.Log("got error on usecase.getWalletLock()")
		return err
	}

	if err := usecase.relaseWalletLock("anty"); err != nil {
		infrastructure.Log("got error on usecase.relaseWalletLock()")
		return err
	}

	return
}

func (usecase *walletUsecase) ReduceWalletBalance(ctx context.Context, req wallet.WalletTransactionRequest) (err error) {

	if err := usecase.getWalletLock("anty"); err != nil {
		infrastructure.Log("got error on usecase.getWalletLock()")
		return err
	}

	if err := usecase.relaseWalletLock("anty"); err != nil {
		infrastructure.Log("got error on usecase.relaseWalletLock()")
		return err
	}

	return
}

func (usecase *walletUsecase) GetWalletTransactions(ctx context.Context, req wallet.GetWalletTransactionRequest) (res []wallet.WalletTransaction, err error) {

	return
}

func (usecase *walletUsecase) getWalletLock(walletId string) error {
	walletMutexKey := fmt.Sprintf(mutexKey, walletId)
	walletMutex := usecase.mutexProvider.NewMutex(walletMutexKey)

	if err := walletMutex.Lock(); err != nil {
		return err
	}

	return nil
}

func (usecase *walletUsecase) relaseWalletLock(walletId string) error {
	walletMutexKey := fmt.Sprintf(mutexKey, walletId)
	walletMutex := usecase.mutexProvider.NewMutex(walletMutexKey)

	if ok, err := walletMutex.Unlock(); !ok || err != nil {
		return err
	}

	return nil
}
