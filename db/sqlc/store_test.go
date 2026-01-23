package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentTx(t *testing.T) {

	store := NewStore(testDB)

	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	fromWallet := createRandomWallet(t, user1.ID)
	toWallet := createRandomWallet(t, user2.ID)
	fmt.Println(">> before", fromWallet.Balance, toWallet.Balance)
	initialFromBalance := fromWallet.Balance
	initialToBalance := toWallet.Balance
	n := 5
	amount := int64(10)

	// run n concurrent payment transaction
	errs := make(chan error)
	results := make(chan PaymentTxResult)
	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.PaymentTx(ctx, PaymentTxParams{
				FromWalletID: fromWallet.ID,
				ToWalletID:   toWallet.ID,
				Amount:       amount,
			})
			errs <- err
			results <- result
		}()
	}
	// check result
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// check payment
		payment := result.Payment
		require.NotEmpty(t, payment)
		require.Equal(t, user1.ID, payment.FromUserID)
		require.Equal(t, user2.ID, payment.ToUserID)
		require.Equal(t, amount, payment.Amount)
		require.NotZero(t, payment.ID)
		require.NotZero(t, payment.CreatedAt)

		_, err = store.GetPayment(context.Background(), payment.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromWallet.ID, fromEntry.WalletID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toWallet.ID, toEntry.WalletID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check wallet
		fromWalletResult := result.FromWallet
		require.NotEmpty(t, fromWalletResult)
		require.Equal(t, user1.ID, fromWalletResult.UserID)

		toWalletResult := result.ToWallet
		require.NotEmpty(t, toWalletResult)
		require.Equal(t, user2.ID, toWalletResult.UserID)

		// check balance
		fmt.Println(">> before", fromWallet.Balance, toWallet.Balance)
		diff := initialFromBalance - fromWalletResult.Balance

		require.True(t, diff > 0)
		require.True(t, diff%amount == 0)

		k := int(diff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	// check the final updated balances
	updatedWallet1, err := testQueries.GetWalletByID(context.Background(), fromWallet.ID)
	require.NoError(t, err)

	updatedWallet2, err := testQueries.GetWalletByID(context.Background(), toWallet.ID)
	require.NoError(t, err)

	fmt.Println(">> before", updatedWallet1.Balance, updatedWallet2.Balance)
	require.Equal(t, initialFromBalance-int64(n)*amount, updatedWallet1.Balance)
	require.Equal(t, initialToBalance+int64(n)*amount, updatedWallet2.Balance)

}

func TestPaymentTxDeadlock(t *testing.T) {

	store := NewStore(testDB)

	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	fromWallet := createRandomWallet(t, user1.ID)
	toWallet := createRandomWallet(t, user2.ID)
	fmt.Println(">> before", fromWallet.Balance, toWallet.Balance)
	initialFromBalance := fromWallet.Balance
	initialToBalance := toWallet.Balance
	n := 10
	amount := int64(10)

	// run n concurrent payment transaction
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromWalletID := fromWallet.ID
		toWalletID := toWallet.ID

		if i%2 == 1 {
			fromWalletID = toWalletID
			toWalletID = fromWallet.ID
		}

		go func() {
			_, err := store.PaymentTx(context.Background(), PaymentTxParams{
				FromWalletID: fromWalletID,
				ToWalletID:   toWalletID,
				Amount:       amount,
			})
			errs <- err
		}()
	}
	// check result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}
	// check the final updated balances
	updatedWallet1, err := testQueries.GetWalletByID(context.Background(), fromWallet.ID)
	require.NoError(t, err)

	updatedWallet2, err := testQueries.GetWalletByID(context.Background(), toWallet.ID)
	require.NoError(t, err)

	fmt.Println(">> before", updatedWallet1.Balance, updatedWallet2.Balance)
	require.Equal(t, initialFromBalance, updatedWallet1.Balance)
	require.Equal(t, initialToBalance, updatedWallet2.Balance)

}
