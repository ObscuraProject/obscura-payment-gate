package cryptocompare

import (
	"encoding/json"
	"strings"

	"github.com/ObscuraProject/obscura-payment-gate/apis"
	"github.com/ObscuraProject/obscura-payment-gate/settings"
)

const CRYPTOCOMPARE_RATES_URL = "https://min-api.cryptocompare.com/data/price?api_key={api_key}&fsym={base_currency}&tsyms=USD,RUB,EUR,GBP,AUD,ETH,BTC,USDT"

type CryptocompateCurrencyResponse struct {
	USD  float64 `json:"USD"`
	RUB  float64 `json:"RUB"`
	EUR  float64 `json:"EUR"`
	AUD  float64 `json:"AUD"`
	GBP  float64 `json:"GBP"`
	ETH  float64 `json:"ETH"`
	BTC  float64 `json:"BTC"`
	USDT float64 `json:"USDT"`
	DOT  float64 `json:"DOT"`
}

func GetCryptocompareCurrencyResponse(baseCurrency string) (CryptocompateCurrencyResponse, error) {
	currencyResponse := CryptocompateCurrencyResponse{}
	url := strings.Replace(CRYPTOCOMPARE_RATES_URL, "{base_currency}", baseCurrency, -1)
	url = strings.Replace(url, "{api_key}", settings.PAYMENT_GATE_SETTINGS.CryptocompareToken, -1)
	response, err := apis.DirectGET(url)
	if err != nil {
		return currencyResponse, err
	}
	currencyResponse.BTC = 1.0
	err = json.Unmarshal([]byte(response), &currencyResponse)
	return currencyResponse, err
}
