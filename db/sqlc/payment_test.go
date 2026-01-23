package db

import (
	"context"
	"testing"
	"time"

	"github.com/ryannguyen1105/Simplepayment/util"
	"github.com/stretchr/testify/require"
)

func createRandomPayment(t *testing.T, fromUserID, toUserID int64) Payment {
	t.Helper()

	arg := CreatePaymentParams{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     util.RandomMoney(),
	}
	payment, err := testQueries.CreatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payment)

	require.Equal(t, arg.FromUserID, payment.FromUserID)
	require.Equal(t, arg.ToUserID, payment.ToUserID)
	require.Equal(t, arg.Amount, payment.Amount)

	require.NotZero(t, payment.ID)
	require.NotZero(t, payment.CreatedAt)

	return payment
}

func TestCreatePayment(t *testing.T) {
	fromUser := createRandomUser(t)
	toUser := createRandomUser(t)

	payment := createRandomPayment(t, fromUser.ID, toUser.ID)

	require.NotEmpty(t, payment)
	require.Equal(t, fromUser.ID, payment.FromUserID)
	require.Equal(t, toUser.ID, payment.ToUserID)
}

func TestGetPayment(t *testing.T) {
	fromUser := createRandomUser(t)
	toUser := createRandomUser(t)

	payment1 := createRandomPayment(t, fromUser.ID, toUser.ID)
	payment2, err := testQueries.GetPayment(context.Background(), payment1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, payment2)

	require.Equal(t, payment1.ID, payment2.ID)
	require.Equal(t, payment1.FromUserID, payment2.FromUserID)
	require.Equal(t, payment1.ToUserID, payment2.ToUserID)
	require.Equal(t, payment1.Amount, payment2.Amount)
	require.Equal(t, payment1.Status, payment2.Status)

	require.WithinDuration(t, payment1.CreatedAt, payment2.CreatedAt, time.Second)
}

func TestListPayments(t *testing.T) {
	fromUser := createRandomUser(t)
	toUser := createRandomUser(t)

	for i := 0; i < 10; i++ {
		createRandomPayment(t, fromUser.ID, toUser.ID)
	}

	arg := ListPaymentsParams{
		FromUserID: fromUser.ID,
		Limit:      5,
		Offset:     5,
	}

	payments, err := testQueries.ListPayments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, payments, 5)

	for _, payment := range payments {
		require.NotEmpty(t, payment)
		require.Equal(t, fromUser.ID, payment.FromUserID)
		require.NotZero(t, payment.ID)
		require.NotZero(t, payment.CreatedAt)
	}
}

func TestCancelPayment(t *testing.T) {
	fromUser := createRandomUser(t)
	toUser := createRandomUser(t)

	payment := createRandomPayment(t, fromUser.ID, toUser.ID)

	cancelPayment, err := testQueries.CancelPayment(context.Background(), payment.ID)
	require.NoError(t, err)
	require.NotEmpty(t, cancelPayment)

	require.Equal(t, payment.ID, cancelPayment.ID)

	require.WithinDuration(t, payment.CreatedAt, cancelPayment.CreatedAt, time.Second)
}
