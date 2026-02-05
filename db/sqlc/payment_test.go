package db

import (
	"context"
	"testing"
	"time"

	"github.com/ryannguyen1105/Simplepayment/util"
	"github.com/stretchr/testify/require"
)

func createRandomPayment(t *testing.T, wallet1, wallet2 Wallet) Payment {
	t.Helper()

	arg := CreatePaymentParams{
		FromWalletID: wallet1.ID,
		ToWalletID:   wallet2.ID,
		Amount:       util.RandomMoney(),
		Status:       util.RandomStatus(),
	}
	payment, err := testQueries.CreatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payment)

	require.Equal(t, arg.FromWalletID, payment.FromWalletID)
	require.Equal(t, arg.ToWalletID, payment.ToWalletID)
	require.Equal(t, arg.Amount, payment.Amount)

	require.NotZero(t, payment.ID)
	require.NotZero(t, payment.CreatedAt)

	return payment
}

func TestCreatePayment(t *testing.T) {
	wallet1 := createRandomWallet(t)
	wallet2 := createRandomWallet(t)
	createRandomPayment(t, wallet1, wallet2)
}

func TestGetPayment(t *testing.T) {
	wallet1 := createRandomWallet(t)
	wallet2 := createRandomWallet(t)
	payment1 := createRandomPayment(t, wallet1, wallet2)

	payment2, err := testQueries.GetPayment(context.Background(), payment1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, payment2)

	require.Equal(t, payment1.ID, payment2.ID)
	require.Equal(t, payment1.FromWalletID, payment2.FromWalletID)
	require.Equal(t, payment1.ToWalletID, payment2.ToWalletID)
	require.Equal(t, payment1.Amount, payment2.Amount)
	require.WithinDuration(t, payment1.CreatedAt, payment2.CreatedAt, time.Second)
}

func TestListPayments(t *testing.T) {
	wallet1 := createRandomWallet(t)
	wallet2 := createRandomWallet(t)

	for i := 0; i < 5; i++ {
		createRandomPayment(t, wallet1, wallet2)
		createRandomPayment(t, wallet2, wallet1)
	}
	arg := ListPaymentsParams{
		FromWalletID: wallet1.ID,
		ToWalletID:   wallet2.ID,
		Limit:        5,
		Offset:       0,
	}
	payments, err := testQueries.ListPayments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, payments, 5)

	for _, payment := range payments {
		require.NotEmpty(t, payment)
		require.Equal(t, arg.FromWalletID, payment.FromWalletID)
		require.NotZero(t, payment.ID)
		require.NotZero(t, payment.CreatedAt)
	}
}

func TestCancelPayment(t *testing.T) {
	wallet1 := createRandomWallet(t)
	wallet2 := createRandomWallet(t)

	payment := createRandomPayment(t, wallet1, wallet2)

	cancelPayment, err := testQueries.CancelPayment(context.Background(), payment.ID)
	require.NoError(t, err)
	require.NotEmpty(t, cancelPayment)

	require.Equal(t, payment.ID, cancelPayment.ID)

	require.WithinDuration(t, payment.CreatedAt, cancelPayment.CreatedAt, time.Second)
}
