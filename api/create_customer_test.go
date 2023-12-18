package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	mockdb "github.com/slamchillz/xchange/db/mock"
	db "github.com/slamchillz/xchange/db/sqlc"
	mockredisdb "github.com/slamchillz/xchange/redisdb/mock"
	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)


type CreateCustomerParamsMatcher struct {
	arg db.CreateCustomerParams
	password string
}

type CreateOTPMatcher struct {
	otp *string
}

func (c *CreateOTPMatcher) Matches(x interface{}) bool {
	arg, ok := x.(string)
	if !ok {
		return false
	}
	value, err := strconv.Atoi(arg)
	if err != nil || value <= 0 {
		return false
	}
	*c.otp = arg
	return true
}

func (c *CreateOTPMatcher) String() string {
	return fmt.Sprintf("matches create otp %v", *c.otp)
}

func EqCreateOTP(otp *string) gomock.Matcher {
	return &CreateOTPMatcher{otp}
}

func (c *CreateCustomerParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateCustomerParams)
	if !ok {
		return false
	}
	err := utils.CheckPassword(arg.Password, c.password)
	if err != nil {
		return false
	}
	c.arg.Password = arg.Password
	return reflect.DeepEqual(c.arg, arg)
}

func (c *CreateCustomerParamsMatcher) String() string {
	return fmt.Sprintf("matches create customers params %v and password %v", c.arg, c.password)
}

func EqCreateCustomerParams(arg db.CreateCustomerParams, password string) gomock.Matcher {
	return &CreateCustomerParamsMatcher{arg, password}
}

func TestCreateCustomerRequest(t *testing.T) {
	password := "password"
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	customer := db.Customer {
		ID: utils.RandomNumber(),
		FirstName: "John",
		LastName: "Benjamin",
		Email: "slamchillz@gmail.com",
		Phone: "08054444667",
		Password: hashedPassword,
	}
	testCases := []struct {
		name string
		body gin.H
		stubs func (
			*mockdb.MockStore, 
			*mockredisdb.MockRedisClient,
		)
		response func (*testing.T, *httptest.ResponseRecorder)
	}{
		{name: "OK",
		 body: gin.H{
			 "first_name": customer.FirstName,
			 "last_name": customer.LastName,
			 "email": customer.Email,
			 "phone_number": customer.Phone,
			 "password": password,
			 "confirm_password": password,
		 },
		 stubs: func(
			storage *mockdb.MockStore,
			rdbClient *mockredisdb.MockRedisClient,
		) {
			var otp string
			arg := db.CreateCustomerParams{
				FirstName: customer.FirstName,
				LastName: customer.LastName,
				Email: customer.Email,
				Phone: "+234" + customer.Phone[1:],
			}
			storage.EXPECT().
				CreateCustomer(gomock.Any(), EqCreateCustomerParams(arg, password)).
				Return(customer, nil).
				Times(1)
			rdbClient.EXPECT().
				Set(gomock.Any(), EqCreateOTP(&otp), gomock.Eq(customer.Email), gomock.Eq(5 * time.Minute)).
				Return(otp, nil).
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
			tc.stubs(store, redisClient)

			server := newTestServer(t, store, redisClient)
			recorder := httptest.NewRecorder()

			reqBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/user/signup"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// body := recorder.Body.String()
			// fmt.Println(body)
			tc.response(t, recorder)
		})
	}
}
