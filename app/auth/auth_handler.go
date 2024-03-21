package auth

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/auth"
	"mini-wallet/domain/wallet"
	"net/http"

	"mini-wallet/domain/common/response"

	"github.com/go-chi/chi/v5"
)

type authHandler struct {
	authUsecase auth.AuthUsecase
}

func SetAuthHandler(router *chi.Mux, usecases domain.Usecases) {
	authHandler := authHandler{
		authUsecase: usecases.AuthUsecase,
	}

	router.Route("/api/v1/", func(r chi.Router) {
		r.Post("/init", authHandler.InitUser)

	})
}

func (authHandler *authHandler) InitUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := wallet.WalletCreationRequest{}
	resp := &response.Response[auth.Token]{}

	req.CustomerId = r.FormValue("customer_xid")
	err := req.Validate()
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Fail()
		errResp.WriteResponse(w)
		return
	}

	resp, err = authHandler.authUsecase.InitUser(ctx, req.CustomerId)
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

	resp.StatusCode = http.StatusOK
	resp.Status = response.STATUS_SUCCESS
	resp.WriteResponse(w)
}
