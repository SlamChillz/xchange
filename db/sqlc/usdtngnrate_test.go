package db

import (
	"fmt"
	"context"
	"testing"

	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUsdtNgnRate(t *testing.T) Usdtngnrate {
	swaprate := utils.RandomCoinswapRate()
	usdtngnrate, err := testQueries.RecentUsdtNgnRate(context.Background(), fmt.Sprintf("%.8f", swaprate))
	require.NoError(t, err)
	require.NotEmpty(t, usdtngnrate)
	require.Equal(t, fmt.Sprintf("%.8f", swaprate), usdtngnrate.UsdtNgnRate)
	require.NotZero(t, usdtngnrate.CreatedAt)
	require.NotZero(t, usdtngnrate.UpdatedAt)
	return usdtngnrate
}

func TestRecentUsdtNgnRate(t *testing.T) {
	createRandomUsdtNgnRate(t)
}
