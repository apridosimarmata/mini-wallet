package auth

import (
	"context"
	"crypto/sha1"
	"log"
	"mini-wallet/domain"
	"mini-wallet/domain/auth"
	"mini-wallet/domain/wallet"
	"mini-wallet/infrastructure"
	"net/http"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

type authUsecase struct {
	APIKey           string
	walletRepository wallet.WalletRepository
}

func NewAuthUsecase(repositories domain.Repositories) auth.AuthUsecase {
	return &authUsecase{
		APIKey: "any",
	}
}

// InitUser should create a new Wallet and set customer_xid provided as the owner
// if there is already a Wallet of this customer -> provide a new token with assumption that:
// 1. Previous issued token is already expired (being deleted from Redis depends on its TTL)
// 2. The request is being made from different device/client
func (usecase *authUsecase) InitUser(ctx context.Context, customerId string) (token string, err error) {
	customerWallet, err := usecase.walletRepository.GetCustomerWallet(ctx, customerId)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.GetCustomerWallet()")
		return "", err
	}

	if customerWallet == nil {
		walletId, err := uuid.NewV6()
		if err != nil {
			infrastructure.Log("got error on uuid.NewV6()")
			return "", err
		}

		err = usecase.walletRepository.InsertWallet(ctx, wallet.Wallet{
			WalletId: walletId.String(),
			OwnedBy:  customerId,
			Balance:  0,
			Status:   wallet.WALLET_STATUS_DISABLED,
		})
		if err != nil {
			log.Default().Printf("got error on usecase.walletRepository.InsertWallet()")
			return "", err
		}
	}

	return usecase.generateToken(customerId), nil
}

func (usecase *authUsecase) generateToken(customerId string) string {
	// avoid generating the same token over and over again -> add time
	var sha = sha1.New()
	customerId = customerId + time.Now().String()
	sha.Write([]byte(customerId))

	var token = sha.Sum(nil)

	return string(token)
}

func (usecase *authUsecase) AuthorizeRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 || (authHeader[1] != usecase.APIKey) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
