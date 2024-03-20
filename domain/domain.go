package domain

import (
	"mini-wallet/domain/auth"
	"mini-wallet/domain/wallet"
)

type Repositories struct {
	WalletRepository wallet.WalletRepository
}

type Usecases struct {
	WalletUsecase wallet.WalletUsecase
	AuthUsecase   auth.AuthUsecase
}
