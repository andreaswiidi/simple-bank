package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/andreaswiidi/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount() (Account, CreateAccountParams, error) {
	arg := CreateAccountParams{
		Fullname: util.RandomOwner(),
		Username: util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	return account, arg, err
}

func TestCreateAccount(t *testing.T) {
	account, arg, err := createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Fullname, account.Fullname)
	require.Equal(t, arg.Username, account.Username)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	account1, _, err := createRandomAccount()
	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Fullname, account2.Fullname)
	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1, _, err := createRandomAccount()
	require.NoError(t, err)

	arg := UpdateAccountParams{
		ID:        account1.ID,
		Balance:   util.RandomMoney(),
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Fullname, account2.Fullname)
	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1, _, err := createRandomAccount()
	require.NoError(t, err)
	err = testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		tempAccount, _, err := createRandomAccount()
		require.NoError(t, err)
		lastAccount = tempAccount
	}

	arg := ListAccountsParams{
		Username: lastAccount.Username,
		Limit:    5,
		Offset:   0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Fullname, account.Fullname)
	}
}
