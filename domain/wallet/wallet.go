package wallet

import (
	"context"
	"errors"
)

const (
	WALLET_TRANSACTION_DEPOSIT    = "deposit"
	WALLET_TRANSACTION_WITHDRAWAL = "withdrawal"
	WALLET_STATUS_DISABLED        = "disabled"
	WALLET_STATUS_ENABLED         = "enabled"
)

type Wallet struct {
	WalletId   string  `json:"wallet_id"`
	OwnedBy    string  `json:"owned_by"` // customer_xid on wallet creation
	EnabledAt  *string `json:"enabled_at"`
	DisabledAt *string `json:"disabled_at"`
	Balance    int     `json:"balance"`
	Status     string  `json:"status"`
}

type WalletTransaction struct {
	WalletId  string `json:"wallet_id"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
	Type      string `json:"type"`
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
	CreateWallet(ctx context.Context, req WalletCreationRequest) (err error)
	EnableWallet(ctx context.Context, token string, walletId string) (err error)
	DisableWallet(ctx context.Context, token string, walletId string) (err error)
	AddWalletBalance(ctx context.Context, req WalletTransactionRequest) (err error)
	ReduceWalletBalance(ctx context.Context, req WalletTransactionRequest) (err error)
	GetWalletTransactions(ctx context.Context, req GetWalletTransactionRequest) (res []WalletTransaction, err error)
}

type WalletRepository interface {
	GetCustomerWallet(ctx context.Context, customerId string) (res *Wallet, err error)
	InsertWallet(ctx context.Context, wallet Wallet) (err error)
	UpdateWallet(ctx context.Context, wallet Wallet) (err error)
}
