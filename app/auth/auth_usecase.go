package auth

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"mini-wallet/domain"
	"mini-wallet/domain/auth"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/wallet"
	"mini-wallet/infrastructure"
	"net/http"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

type authUsecase struct {
	walletRepository wallet.WalletRepository
	authRepository   auth.AuthRepository
}

func NewAuthUsecase(repositories domain.Repositories) auth.AuthUsecase {
	return &authUsecase{
		walletRepository: repositories.WalletRepository,
		authRepository:   repositories.AuthRepository,
	}
}

// InitUser should create a new Wallet and set customer_xid provided as the owner
// if there is already a Wallet of this customer -> provide a new token with assumption that:
// 1. Previous issued token is already expired (being deleted from Redis depends on its TTL)
// 2. The request is being made from different device/client
func (usecase *authUsecase) InitUser(ctx context.Context, customerId string) (token *response.Response[auth.Token], err error) {
	var walletId string

	customerWallet, err := usecase.walletRepository.GetCustomerWallet(ctx, customerId)
	if err != nil {
		infrastructure.Log("got error on usecase.walletRepository.GetCustomerWallet()")
		return nil, err
	}

	if customerWallet == nil {
		walletId, err := uuid.NewV6()
		if err != nil {
			infrastructure.Log("got error on uuid.NewV6()")
			return nil, err
		}

		err = usecase.walletRepository.InsertWallet(ctx, wallet.Wallet{
			Id:      walletId.String(),
			OwnedBy: customerId,
			Balance: 0,
			Status:  wallet.WALLET_STATUS_DISABLED,
		})
		if err != nil {
			log.Default().Printf("got error on usecase.walletRepository.InsertWallet()")
			return nil, err
		}
	} else {
		walletId = customerWallet.Id
	}

	return &response.Response[auth.Token]{
		Data: &auth.Token{
			Token: usecase.generateToken(walletId),
		},
	}, nil
}

func (usecase *authUsecase) generateToken(walletId string) string {
	// avoid generating the same token over and over again -> add time
	var sha = sha1.New()
	walletIdTimestamp := walletId + time.Now().String()
	sha.Write([]byte(walletIdTimestamp))

	var token = sha.Sum(nil)

	usecase.authRepository.AddToken(context.Background(), fmt.Sprintf("%x", token), walletId)

	return fmt.Sprintf("%x", token)
}

func (usecase *authUsecase) AuthorizeRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		unauthorizedResp := response.Response[response.Error]{}

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			unauthorizedResp.Data = &response.Error{
				Error: "unauthorized",
			}
			unauthorizedResp.Fail()
			unauthorizedResp.WriteResponse(w)
			return
		}

		walletId, err := usecase.authRepository.GetTokenWalletId(ctx, authHeader[1])
		if err != nil || walletId == "" {
			unauthorizedResp.Data = &response.Error{
				Error: "unauthorized",
			}
			unauthorizedResp.Fail()
			unauthorizedResp.WriteResponse(w)
			return
		}

		ctx = context.WithValue(ctx, "walletId", walletId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
