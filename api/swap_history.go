package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/token"
)

func parseDate(date string) (time.Time, error) {
	return time.Parse(time.RFC3339, date)
}

func (server *Server) ListCoinSwapHistory(ctx *gin.Context) {
	var err error
	var page int
	var pageSize int
	var status []string
	var asset []string
	var network []string
	var startDate time.Time
	var endDate time.Time
	payload, ok := ctx.Get(AUTHENTICATIONCONTEXTKEY)
	if !ok{
		logger.Error().Msg("error getting customer payload from authentication context key")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	customerPayload, ok := payload.(*token.Payload)
	if !ok {
		logger.Error().Interface("payload", payload).Msg("error casting customer payload to token payload")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	customerId := customerPayload.CustomerID
	if ps, err := strconv.Atoi(ctx.DefaultQuery("pagesize", "20")); err == nil {
		pageSize = ps
	} else {
		pageSize = 15
	}
	if p, err := strconv.Atoi(ctx.DefaultQuery("page", "1")); err == nil {
		page = p
	} else {
		page = 1
	}
	page = (page - 1) * pageSize
	status = ctx.QueryArray("status")
	asset = ctx.QueryArray("asset")
	network = ctx.QueryArray("network")
	startDate, err = parseDate(ctx.DefaultQuery("start_date", "1970-01-01T00:00:00Z"))
	if err != nil {
		logger.Error().Err(err).Msg("error parsing start date")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date in query params"})
		return
	}
	endDate, err = parseDate(ctx.DefaultQuery("end_date", time.Now().Format(time.RFC3339)))
	if err != nil {
		logger.Error().Err(err).Msg("error parsing end date")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date in query params"})
		return
	}
	if len(status) == 0 {
		status = []string{"PENDING", "CANCELED", "FAILED", "SUCCESS"}
	}
	if len(asset) == 0 {
		asset = []string{"BTC", "Bitcoin", "USDT", "USDT_TRON", "USDT_BSC"}
	}
	if len(network) == 0 {
		network = []string{"BTC", "ETH", "TRC20", "BEP20"}
	}
	params := db.ListAllCoinSwapTransactionsParams{
		// CustomerID: customerId,
		TransactionStatus: status,
		CoinName: asset,
		Network: network,
		Offset: int32(page),
		Limit: int32(pageSize),
		StartDate: startDate,
		EndDate: endDate,
	}
	transactions, err := server.storage.ListAllCoinSwapTransactions(ctx, params)
	if err != nil {
		logger.Error().Err(err).Interface("params", params).Msg("error getting coin swap transactions history")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	logger.Info().Int32("customer_id", customerId).Int("length", (len(transactions))).Msg("coin swap transactions history")
	ctx.JSON(http.StatusOK, transactions)
}
