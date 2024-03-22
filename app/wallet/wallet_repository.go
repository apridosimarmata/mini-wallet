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

func (walletRepository *walletRepository) GetWalletTransactionsByWalletId(ctx context.Context, walletId string, page int, size int) (res []wallet.WalletTransactionEntity, err error) {
	builder := sq.Select("*").From("tr_wallet_transaction").Where(sq.Eq{"wallet_id": walletId}).Limit(uint64(size)).Offset(uint64((page - 1) * size))
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

	return res, nil
}

func (walletRepository *walletRepository) GetWalletTransactionByReferenceId(ctx context.Context, referenceId string) (res *wallet.WalletTransactionEntity, err error) {
	builder := sq.Select("*").From("tr_wallet_transaction").Where(sq.Eq{"reference_id": referenceId})
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

func (walletRepository *walletRepository) CreateWalletTransaction(ctx context.Context, updatedWallet wallet.Wallet, walletTransaction wallet.WalletTransactionEntity) (err error) {
	tx := walletRepository.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = tx.WithContext(ctx).Table("tr_wallet_transaction").Create(walletTransaction).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.WithContext(ctx).Table("ms_wallet").Where("id", walletTransaction.WalletId).Update("balance", updatedWallet.Balance).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	res := tx.Commit()
	if err = res.Error; err != nil {
		return err
	}

	return nil
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
