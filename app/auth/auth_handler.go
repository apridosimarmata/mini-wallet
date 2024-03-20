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
	var resp response.Response[string]

	req.CustomerId = r.FormValue("customer_xid")
	err := req.Validate()
	if err != nil {
		resp.BadRequest("invalid payload")
		resp.WriteResponse(w)
		return
	}

	token, err := authHandler.authUsecase.InitUser(ctx, req.CustomerId)
	if err != nil {
		resp.InternalServerError(err.Error())
		resp.WriteResponse(w)
		return
	}

	resp.Data = token
	resp.WriteResponse(w)
}
