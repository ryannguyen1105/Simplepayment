package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	PaymentTx(ctx context.Context, arg PaymentTxParams) (PaymentTxResult, error)
}
type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type PaymentTxParams struct {
	FromWalletID int64 `json:"from_wallet_id"`
	ToWalletID   int64 `json:"to_wallet_id"`
	Amount       int64 `json:"amount"`
}

type PaymentTxResult struct {
	Payment    Payment `json:"payment"`
	FromWallet Wallet  `json:"from_wallet"`
	ToWallet   Wallet  `json:"to_wallet"`
	FromEntry  Entry   `json:"from_entry"`
	ToEntry    Entry   `json:"to_entry"`
}

func (store *SQLStore) PaymentTx(ctx context.Context, arg PaymentTxParams) (PaymentTxResult, error) {
	var result PaymentTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Payment, err = q.CreatePayment(ctx, CreatePaymentParams{
			FromWalletID: arg.FromWalletID,
			ToWalletID:   arg.ToWalletID,
			Amount:       arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			WalletID: arg.FromWalletID,
			Amount:   -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			WalletID: arg.ToWalletID,
			Amount:   arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromWalletID < arg.ToWalletID {
			result.FromWallet, result.ToWallet, err = addMoney(ctx, q, arg.FromWalletID, -arg.Amount, arg.ToWalletID, arg.Amount)
		} else {
			result.ToWallet, result.FromWallet, err = addMoney(ctx, q, arg.ToWalletID, arg.Amount, arg.FromWalletID, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	walletID1 int64,
	amount1 int64,
	walletID2 int64,
	amount2 int64,
) (wallet1 Wallet, wallet2 Wallet, err error) {
	wallet1, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
		ID:     walletID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	wallet2, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
		ID:     walletID2,
		Amount: amount2,
	})
	return
}
