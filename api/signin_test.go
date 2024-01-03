package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mockdb "github.com/slamchillz/xchange/db/mock"
	db "github.com/slamchillz/xchange/db/sqlc"
	mockredisdb "github.com/slamchillz/xchange/redisdb/mock"
	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestLoginCustomerRequest(t *testing.T) {
	password := "password"
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	customer := db.Customer {
		FirstName: "John",
		LastName: "Benjamin",
		Email: "slamchillz@gmail.com",
		Phone: sql.NullString{
			String: "+2347030000000",
			Valid: true,
		},
		Password: sql.NullString{
			String: hashedPassword,
			Valid: true,
		},
	}
	testCases := []struct {
		name string
		body gin.H
		stubs func (*mockdb.MockStore)
		response func (*testing.T, *httptest.ResponseRecorder)
	}{
		{name: "OK",
		 body: gin.H{
			 "email": customer.Email,
			 "password": password,
		 },
		 stubs: func(storage *mockdb.MockStore) {
			storage.EXPECT().
				GetCustomerByEmail(gomock.Any(), gomock.Eq(customer.Email)).
				Return(customer, nil).
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
			redisClient := mockredisdb.NewMockRedisClient(ctrl)
			tc.stubs(store)

			server := newTestServer(t, store, redisClient)
			recorder := httptest.NewRecorder()

			reqBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/user/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// body := recorder.Body.String()
			// fmt.Println(body)
			tc.response(t, recorder)
		})
	}
}
