package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentTx(t *testing.T) {
	store := NewStore(testDB)

	wallet1 := createRandomWallet(t)
	wallet2 := createRandomWallet(t)
	fmt.Println(">> before:", wallet1.Balance, wallet2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan PaymentTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.PaymentTx(ctx, PaymentTxParams{
				FromWalletID: wallet1.ID,
				ToWalletID:   wallet2.ID,
				Amount:       amount,
			})
			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		payment := result.Payment
		require.NotEmpty(t, payment)
		require.Equal(t, wallet1.ID, payment.FromWalletID)
		require.Equal(t, wallet2.ID, payment.ToWalletID)
		require.Equal(t, amount, payment.Amount)
		require.NotZero(t, payment.ID)
		require.NotZero(t, payment.CreatedAt)

		_, err = store.GetWallet(context.Background(), wallet1.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, wallet1.ID, fromEntry.WalletID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, wallet2.ID, toEntry.WalletID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)

		fromWallet := result.FromWallet
		require.NotEmpty(t, fromWallet)
		require.Equal(t, wallet1.ID, fromWallet.ID)

		toWallet := result.ToWallet
		require.NotEmpty(t, toWallet)
		require.Equal(t, wallet2.ID, toWallet.ID)

		fmt.Println(">> tx:", fromWallet.Balance, toWallet.Balance)
		diff1 := wallet1.Balance - fromWallet.Balance
		diff2 := toWallet.Balance - wallet2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updateWallet1, err := testQueries.GetWallet(context.Background(), wallet1.ID)
	require.NoError(t, err)

	updateWallet2, err := testQueries.GetWallet(context.Background(), wallet2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateWallet1.Balance, updateWallet2.Balance)
	require.Equal(t, wallet1.Balance-int64(n)*amount, updateWallet1.Balance)
	require.Equal(t, wallet2.Balance+int64(n)*amount, updateWallet2.Balance)

}

func TestPaymentTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	wallet1 := createRandomWallet(t)
	wallet2 := createRandomWallet(t)
	fmt.Println(">> before:", wallet1.Balance, wallet2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromWalletID := wallet1.ID
		toWalletID := wallet2.ID

		if i%2 == 1 {
			fromWalletID = wallet2.ID
			toWalletID = wallet1.ID
		}
		go func() {
			ctx := context.Background()
			_, err := store.PaymentTx(ctx, PaymentTxParams{
				FromWalletID: fromWalletID,
				ToWalletID:   toWalletID,
				Amount:       amount,
			})
			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	updateWallet1, err := testQueries.GetWallet(context.Background(), wallet1.ID)
	require.NoError(t, err)

	updateWallet2, err := testQueries.GetWallet(context.Background(), wallet2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateWallet1.Balance, updateWallet2.Balance)
	require.Equal(t, wallet1.Balance, updateWallet1.Balance)
	require.Equal(t, wallet2.Balance, updateWallet2.Balance)

}
