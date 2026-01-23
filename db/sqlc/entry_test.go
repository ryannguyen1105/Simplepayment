package db

import (
	"context"
	"testing"
	"time"

	"github.com/ryannguyen1105/Simplepayment/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, walletID int64, paymentID int64) Entry {
	t.Helper()

	arg := CreateEntryParams{
		WalletID:  walletID,
		PaymentID: paymentID, // 2. Sử dụng ID vừa tạo
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotZero(t, entry.ID)
	require.Equal(t, walletID, entry.WalletID)
	require.Equal(t, paymentID, entry.PaymentID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.WithinDuration(t, time.Now(), entry.CreatedAt, time.Second)

	return entry
}

func TestCreateEntry(t *testing.T) {
	fromUser := createRandomUser(t)
	toUser := createRandomUser(t)

	wallet := createRandomWallet(t, fromUser.ID)
	payment := createRandomPayment(t, fromUser.ID, toUser.ID)

	createRandomEntry(t, wallet.ID, payment.ID)
}

func TestGetEntry(t *testing.T) {
	fromUser := createRandomUser(t)
	toUser := createRandomUser(t)

	wallet := createRandomWallet(t, fromUser.ID)
	payment := createRandomPayment(t, fromUser.ID, toUser.ID)

	entry1 := createRandomEntry(t, wallet.ID, payment.ID)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.WalletID, wallet.ID)
	require.Equal(t, entry1.PaymentID, payment.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)

	require.WithinDuration(t, time.Now(), entry1.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	fromUser := createRandomUser(t)
	toUser := createRandomUser(t)

	wallet := createRandomWallet(t, fromUser.ID)
	payment := createRandomPayment(t, fromUser.ID, toUser.ID)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, wallet.ID, payment.ID)
	}

	arg := ListEntriesParams{
		WalletID: wallet.ID,
		Limit:    5,
		Offset:   5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, wallet.ID, entry.WalletID)
		require.Equal(t, payment.ID, entry.PaymentID)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)
	}
}
