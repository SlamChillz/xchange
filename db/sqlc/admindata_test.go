package db

import (
	"context"
	"testing"

	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAdmindata(t *testing.T) Admindatum {
	arg := CreateAdmindataParams {
		BitpowrAccountID: utils.RandomString(10),
		BtcAddress: utils.RandomBtcAddress(),
		UsdtAddress: utils.RandomUsdtAddress(),
		UsdtTronAddress: utils.RandomUsdtTronAddress(),
	}

	admindata, err := testQueries.CreateAdmindata(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, admindata)
	require.Equal(t, arg.BitpowrAccountID, admindata.BitpowrAccountID)
	require.Equal(t, arg.BtcAddress, admindata.BtcAddress)
	require.Equal(t, arg.UsdtAddress, admindata.UsdtAddress)
	require.Equal(t, arg.UsdtTronAddress, admindata.UsdtTronAddress)
	require.NotZero(t, admindata.CreatedAt)
	require.NotZero(t, admindata.UpdatedAt)
	return admindata
}

func TestCreateAdmindata(t *testing.T) {
	createRandomAdmindata(t)
}
