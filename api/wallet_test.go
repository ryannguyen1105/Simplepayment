package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/ryannguyen1105/Simplepayment/db/mock"
	db "github.com/ryannguyen1105/Simplepayment/db/sqlc"
	"github.com/ryannguyen1105/Simplepayment/util"
	"github.com/stretchr/testify/require"
)

func TestGetWalletApi(t *testing.T) {
	wallet := randomWallet()

	testCases := []struct {
		name          string
		walletID      int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			walletID: wallet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Eq(wallet.ID)).
					Times(1).
					Return(wallet, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchAccount(t, recoder.Body, wallet)
			},
		},
		{
			name:     "NotFound",
			walletID: wallet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Eq(wallet.ID)).
					Times(1).
					Return(db.Wallet{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},
		{
			name:     "InternalError",
			walletID: wallet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Eq(wallet.ID)).
					Times(1).
					Return(db.Wallet{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name:     "InvalidID",
			walletID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			//store.EXPECT().
			//	GetWallet(gomock.Any(), gomock.Eq(wallet.ID)).
			//	Times(1).
			//	Return(wallet, nil)

			server := NewServer(store)
			recoder := httptest.NewRecorder()

			url := fmt.Sprintf("/wallets/%d", tc.walletID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			//require.Equal(t, http.StatusOK, recoder.Code)
			//requireBodyMatchAccount(t, recoder.Body, wallet)
			tc.checkResponse(t, recoder)
		})

	}
}

func TestCreateWalletApi(t *testing.T) {
	wallet := randomWallet()

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"owner":    wallet.Owner,
				"balance":  wallet.Balance,
				"currency": wallet.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateWallet(gomock.Any(), gomock.Any()).
					Times(1).
					Return(wallet, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchAccount(t, recoder.Body, wallet)
			},
		},
		{
			name: "InvalidInput",
			body: map[string]interface{}{
				"owner": "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InternalError",
			body: map[string]interface{}{
				"owner":    wallet.Owner,
				"balance":  wallet.Balance,
				"currency": wallet.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateWallet(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Wallet{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recoder := httptest.NewRecorder()

			url := "/wallets"
			data, _ := json.Marshal(tc.body)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.checkResponse(t, recoder)
		})
	}
}

func randomWallet() db.Wallet {
	return db.Wallet{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: "USD",
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, wallet db.Wallet) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotWallet db.Wallet
	err = json.Unmarshal(data, &gotWallet)
	require.NoError(t, err)
	require.Equal(t, wallet, gotWallet)
}
