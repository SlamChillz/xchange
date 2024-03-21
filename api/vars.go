package api

import (
	"errors"

	log "github.com/slamchillz/xchange/logger"
)

const BTCUSDT = 34403.000000

var (
	logger = log.GetLogger()

	ErrInvalidBankAccount = errors.New("unable to verify account number")

	BTCDATA = []map[string]string{
		{"coin_name": "btc", "coin_id": "bitcoin", "coin_symbol": "BTC"},
	}

	STABLECOINSLIST = []string{"USDT", "USDT_TRON", "USDT_BSC"}

	BITPOWRCOINTICKER = map[string]string{
		"BTC": "BTC",
		"USDT_TRON": "USDT_TRON",
		"USDT_BSC": "USDT_BSC",
		"USDT": "USDT",
	}

	CHAINS = map[string]string {
		"BTC": "BTC",
		"USDT_TRON": "TRON",
		"USDT_BSC": "BSC",
		"USDT": "ETH",
	}
	
	NETWORKS = map[string]string {
		"BTC": "BTC",
		"TRON": "USDT_TRON",
		"BSC": "USDT_BSC",
		"ETH": "USDT",
	}
)
