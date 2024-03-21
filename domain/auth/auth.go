package auth

import (
	"context"
	"mini-wallet/domain/common/response"
	"net/http"
)

type Token struct {
	Token string `json:"token"`
}

type AuthUsecase interface {
	AuthorizeRequestMiddleware(next http.Handler) http.Handler
	InitUser(ctx context.Context, customerId string) (token *response.Response[Token], err error)
}

type AuthRepository interface {
	AddToken(ctx context.Context, token string, walletId string) (err error)
	GetTokenWalletId(ctx context.Context, token string) (walletId string, err error)
}
