package api

import (
	"log"
	"os"
)

var BTCDATA = []map[string]string{
	{"coin_name": "btc", "coin_id": "bitcoin", "coin_symbol": "BTC"},
}

var STABLECOINSLIST = []string{"USDT", "USDT_TRON", "USDT_BSC"}

var BITPOWRCOINTICKER = map[string]string{
	"BTC": "BTC",
	"USDT_TRON": "USDT_TRON",
	"USDT_BSC": "USDT_BSC",
	"USDT": "USDT",
}

var CHAINS = map[string]string {
	"BTC": "BTC",
	"USDT_TRON": "TRON",
	"USDT_BSC": "BSC",
	"USDT": "ETH",
}

var logger = log.New(os.Stdout, "api: ", log.Llongfile)
