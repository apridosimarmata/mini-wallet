package wallet

import (
	"context"
	"database/sql"
	"mini-wallet/domain/wallet"
	"mini-wallet/infrastructure"

	sq "github.com/Masterminds/squirrel"

	"gorm.io/gorm"
)

type walletRepository struct {
	db    *gorm.DB
	cache infrastructure.Cache
}

func NewWalletRepository(db *gorm.DB, cache infrastructure.Cache) wallet.WalletRepository {
	return &walletRepository{
		db:    db,
		cache: cache,
	}
}

func (walletRepository *walletRepository) GetWalletById(ctx context.Context, walletId string) (res *wallet.Wallet, err error) {
	builder := sq.Select("*").From("ms_wallet").Where(sq.Eq{"id": walletId})
	qry, args, err := builder.ToSql()
	if err != nil {
		return res, err
	}

	err = walletRepository.db.WithContext(ctx).Raw(qry, args...).Scan(&res).Error
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return
}

func (walletRepository *walletRepository) GetCustomerWallet(ctx context.Context, customerId string) (res *wallet.Wallet, err error) {
	builder := sq.Select("*").From("ms_wallet").Where(sq.Eq{"owned_by": customerId})
	qry, args, err := builder.ToSql()
	if err != nil {
		return res, err
	}

	err = walletRepository.db.WithContext(ctx).Raw(qry, args...).Scan(&res).Error
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return
}

func (walletRepository *walletRepository) InsertWallet(ctx context.Context, wallet wallet.Wallet) (err error) {
	err = walletRepository.db.WithContext(ctx).Table("ms_wallet").Create(wallet).Error
	if err != nil {
		return err
	}

	return nil
}

func (walletRepository *walletRepository) UpdateWallet(ctx context.Context, wallet wallet.Wallet) (err error) {
	err = walletRepository.db.WithContext(ctx).Table("ms_wallet").UpdateColumns(wallet).Error
	if err != nil {
		return err
	}
	return
}
