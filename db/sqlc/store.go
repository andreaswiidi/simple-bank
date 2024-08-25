package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/andreaswiidi/simple-bank/util"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer               TransfersHistory   `json:"transfer"`
	FromAccount            Account            `json:"from_account"`
	ToAccount              Account            `json:"to_account"`
	SenderTransaction      TransactionHistory `json:"sender_transaction"`
	BeneficiaryTransaction TransactionHistory `json:"beneficiary_transaction"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	timeEdit := time.Now()

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.SenderTransaction, err = q.CreateTransactionHistory(ctx, CreateTransactionHistoryParams{
			AccountID:         arg.FromAccountID,
			Amount:            -arg.Amount,
			TransactionType:   util.TRANSACTION_TYPE_TRANSFER,
			TransferHistoryID: sql.NullInt64{Int64: result.Transfer.ID, Valid: true},
		})
		if err != nil {
			return err
		}

		result.BeneficiaryTransaction, err = q.CreateTransactionHistory(ctx, CreateTransactionHistoryParams{
			AccountID:         arg.ToAccountID,
			Amount:            arg.Amount,
			TransactionType:   util.TRANSACTION_TYPE_TRANSFER,
			TransferHistoryID: sql.NullInt64{Int64: result.Transfer.ID, Valid: true},
		})
		if err != nil {
			return err
		}

		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:        arg.FromAccountID,
			Balance:   account1.Balance - arg.Amount,
			UpdatedAt: sql.NullTime{Time: timeEdit, Valid: true},
		})
		if err != nil {
			return err
		}

		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:        arg.ToAccountID,
			Balance:   account2.Balance + arg.Amount,
			UpdatedAt: sql.NullTime{Time: timeEdit, Valid: true},
		})
		if err != nil {
			return err
		}

		return err
	})

	return result, err
}
