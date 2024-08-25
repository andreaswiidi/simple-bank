package db

import (
	"context"
	"fmt"
	"testing"

	sqlc "github.com/andreaswiidi/simple-bank/db/sqlc"
	"github.com/andreaswiidi/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := sqlc.NewStore(testDB)
	var err error
	var account1 sqlc.Account
	var account2 sqlc.Account
	account1, _, err = createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	account2, _, err = createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan sqlc.TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		txKey := struct{}{}
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, sqlc.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetDetilTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check transaction
		fromSender := result.SenderTransaction
		require.NotEmpty(t, fromSender)
		require.Equal(t, account1.ID, fromSender.AccountID)
		require.Equal(t, -amount, fromSender.Amount)
		require.Equal(t, util.TRANSACTION_TYPE_TRANSFER, fromSender.TransactionType)
		require.NotZero(t, fromSender.ID)
		require.NotZero(t, fromSender.CreatedAt)

		_, err = store.GetTransactionHistory(context.Background(), fromSender.ID)
		require.NoError(t, err)

		toBeneficiary := result.BeneficiaryTransaction
		require.NotEmpty(t, toBeneficiary)
		require.Equal(t, account2.ID, toBeneficiary.AccountID)
		require.Equal(t, amount, toBeneficiary.Amount)
		require.NotZero(t, toBeneficiary.ID)
		require.NotZero(t, toBeneficiary.CreatedAt)

		_, err = store.GetTransactionHistory(context.Background(), toBeneficiary.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		BeneficiaryAccount := result.ToAccount
		require.NotEmpty(t, BeneficiaryAccount)
		require.Equal(t, account2.ID, BeneficiaryAccount.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, BeneficiaryAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := BeneficiaryAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
