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
	config, err := utils.LoadConfig("../")
	require.NoError(t, err)
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
