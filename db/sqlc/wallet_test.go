package db

import (
	"context"
	"testing"
	"time"

	"github.com/ryannguyen1105/Simplepayment/util"
	"github.com/stretchr/testify/require"
)

func createRandomWallet(t *testing.T, userID int64) Wallet {
	t.Helper()
	arg := CreateWalletParams{
		UserID:  userID,
		Balance: util.RandomMoney(),
	}
	wallet, err := testQueries.CreateWallet(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	require.Equal(t, arg.UserID, wallet.UserID)
	require.Equal(t, arg.Balance, wallet.Balance)

	require.NotZero(t, wallet.CreatedAt)

	return wallet

}

func TestCreateWallet(t *testing.T) {
	user := createRandomUser(t)
	createRandomWallet(t, user.ID)
}

func TestGetWallet(t *testing.T) {
	user := createRandomUser(t)

	wallet1 := createRandomWallet(t, user.ID)
	wallet2, err := testQueries.GetWalletByID(context.Background(), wallet1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	require.Equal(t, wallet1.ID, wallet2.ID)
	require.Equal(t, user.ID, wallet2.UserID)
	require.Equal(t, wallet1.Balance, wallet2.Balance)

	require.WithinDuration(t, wallet1.CreatedAt, wallet2.CreatedAt, time.Second)
}

func TestUpdateWallet(t *testing.T) {
	user := createRandomUser(t)
	wallet1 := createRandomWallet(t, user.ID)

	arg := UpdateWalletBalanceParams{
		ID:      wallet1.ID,
		Balance: util.RandomMoney(),
	}
	wallet2, err := testQueries.UpdateWalletBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	require.Equal(t, wallet1.ID, wallet2.ID)
	require.Equal(t, wallet1.UserID, wallet2.UserID)
	require.Equal(t, arg.Balance, wallet2.Balance)
	require.WithinDuration(t, wallet1.CreatedAt, wallet2.CreatedAt, time.Second)

}

func TestListWallets(t *testing.T) {
	for i := 0; i < 10; i++ {
		user := createRandomUser(t)
		createRandomWallet(t, user.ID)
	}
	arg := ListWalletsParams{
		Limit:  5,
		Offset: 5,
	}
	wallets, err := testQueries.ListWallets(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, wallets, 5)

	for _, wallet := range wallets {
		require.NotEmpty(t, wallet)
		require.NotZero(t, wallet.UserID)
	}
}
