package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"database/sql"
)

const numbers = "0123456789"
const laplhas = "abcdefghijklmnopqrstuvwxyz"
const calphas = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const hexvalues = "0123456789aAbBcCdDeEfF"

var charset = laplhas

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomBtcAddres generates a random bitcoin address
func RandomBtcAddress() sql.NullString {
	charset = charset + numbers
	var address strings.Builder
	address.WriteString("bc")
	address.WriteString(RandomString(40))
	return sql.NullString{String: address.String(), Valid: true}
}

// RandomBepAddress generates random bep20 address
func RandomBepAddress() sql.NullString {
	return sql.NullString{String: RandomEthAddress().String, Valid: true}
}

// RandomEthAddress generates a random ethereum address
func RandomEthAddress() sql.NullString {
	var address strings.Builder
	charset = hexvalues
	address.WriteString("0x")
	address.WriteString(RandomString(40))
	return sql.NullString{String: address.String(), Valid: true}
}

// RandomTrcAddress generates a random tron address
func RandomTrcAddress() sql.NullString {
	charset = charset + calphas + numbers
	var address strings.Builder
	address.WriteByte('T')
	address.WriteByte(calphas[rand.Intn(len(calphas))])
	address.WriteString(RandomString(32))
	return sql.NullString{String: address.String(), Valid: true}
}

func RandomUsdtAddress() sql.NullString {
	return sql.NullString{String: RandomEthAddress().String, Valid: true}
}

func RandomUsdtTronAddress() sql.NullString {
	return sql.NullString{String: RandomTrcAddress().String, Valid: true}
}

// RandomCoinAddress generates a random coin address
func RandomCoinAddress(network string) sql.NullString {
	switch network {
	case "BTC":
		return RandomBtcAddress()
	case "ETH":
		return RandomEthAddress()
	case "TRC20":
		return RandomTrcAddress()
	case "BEP20":
		return RandomBepAddress()
	default:
		return sql.NullString{String: "", Valid: false}
	}
}

// RandomString generate a random string
func RandomString(n int) string {
	var stringBuilder strings.Builder
	l := len(charset)
	for i := 0; i < n; i++ {
		v := charset[rand.Intn(l)]
		stringBuilder.WriteByte(v)
	}
	return stringBuilder.String()
}

func RandomNumber() int32 {
	return int32(RandomInt(0, 100) + RandomInt(100, 200))
}

// RandomName genrates a random name
func RandomName() string {
	charset = laplhas
	return RandomString(20)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(8))
}

func RandomCoinName() string {
	coins := []string{"BTC", "ETHUSDT", "TRNUSDT", "BEPUSDT", "BUSD"}
	return coins[rand.Intn(len(coins))]
}

func RandomPhoneNumber() string {
	charset = numbers
	var phoneNumber strings.Builder
	phoneNumber.WriteString("080")
	phoneNumber.WriteString(RandomString(8))
	return phoneNumber.String()
}

func RandomCoinNetwork() string {
	network := []string{"BTC", "ETH", "TRC20", "BEP20"}
	return network[rand.Intn(len(network))]
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64() * (max - min + 1)
}

func RandomCoinSwapAmount() float64 {
	return RandomFloat(20.00, 1000.00)
}

func RandomTranstatus() string {
	status := []string{"PENDING", "CANCELED", "FAILED", "SUCCESS"}
	return status[rand.Intn(len(status))]
}

func RandomCoinswapRate() float64 {
	return RandomFloat(900.00, 1000.00)
}

func RandomPayoutStatus() string {
	status := []string{"PAID", "UNPAID"}
	return status[rand.Intn(len(status))]
}

func RandomBankName() string {
	var bankname strings.Builder
	bankname.WriteString(RandomString(int(RandomInt(4, 10))))
	return bankname.String()
}

func RandomBankAccount() string {
	charset = numbers
	return RandomString(10)
}

func RandomBankCode() string {
	charset = numbers
	return RandomString(3)
}
