package bre

import (
	"encoding/json"

	"github.com/ObscuraProject/obscura-payment-gate/apis"
	"github.com/ObscuraProject/obscura-payment-gate/settings"
)

/*
Globals
*/
var (
	API = BREAPI{APIEndpoint: settings.PAYMENT_GATE_SETTINGS.BitcoinRPCExplorerURL}
)

/*
	API Definition
*/

type BREAPI struct {
	APIEndpoint string
}

/*
	Responses
*/

type TxHistoryResponse struct {
	TxIds      []string `json:"txids"`
	TxCount    int      `json:"txCount"`
	BalanceSat int      `json:"balanceSat"`
}

type AddressResponse struct {
	TxHistory TxHistoryResponse `json:"txHistory"`
}

type NextBlockResponse struct {
	Smart  int `json:"smart"`
	Min    int `json:"min" `
	Max    int `json:"max"`
	Median int `json:"median"`
}

type MempoolFeesResponse struct {
	NextBlock NextBlockResponse `json:"nextBlock"`
}

/*
	Methods
*/

func (api BREAPI) GetWalletBalance(publicKey string) (TxHistoryResponse, error) {
	addr := api.APIEndpoint + "api/address/" + publicKey
	response, err := apis.DirectGET(addr)
	if err != nil {
		return TxHistoryResponse{}, err
	}
	wbr := AddressResponse{}
	err = json.Unmarshal([]byte(response), &wbr)
	if err != nil {
		return TxHistoryResponse{}, err
	}
	return wbr.TxHistory, nil
}

func (api BREAPI) MempoolFees() (NextBlockResponse, error) {
	addr := api.APIEndpoint + "api/mempool/fees/"
	response, err := apis.DirectGET(addr)
	if err != nil {
		return NextBlockResponse{}, err
	}
	wbr := MempoolFeesResponse{}
	err = json.Unmarshal([]byte(response), &wbr)
	if err != nil {
		return NextBlockResponse{}, err
	}
	return wbr.NextBlock, nil
}
