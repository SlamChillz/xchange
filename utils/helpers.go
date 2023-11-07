package utils

import (
	"fmt"
	"net/http"
	"encoding/json"
)

func CurrentCoinPriceInUSD(coin string) (float64, error) {
	pricedata := map[string]map[string]float64 {}
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coin)
	resp, err := http.Get(url)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&pricedata); err != nil {
		return 0.0, err
	}
	if _, ok := pricedata[coin]; !ok {
		return 0.0, fmt.Errorf("coin not found")
	}
	if _, ok := pricedata[coin]["usd"]; !ok {
		return 0.0, fmt.Errorf("coin price in usd not found")
	}
	price := pricedata[coin]["usd"]
	return price, nil
}
