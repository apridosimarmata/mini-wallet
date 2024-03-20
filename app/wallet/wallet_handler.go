package wallet

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/auth"
	"mini-wallet/domain/wallet"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type walletHandler struct {
	walletUsecase wallet.WalletUsecase
	authUsecase   auth.AuthUsecase
}

func SetWalletHandler(router *chi.Mux, usecases domain.Usecases) {
	walletHandler := walletHandler{
		walletUsecase: usecases.WalletUsecase,
		authUsecase:   usecases.AuthUsecase,
	}

	router.Route("/api/v1/wallet", func(r chi.Router) {
		r.Use(usecases.AuthUsecase.AuthorizeRequestMiddleware)

		// GET
		r.Get("/wallet", walletHandler.GetWalletBalance)
		r.Get("/wallet/transactions", walletHandler.GetWalletTransactions)

		// POST
		r.Post("/wallet", walletHandler.ChangeWalletStatus)
		r.Post("/deposits", walletHandler.AddWalletBalance)
		r.Post("/withdrawals", walletHandler.WithdrawFromWallet)

		// PATCH
		r.Patch("/wallet", walletHandler.DisableWallet)

	})

}

func (handler *walletHandler) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	err := handler.walletUsecase.AddWalletBalance(ctx, wallet.WalletTransactionRequest{})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return
	}
}

func (handler *walletHandler) GetWalletTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	err := handler.walletUsecase.AddWalletBalance(ctx, wallet.WalletTransactionRequest{})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return
	}
}

func (handler *walletHandler) ChangeWalletStatus(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	err := handler.walletUsecase.AddWalletBalance(ctx, wallet.WalletTransactionRequest{})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return
	}
}

func (handler *walletHandler) AddWalletBalance(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	err := handler.walletUsecase.AddWalletBalance(ctx, wallet.WalletTransactionRequest{})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return
	}
}

func (handler *walletHandler) WithdrawFromWallet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	err := handler.walletUsecase.AddWalletBalance(ctx, wallet.WalletTransactionRequest{})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return
	}
}

func (handler *walletHandler) DisableWallet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	err := handler.walletUsecase.AddWalletBalance(ctx, wallet.WalletTransactionRequest{})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return
	}
}
