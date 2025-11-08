package etherscan

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ObscuraProject/obscura-payment-gate/apis"
	"github.com/ObscuraProject/obscura-payment-gate/settings"
)

type EtherscanPushTransactionResponse struct {
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result *string `json:"result"`
}

func PushEtherscanTransacation(rawtx string) (string, error) {
	var (
		etherscanURL = fmt.Sprintf(
			"https://api.etherscan.io/api?module=proxy&action=eth_sendRawTransaction&hex=%s&apikey=%s",
			rawtx,
			settings.PAYMENT_GATE_SETTINGS.EtherscanAPIKey,
		)
		result = EtherscanPushTransactionResponse{}
	)

	response, err := apis.TorGET(etherscanURL)
	if err != nil {
		return response, errors.New("can't push to etherscan")
	}

	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		return response, errors.New("can't push to etherscan")
	}

	if result.Error != nil {
		return response, errors.New("can't push to etherscan")
	}

	if result.Result != nil {
		return *result.Result, nil
	}

	return "", errors.New("Unexpected result")
}
