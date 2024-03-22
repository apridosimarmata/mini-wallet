package wallet

import (
	"context"
	"errors"
	"mini-wallet/domain/common/response"
)

const (
	WALLET_TRANSACTION_DEPOSIT        = "deposit"
	WALLET_TRANSACTION_WITHDRAWAL     = "withdrawal"
	WALLET_STATUS_DISABLED            = "disabled"
	WALLET_STATUS_ENABLED             = "enabled"
	WALLET_TRANSACTION_STATUS_SUCCESS = "success"
)

type Wallet struct {
	Id        string  `json:"id" gorm:"column:id"`
	OwnedBy   string  `json:"owned_by" gorm:"column:owned_by"` // customer_xid on wallet creation
	EnabledAt *string `json:"enabled_at" gorm:"column:enabled_at"`
	Balance   int     `json:"balance" gorm:"column:balance"`
	Status    string  `json:"status" gorm:"column:status"`
}

func (wallet *Wallet) ValidateWalletStatus() error {
	if wallet.Status != WALLET_STATUS_ENABLED {
		return errors.New(response.ERROR_WALLET_DISABLED)
	}

	return nil
}

type WalletTransactionEntity struct {
	Id          string `json:"id" gorm:"column:id"`
	WalletId    string `json:"wallet_id" gorm:"column:wallet_id"`
	Amount      int    `json:"amount" gorm:"column:amount"`
	CreatedAt   string `json:"created_at" gorm:"column:created_at"`
	CreatedBy   string `json:"created_by" gorm:"column:created_by"`
	Type        string `json:"type" gorm:"column:type"`
	Status      string `json:"status" gorm:"column:status"`
	ReferenceId string `json:"reference_id" gorm:"column:reference_id"`
}

type WalletTransaction struct {
	Id          string  `json:"id"`
	Amount      int     `json:"amount"`
	Status      string  `json:"status"`
	ReferenceId string  `json:"reference_id"`
	DepositedAt *string `json:"deposited_at,omitempty"`
	DepositedBy *string `json:"deposited_by,omitempty"`
	WithdrawnAt *string `json:"withdrawn_at,omitempty"`
	WithdrawnBy *string `json:"withdrawn_by,omitempty"`
}

func (walletTransaction *WalletTransactionEntity) ToWithdrawalTransaction() WalletTransaction {
	return WalletTransaction{
		Id:          walletTransaction.Id,
		Amount:      walletTransaction.Amount,
		Status:      walletTransaction.Status,
		ReferenceId: walletTransaction.ReferenceId,
		WithdrawnAt: &walletTransaction.CreatedAt,
		WithdrawnBy: &walletTransaction.CreatedBy,
	}
}

func (walletTransaction *WalletTransactionEntity) ToDepositTransaction() WalletTransaction {
	return WalletTransaction{
		Id:          walletTransaction.Id,
		Amount:      walletTransaction.Amount,
		Status:      walletTransaction.Status,
		ReferenceId: walletTransaction.ReferenceId,
		DepositedAt: &walletTransaction.CreatedAt,
		DepositedBy: &walletTransaction.CreatedBy,
	}
}

type WalletTransactionRequest struct {
	WalletId    string `json:"wallet_id"`
	Type        string `json:"type"`
	Amount      int    `json:"amount"`
	ReferenceId string `json:"reference_id"`
	Timestamp   int    `json:"timestamp"`
}

func (transactionRequest *WalletTransactionRequest) Validate() error {
	if transactionRequest.Amount <= 0 || len(transactionRequest.ReferenceId) == 0 {
		return errors.New(response.ERROR_BAD_REQUEST)
	}

	return nil
}

type WalletCreationRequest struct {
	CustomerId string `schema:"customer_xid,required"`
}

func (payload *WalletCreationRequest) Validate() error {
	if len(payload.CustomerId) == 0 {
		return errors.New(response.ERROR_BAD_REQUEST)
	}

	return nil
}

type GetWalletTransactionRequest struct {
	WalletId string  `json:"wallet_id"`
	Type     *string `json:"type"`
}

type WalletUsecase interface {
	EnableWallet(ctx context.Context, walletId string) (res *response.Response[Wallet], err error)
	DisableWallet(ctx context.Context, walletId string) (res *response.Response[Wallet], err error)
	GetWalletBalance(ctx context.Context, walletId string) (res *response.Response[Wallet], err error)
	CreateWalletTransaction(ctx context.Context, req WalletTransactionRequest) (res *response.Response[Wallet], err error)
	GetWalletTransactions(ctx context.Context, walletId string) (res *response.Response[[]WalletTransaction], err error)
}

type WalletRepository interface {
	GetCustomerWallet(ctx context.Context, customerId string) (res *Wallet, err error)
	GetWalletById(ctx context.Context, walletId string) (res *Wallet, err error)
	InsertWallet(ctx context.Context, wallet Wallet) (err error)
	UpdateWallet(ctx context.Context, wallet Wallet) (err error)
	CreateWalletTransaction(ctx context.Context, updatedWallet Wallet, walletTransaction WalletTransactionEntity) (err error)
	GetWalletTransactionByReferenceId(ctx context.Context, referenceId string) (res *WalletTransactionEntity, err error)
	GetWalletTransactionsByWalletId(ctx context.Context, walletId string, page int, size int) (res []WalletTransactionEntity, err error)
}
