package auth

import (
	"context"
	"mini-wallet/domain/auth"
	"mini-wallet/infrastructure"
)

type authRepository struct {
	cache infrastructure.Cache
}

func NewAuthRepository(cache infrastructure.Cache) auth.AuthRepository {
	return &authRepository{
		cache: cache,
	}
}

func (authRepository *authRepository) AddToken(ctx context.Context, token string, walletId string) (err error) {
	err = authRepository.cache.SetString(ctx, token, walletId, 6000) // ttl can be defined on .env
	if err != nil {
		infrastructure.Log("got error on authRepository.cache.SetString() - AddToken")
		return nil
	}

	return nil
}
func (authRepository *authRepository) GetTokenWalletId(ctx context.Context, token string) (walletId string, err error) {
	walletId, err = authRepository.cache.GetString(ctx, token)
	if err != nil {
		infrastructure.Log("got error on authRepository.cache.GetString() - GetTokenWalletId")
		return "", err
	}

	return walletId, nil
}
