package wallet

import (
	"context"
	"errors"
	"fmt"
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
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

func (usecase *walletUsecase) GetWalletBalance(ctx context.Context, walletId string) (res *response.Response[wallet.Wallet], err error) {
	walletResult, err := usecase.walletRepository.GetWalletById(ctx, walletId)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.GetWalletById() - GetWalletBalance")
		return nil, err
	}

	if walletResult == nil {
		return nil, errors.New("wallet not found")
	}

	if err = walletResult.ValidateWalletStatus(); err != nil {
		return nil, errors.New(response.WALLET_DISABLED_ERROR)
	}

	return &response.Response[wallet.Wallet]{
		Data: walletResult,
	}, nil
}

func (usecase *walletUsecase) EnableWallet(ctx context.Context, walletId string) (res *response.Response[wallet.Wallet], err error) {
	walletResult, err := usecase.walletRepository.GetWalletById(ctx, walletId)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.GetWalletById() - EnableWallet")
		return nil, err
	}

	if walletResult == nil {
		return nil, errors.New("wallet not found")
	}

	walletResult.Status = wallet.WALLET_STATUS_ENABLED
	err = usecase.walletRepository.UpdateWallet(ctx, *walletResult)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.UpdateWallet() - EnableWallet")
		return nil, err
	}

	return &response.Response[wallet.Wallet]{
		Data: walletResult,
	}, nil
}

func (usecase *walletUsecase) DisableWallet(ctx context.Context, walletId string) (res *response.Response[wallet.Wallet], err error) {
	walletResult, err := usecase.walletRepository.GetWalletById(ctx, walletId)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.GetWalletById() - EnableWallet")
		return nil, err
	}

	if walletResult == nil {
		return nil, errors.New("wallet not found")
	}

	walletResult.Status = wallet.WALLET_STATUS_DISABLED
	err = usecase.walletRepository.UpdateWallet(ctx, *walletResult)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.UpdateWallet() - EnableWallet")
		return nil, err
	}

	return &response.Response[wallet.Wallet]{
		Data: walletResult,
	}, nil
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
