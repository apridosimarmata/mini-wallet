package presentation

import (
	"context"
	"mini-wallet/app/auth"
	"mini-wallet/app/wallet"

	"mini-wallet/domain"
	"mini-wallet/infrastructure"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"
)

func InitServer() chi.Router {
	ctx := context.Background()
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	config := infrastructure.GetConfig()

	postgresDb := infrastructure.NewPostgresConn(config)
	redisClient := infrastructure.NewRedisClient(ctx, config)
	cache := infrastructure.NewCache(redisClient)

	pool := goredis.NewPool(&redisClient)
	mutexProvider := redsync.New(pool)

	// redsync for distributed mutual exclusion

	repositories := domain.Repositories{
		WalletRepository: wallet.NewWalletRepository(postgresDb, cache),
	}

	usecases := domain.Usecases{
		AuthUsecase:   auth.NewAuthUsecase(repositories),
		WalletUsecase: wallet.NewWalletUsecase(repositories, cache, mutexProvider),
	}

	wallet.SetWalletHandler(router, usecases)

	// in terms of authorization, a token should not be a forever-lived value
	// provided a /refresh endpoint to get fresh token
	auth.SetAuthHandler(router, usecases)

	return router
}

func StopServer() {

}
