package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ryannguyen1105/Simplepayment/util"
	"github.com/stretchr/testify/require"
)

func createRandomWallet(t *testing.T) Wallet {
	t.Helper()
	arg := CreateWalletParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: "USD",
	}
	wallet, err := testQueries.CreateWallet(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	require.Equal(t, arg.Owner, wallet.Owner)
	require.Equal(t, arg.Balance, wallet.Balance)
	require.Equal(t, arg.Currency, wallet.Currency)

	require.NotZero(t, wallet.ID)
	require.NotZero(t, wallet.CreatedAt)

	return wallet

}

func TestCreateWallet(t *testing.T) {
	createRandomWallet(t)
}

func TestGetWallet(t *testing.T) {
	wallet1 := createRandomWallet(t)
	wallet2, err := testQueries.GetWallet(context.Background(), wallet1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	require.Equal(t, wallet1.ID, wallet2.ID)
	require.Equal(t, wallet1.Owner, wallet2.Owner)
	require.Equal(t, wallet1.Balance, wallet2.Balance)
	require.Equal(t, wallet1.Currency, wallet2.Currency)
	require.WithinDuration(t, wallet1.CreatedAt, wallet2.CreatedAt, time.Second)
}

func TestUpdateWallet(t *testing.T) {
	wallet1 := createRandomWallet(t)

	arg := UpdateWalletParams{
		ID:      wallet1.ID,
		Balance: util.RandomMoney(),
	}
	wallet2, err := testQueries.UpdateWallet(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	require.Equal(t, wallet1.ID, wallet2.ID)
	require.Equal(t, wallet1.Owner, wallet2.Owner)
	require.Equal(t, arg.Balance, wallet2.Balance)
	require.Equal(t, wallet1.Currency, wallet2.Currency)
	require.WithinDuration(t, wallet1.CreatedAt, wallet2.CreatedAt, time.Second)
}

func TestDeleteWallet(t *testing.T) {
	wallet1 := createRandomWallet(t)
	err := testQueries.DeleteWallet(context.Background(), wallet1.ID)
	require.NoError(t, err)

	wallet2, err := testQueries.GetWallet(context.Background(), wallet1.ID)
	require.Error(t, err)
	require.Equal(t, err, sql.ErrNoRows)
	require.Empty(t, wallet2)
}

func TestListWallets(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomWallet(t)
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
	}
}
