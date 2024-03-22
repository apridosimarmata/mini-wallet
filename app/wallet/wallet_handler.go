package wallet

import (
	"mini-wallet/domain"
	"mini-wallet/domain/auth"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/wallet"
	"net/http"
	"strconv"
	"time"

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
		r.Post("/deposits", walletHandler.CreateWalletDepositTransaction)
		r.Post("/withdrawals", walletHandler.CreateWalletWithdrawalTransaction)

		// PATCH
		r.Patch("/", walletHandler.DisableWallet)

	})

}

func (handler *walletHandler) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	walletId := r.Context().Value("walletId")
	resp := &response.Response[wallet.Wallet]{}

	result, err := handler.walletUsecase.GetWalletBalance(r.Context(), walletId.(string))
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
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
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	resp = result
	resp.Success(response.STATUS_SUCCESS, *resp.Data)
	resp.WriteResponse(w)
}

func (handler *walletHandler) DisableWallet(w http.ResponseWriter, r *http.Request) {
	walletId := r.Context().Value("walletId")

	result, err := handler.walletUsecase.DisableWallet(r.Context(), walletId.(string))
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	resp := &response.Response[wallet.Wallet]{}
	resp = result
	resp.Success(response.STATUS_SUCCESS, *resp.Data)
	resp.WriteResponse(w)
}

func (handler *walletHandler) CreateWalletDepositTransaction(w http.ResponseWriter, r *http.Request) {
	walletId := r.Context().Value("walletId")
	req := wallet.WalletTransactionRequest{
		WalletId: walletId.(string),
	}

	transactionAmount := r.FormValue("amount")
	referenceId := r.FormValue("reference_id")
	transactionAmountInt, err := strconv.Atoi(transactionAmount)
	if err != nil {
		transactionAmountInt = int(0)
	}

	req.Amount = transactionAmountInt
	req.ReferenceId = referenceId
	req.Type = wallet.WALLET_TRANSACTION_DEPOSIT
	req.Timestamp = int(time.Now().Unix())
	if err = req.Validate(); err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	result, err := handler.walletUsecase.CreateWalletTransaction(r.Context(), req)
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	resp := &response.Response[wallet.Wallet]{}
	resp = result
	resp.Success(response.STATUS_SUCCESS, *resp.Data)
	resp.WriteResponse(w)
}

func (handler *walletHandler) CreateWalletWithdrawalTransaction(w http.ResponseWriter, r *http.Request) {
	walletId := r.Context().Value("walletId")
	req := wallet.WalletTransactionRequest{
		WalletId: walletId.(string),
	}

	transactionAmount := r.FormValue("amount")
	referenceId := r.FormValue("reference_id")
	transactionAmountInt, err := strconv.Atoi(transactionAmount)
	if err != nil {
		transactionAmountInt = int(0)
	}

	req.Amount = transactionAmountInt
	req.ReferenceId = referenceId
	req.Type = wallet.WALLET_TRANSACTION_WITHDRAWAL
	req.Timestamp = int(time.Now().Unix())
	if err = req.Validate(); err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	result, err := handler.walletUsecase.CreateWalletTransaction(r.Context(), req)
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	resp := &response.Response[wallet.Wallet]{}
	resp = result
	resp.Success(response.STATUS_SUCCESS, *resp.Data)
	resp.WriteResponse(w)
}

func (handler *walletHandler) GetWalletTransactions(w http.ResponseWriter, r *http.Request) {
	walletId := r.Context().Value("walletId")

	result, err := handler.walletUsecase.GetWalletTransactions(r.Context(), walletId.(string))
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	resp := &response.Response[[]wallet.WalletTransaction]{}
	resp = result
	resp.Success(response.STATUS_SUCCESS, *resp.Data)
	resp.WriteResponse(w)
}
