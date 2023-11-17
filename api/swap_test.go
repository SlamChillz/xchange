package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slamchillz/xchange/token"
	mockdb "github.com/slamchillz/xchange/db/mock"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestCoinSwapRequest(t *testing.T) {
	customerId := utils.RandomNumber()
	swapReq := CoinSwapRequest{
		CoinName: "BTC",
		CoinAmountToSwap: 0.01,
		Network: "BTC",
		PhoneNumber: "08023222554",
		BankAccName: "Access Bank Nigeria",
		BankAccNumber: "0031961808",
		BankCode: "044",
	}
	arg := db.GetPendingNetworkTransactionParams{
		CustomerID: customerId,
		Network: swapReq.Network,
		TransactionStatus: "PENDING",
	}
	testCases := []struct {
		name string
		body gin.H
		stubs func (*mockdb.MockStore)
		setUpAuth func (t *testing.T, req *http.Request, authToken *token.JWT)
		response func (*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"coin_name": swapReq.CoinName,
				"coin_amount_to_swap": swapReq.CoinAmountToSwap,
				"network": swapReq.Network,
				"phone_number": swapReq.PhoneNumber,
				"bank_acc_name": swapReq.BankAccName,
				"bank_acc_number": swapReq.BankAccNumber,
				"bank_code": swapReq.BankCode,
			},
			setUpAuth: func (t *testing.T, req *http.Request, authToken *token.JWT) {
				addAccessTokenToRequestHeader(t, req, authToken, AUTHENTICATIONSCHEME, customerId, time.Minute)
			},
			stubs: func(storage *mockdb.MockStore) {
				storage.EXPECT().
					GetPendingNetworkTransaction(gomock.Any(), gomock.Eq(arg)).
					Return(int64(0), nil).
					Times(1)
				storage.EXPECT().
					GetBtcAddress(gomock.Any(), gomock.Eq(arg.CustomerID)).
					Return(sql.NullString{}, nil).
					Times(1)
				storage.EXPECT().
					CreateSwap(gomock.Any(), gomock.Any()).
					Return(db.Coinswap{}, nil).
					Times(1)
				storage.EXPECT().
					InsertNewBtcAddress(gomock.Any(), gomock.Any()).
					Return(db.Customerasset{}, nil).
					Times(1)
			},
			response: func (t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func (t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.stubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			reqBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/swap"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
			require.NoError(t, err)

			tc.setUpAuth(t, request, server.token)

			server.router.ServeHTTP(recorder, request)
			// body := recorder.Body.String()
			// fmt.Println(body)
			tc.response(t, recorder)
		})
	}
}
