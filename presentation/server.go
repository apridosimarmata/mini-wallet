package presentation

import (
	"context"
	"fmt"
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

	fmt.Println(fmt.Sprintf("config got: %v", config))

	postgresDb := infrastructure.NewPostgresConn(config)
	redisClient := infrastructure.NewRedisClient(ctx, config)
	cache := infrastructure.NewCache(redisClient)

	// redsync for distributed mutual exclusion
	pool := goredis.NewPool(&redisClient)
	mutexProvider := redsync.New(pool)

	repositories := domain.Repositories{
		WalletRepository: wallet.NewWalletRepository(postgresDb, cache),
		AuthRepository:   auth.NewAuthRepository(cache),
	}

	usecases := domain.Usecases{
		AuthUsecase:   auth.NewAuthUsecase(repositories),
		WalletUsecase: wallet.NewWalletUsecase(repositories, cache, mutexProvider, config),
	}

	wallet.SetWalletHandler(router, usecases)
	// in terms of authorization, a token should not be a forever-lived value
	// provided a /refresh endpoint to get fresh token
	auth.SetAuthHandler(router, usecases)

	// 1.
	// starting worker to listen wallet transaction
	// when getting balance, the requirement expecting a delay (5 seconds at max).
	// it can be occured when the wallet details are being cached
	// or most likely there is a messaging mechanism for each transaction
	/*** the subscriber -> not implemented. ***/
	// workerUsecase := worker.NewWorkerUsecase(redisClient, repositories, config)
	// go workerUsecase.SubscribeWalletTransaction(ctx)

	// 2.
	// I changed my mind, the delay must be caused by cache instead of messaging.
	// why? the response returned to user already stating that the deposit/withdrawal request:
	// success / error
	// the balance must be immediately updated with usecase instead of messaging mechanism.

	// 3. [implemented]
	// but here's another thing.
	// what if the delay means two request is being sent almost at the same time:
	// req I -> deposit/withdraw
	// req II -> check balance sent when req I is still in process.
	/*** there is context with timeout while updating wallet balance ***/

	return router
}

func StopServer() {

}
