package wallet

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/auth"
	"mini-wallet/domain/common/response"
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
		r.Get("/", walletHandler.GetWalletBalance)
		r.Get("/transactions", walletHandler.GetWalletTransactions)

		// POST
		r.Post("/", walletHandler.EnableWallet)
		r.Post("/deposits", walletHandler.AddWalletBalance)
		r.Post("/withdrawals", walletHandler.WithdrawFromWallet)

		// PATCH
		r.Patch("/", walletHandler.DisableWallet)

	})

}

func (handler *walletHandler) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	walletId := r.Context().Value("walletId")
	resp := &response.Response[wallet.Wallet]{}

	result, err := handler.walletUsecase.GetWalletBalance(r.Context(), walletId.(string))
	if err.Error() != response.WALLET_DISABLED_ERROR {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Fail()
		errResp.WriteResponse(w)
		return
	}
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error()
		errResp.WriteResponse(w)
		return
	}

	resp = result
	resp.Success(response.STATUS_SUCCESS, *resp.Data)
	resp.WriteResponse(w)
}

func (handler *walletHandler) EnableWallet(w http.ResponseWriter, r *http.Request) {
	walletId := r.Context().Value("walletId")
	resp := &response.Response[wallet.Wallet]{}

	result, err := handler.walletUsecase.EnableWallet(r.Context(), walletId.(string))
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error()
		errResp.WriteResponse(w)
		return
	}

	resp = result
	resp.Success(response.STATUS_SUCCESS, *resp.Data)
	resp.WriteResponse(w)
}

func (handler *walletHandler) DisableWallet(w http.ResponseWriter, r *http.Request) {
	walletId := r.Context().Value("walletId")
	resp := &response.Response[wallet.Wallet]{}

	result, err := handler.walletUsecase.DisableWallet(r.Context(), walletId.(string))
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error()
		errResp.WriteResponse(w)
		return
	}

	resp = result
	resp.Success(response.STATUS_SUCCESS, *resp.Data)
	resp.WriteResponse(w)
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
