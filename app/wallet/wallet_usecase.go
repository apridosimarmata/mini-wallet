package wallet

import (
	"context"
	"errors"
	"fmt"
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/wallet"
	"mini-wallet/infrastructure"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/google/uuid"
)

const (
	mutexKey = "wallet-%s"
)

type walletUsecase struct {
	walletRepository wallet.WalletRepository
	cache            infrastructure.Cache
	config           infrastructure.Config
	mutexProvider    *redsync.Redsync
}

func NewWalletUsecase(
	repositories domain.Repositories,
	cache infrastructure.Cache,
	mutexProvider *redsync.Redsync,
	config infrastructure.Config) wallet.WalletUsecase {
	return &walletUsecase{
		walletRepository: repositories.WalletRepository,
		cache:            cache,
		mutexProvider:    mutexProvider,
		config:           config,
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
		return nil, errors.New(response.ERROR_WALLET_DISABLED)
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
	nowString := time.Now().Format(time.RFC3339)
	walletResult.EnabledAt = &nowString
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
	walletResult.EnabledAt = nil
	err = usecase.walletRepository.UpdateWallet(ctx, *walletResult)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.UpdateWallet() - EnableWallet")
		return nil, err
	}

	return &response.Response[wallet.Wallet]{
		Data: walletResult,
	}, nil
}

func (usecase *walletUsecase) CreateWalletTransaction(ctx context.Context, req wallet.WalletTransactionRequest) (res *response.Response[wallet.Wallet], err error) {
	var successResponse *response.Response[wallet.Wallet]
	var walletLock *redsync.Mutex

	walletResult, err := usecase.walletRepository.GetWalletById(ctx, req.WalletId)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.GetWalletById() - CreateWalletTransaction")
		return nil, err
	}

	if walletResult == nil {
		return nil, errors.New("wallet not found")
	}

	if err = walletResult.ValidateWalletStatus(); err != nil {
		return nil, errors.New(response.ERROR_WALLET_DISABLED)
	}

	// check if reference id already used before
	walletTransaction, err := usecase.walletRepository.GetWalletTransactionByReferenceId(ctx, req.ReferenceId)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.GetWalletTransactionByReferenceId() - CreateWalletTransaction")
		return nil, err
	}

	if walletTransaction != nil {
		return nil, errors.New("reference id already used")
	}

	successResponse = &response.Response[wallet.Wallet]{
		Data: walletResult,
	}

	// err = usecase.cache.Publish(ctx, usecase.config.WALLET_TRANSACTION_CHANNEL, req)
	// if err != nil {
	// 	infrastructure.Log("got error on usecase.cache.Publish() - CreateWalletTransaction")
	// 	return nil, err
	// }

	transactionId, err := uuid.NewV6()
	if err != nil {
		infrastructure.Log("got error on uuid.NewV6()")
		return nil, err
	}

	transactionEntity := wallet.WalletTransactionEntity{
		Id:          transactionId.String(),
		WalletId:    walletResult.Id,
		Amount:      req.Amount,
		CreatedAt:   time.Now().Format(time.RFC3339),
		CreatedBy:   walletResult.OwnedBy,
		Status:      wallet.WALLET_TRANSACTION_STATUS_SUCCESS,
		ReferenceId: req.ReferenceId,
	}

	if walletLock, err = usecase.getWalletLock(walletResult.Id); err != nil {
		infrastructure.Log("got error on usecase.getWalletLock()")
		return nil, errors.New("another process maybe still modifying this wallet")
	}

	switch req.Type {
	case wallet.WALLET_TRANSACTION_DEPOSIT:
		walletResult.Balance += req.Amount
		transactionEntity.Type = wallet.WALLET_TRANSACTION_DEPOSIT

		err = usecase.walletRepository.CreateWalletTransaction(ctx, *walletResult, transactionEntity)
		if err != nil {
			infrastructure.Log("got error on usecase.walletRepository.CreateWalletDeposit() - CreateWalletTransaction")
			return nil, err
		}
	case wallet.WALLET_TRANSACTION_WITHDRAWAL:
		if walletResult.Balance < req.Amount {
			return nil, errors.New(response.ERROR_INSSUFICIENT_FUND)
		}

		walletResult.Balance -= req.Amount
		transactionEntity.Type = wallet.WALLET_TRANSACTION_WITHDRAWAL
		err = usecase.walletRepository.CreateWalletTransaction(ctx, *walletResult, transactionEntity)
		if err != nil {
			infrastructure.Log("got error on usecase.walletRepository.CreateWalletWithdrawal() - CreateWalletTransaction")
			return nil, err
		}
	}

	if ok, err := walletLock.Unlock(); !ok || err != nil {
		infrastructure.Log("got error on usecase.relaseWalletLock()")
		return nil, err
	}

	return successResponse, nil
}

func (usecase *walletUsecase) GetWalletTransactions(ctx context.Context, walletId string) (res *response.Response[[]wallet.WalletTransaction], err error) {
	walletTransactions, err := usecase.walletRepository.GetWalletTransactionsByWalletId(ctx, walletId, 1, 10)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.GetWalletTransactionByReferenceId() - CreateWalletTransaction")
		return nil, err
	}

	transactions := []wallet.WalletTransaction{}

	for _, walletTransaction := range walletTransactions {
		switch walletTransaction.Type {
		case wallet.WALLET_TRANSACTION_DEPOSIT:
			transactions = append(transactions, walletTransaction.ToDepositTransaction())
		case wallet.WALLET_TRANSACTION_WITHDRAWAL:
			transactions = append(transactions, walletTransaction.ToWithdrawalTransaction())
		}
	}

	return &response.Response[[]wallet.WalletTransaction]{
		Data: &transactions,
	}, nil
}

func (usecase *walletUsecase) getWalletLock(walletId string) (mutex *redsync.Mutex, err error) {
	walletMutexKey := fmt.Sprintf(mutexKey, walletId)
	walletMutex := usecase.mutexProvider.NewMutex(walletMutexKey)

	if err := walletMutex.Lock(); err != nil {
		return nil, err
	}

	return walletMutex, nil
}
