package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		BITPOWR_ACCOUNT_ID: os.Getenv("BITPOWR_ACCOUNT_ID"),
		BITPOWR_API_KEY: os.Getenv("BITPOWR_API_KEY"),
		SHUTTER_PUBLIC_KEY: os.Getenv("SHUTTER_PUBLIC_KEY"),
		VALIDATE_BANK_URL: os.Getenv("VALIDATE_BANK_URL"),
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
