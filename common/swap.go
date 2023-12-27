package common

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"net/http"

	"github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/utils"
)


func GenerateNewAddress(config utils.Config, customerId int32, label, asset, accountId string) (map[string]interface{}, error) {
	postData := map[string]string {
		"label": label,
		"asset": asset,
		"accountId": accountId,
	}
	reqData := new(bytes.Buffer)
	err := json.NewEncoder(reqData).Encode(postData)
	if err != nil {
		return nil, err
	}
	url, err := url.Parse("https://developers.bitpowr.com/api/v1/addresses")
	if err != nil {
		return nil, err
	}
	req := http.Request{
		Method: http.MethodPost,
		URL: url,
		Header: map[string][]string {
			"Content-Type": {"application/json"},
			"Authorization": {fmt.Sprintf("Bearer %s", config.BITPOWR_API_KEY)},
		},
		Body: io.NopCloser(reqData),
	}
	response, err := http.DefaultClient.Do(&req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to generate new %s address. Bitpowr StatusCoode: %v", asset, response.StatusCode)
	}
	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	data = data["data"].(map[string]interface{})
	return data, nil
}


func GetSwapAddress(config utils.Config, storage db.Store, coinName, coinNetwork string, customerId int32) (address string, err error) {
	var nullAddress sql.NullString
	var newAaddressData map[string]interface{}
	switch coinNetwork {
	case "BTC":
		nullAddress, err = storage.GetBtcAddress(context.Background(), customerId)
	case "BSC":
		nullAddress, err = storage.GetUsdtBscAddress(context.Background(), customerId)
	case "TRON":
		nullAddress, err = storage.GetUsdtTronAddress(context.Background(), customerId)
	case "ETH":
		nullAddress, err = storage.GetUsdtAddress(context.Background(), customerId)
	}
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return
		}
	}
	address = nullAddress.String
	if address == "" {
		label := coinName + "_" + coinNetwork
		asset := coinName
		accountId := config.BITPOWR_ACCOUNT_ID
		newAaddressData, err = GenerateNewAddress(config, customerId, label, asset, accountId)
		if err != nil {
			return
		}
	} else {
		return address, nil
	}
	address = newAaddressData["address"].(string)
	addressUid := newAaddressData["uid"].(string)
	chain := newAaddressData["chain"].(string)
	switch coinNetwork {
	case "BTC":
		_, err = storage.InsertNewBtcAddress(context.Background(), db.InsertNewBtcAddressParams{
			BtcAddress: sql.NullString{String: address, Valid: true},
			BtcAddressUid: sql.NullString{String: addressUid, Valid: true},
			BtcNetwork: sql.NullString{String: chain, Valid: true},
			CustomerID: customerId,
		})
	case "BSC":
		_, err = storage.InsertNewUsdtBscAddress(context.Background(), db.InsertNewUsdtBscAddressParams{
			UsdtBscAddress: sql.NullString{String: address, Valid: true},
			UsdtBscAddressUid: sql.NullString{String: addressUid, Valid: true},
			UsdtBscNetwork: sql.NullString{String: chain, Valid: true},
			CustomerID: customerId,
		})
	case "TRON":
		_, err = storage.InsertNewUsdtTronAddress(context.Background(), db.InsertNewUsdtTronAddressParams{
			UsdtTronAddress: sql.NullString{String: address, Valid: true},
			UsdtTronAddressUid: sql.NullString{String: addressUid, Valid: true},
			UsdtTronNetwork: sql.NullString{String: chain, Valid: true},
			CustomerID: customerId,
		})
	case "ETH":
		_, err = storage.InsertNewUsdtAddress(context.Background(), db.InsertNewUsdtAddressParams{
			UsdtAddress: sql.NullString{String: address, Valid: true},
			UsdtAddressUid: sql.NullString{String: addressUid, Valid: true},
			UsdtNetwork: sql.NullString{String: chain, Valid: true},
			CustomerID: customerId,
		})
	}
	if err != nil {
		return
	}
	return address, nil
}
