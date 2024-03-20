package auth

import (
	"context"
	"net/http"
)

type Token struct {
	Token    string  `json:"token"`
	WalletId *string `json:"wallet_id"`
}

type AuthUsecase interface {
	AuthorizeRequestMiddleware(next http.Handler) http.Handler
	InitUser(ctx context.Context, customerId string) (token string, err error)
}
