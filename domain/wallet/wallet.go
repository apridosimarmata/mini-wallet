package wallet

import (
	"context"
	"errors"
	"mini-wallet/domain/common/response"
)

const (
	WALLET_TRANSACTION_DEPOSIT    = "deposit"
	WALLET_TRANSACTION_WITHDRAWAL = "withdrawal"
	WALLET_STATUS_DISABLED        = "disabled"
	WALLET_STATUS_ENABLED         = "enabled"
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
		return errors.New(response.WALLET_DISABLED_ERROR)
	}

	return nil
}

type WalletTransactionEntity struct {
	Id          string `json:"id"`
	WalletId    string `json:"wallet_id"`
	Amount      int    `json:"amount"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   string `json:"created_by"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	ReferenceId string `json:"reference_id"`
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

type WalletTransactionRequest struct {
	WalletId    string `json:"wallet_id"`
	Amount      int    `json:"amount"`
	ReferenceId string `json:"reference_id"`
}

type WalletCreationRequest struct {
	CustomerId string `schema:"customer_xid,required"`
}

func (payload *WalletCreationRequest) Validate() error {
	if len(payload.CustomerId) == 0 {
		return errors.New("invalid value")
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
	AddWalletBalance(ctx context.Context, req WalletTransactionRequest) (err error)
	ReduceWalletBalance(ctx context.Context, req WalletTransactionRequest) (err error)
	GetWalletTransactions(ctx context.Context, req GetWalletTransactionRequest) (res []WalletTransaction, err error)
}

type WalletRepository interface {
	GetCustomerWallet(ctx context.Context, customerId string) (res *Wallet, err error)
	GetWalletById(ctx context.Context, walletId string) (res *Wallet, err error)
	InsertWallet(ctx context.Context, wallet Wallet) (err error)
	UpdateWallet(ctx context.Context, wallet Wallet) (err error)
}
