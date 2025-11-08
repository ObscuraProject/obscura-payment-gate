package monero

import (
	"context"
	"fmt"

	"github.com/cirocosta/go-monero/pkg/http"
	"github.com/cirocosta/go-monero/pkg/rpc"
	"github.com/cirocosta/go-monero/pkg/rpc/wallet"
	"github.com/ObscuraProject/obscura-payment-gate/settings"
)

type MoneroAPI struct {
	daemonClient *rpc.Client
	walletClient *rpc.Client
}

var API = NewAPI(
	settings.PAYMENT_GATE_SETTINGS.MoneroRPCURL,
	settings.PAYMENT_GATE_SETTINGS.MoneroWalletRPCURL,
)

func NewAPI(daemonEndpoint, walletEndpoint string) MoneroAPI {
	daemonClient, err := rpc.NewClient(daemonEndpoint)
	if err != nil {
		panic(fmt.Errorf("new client for '%s': %w", daemonEndpoint, err))
	}

	walletHttpClientConfig := http.ClientConfig{
		TLSSkipVerify: true,
		Verbose:       false,
		Username:      settings.PAYMENT_GATE_SETTINGS.MoneroWalletRPCUser,
		Password:      settings.PAYMENT_GATE_SETTINGS.MoneroWalletRPCPassword,
	}

	walletRpcHttpClient, err := http.NewClient(walletHttpClientConfig)
	if err != nil {
		panic(fmt.Errorf("new client for '%s': %w", daemonEndpoint, err))
	}

	walletClient, err := rpc.NewClient(walletEndpoint, rpc.WithHTTPClient(walletRpcHttpClient))
	if err != nil {
		panic(fmt.Errorf("new client for '%s': %w", daemonEndpoint, err))
	}
	return MoneroAPI{
		daemonClient: daemonClient,
		walletClient: walletClient,
	}
}

func (api MoneroAPI) CreateWallet(label string, index uint) (*string, *uint, error) {
	ctx := context.TODO()
	result, err := wallet.
		NewClient(api.walletClient).
		CreateAddress(ctx, 0, 1, label)
	if err != nil {
		return nil, nil, err
	}

	return &result.Address, &result.AddressIndex, nil
}

func (api MoneroAPI) Balance(addressIndex uint) (*int64, *uint64, error) {
	ctx := context.TODO()
	result, err := wallet.NewClient(api.walletClient).GetBalance(
		ctx,
		wallet.GetBalanceRequestParameters{
			AccountIndex:   0,
			AddressIndices: []uint{addressIndex},
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return &result.PerSubaddress[0].UnlockedBalance, &result.PerSubaddress[0].Balance, nil
}

type AmountRecipient struct {
	Address string `json:"address"`
	Amount  uint   `json:"amount"`
}

type MoneroTransferResult struct {
	Amount        int64  `json:"amount"`
	Fee           int64  `json:"fee"`
	MultisigTxset string `json:"multisig_txset"`
	TxBlob        string `json:"tx_blob"`
	TxHash        string `json:"tx_hash"`
	TxMetadata    string `json:"tx_metadata"`
	TxKey         string `json:"tx_key"`
	UnsignedTxset string `json:"unsigned_txset"`
}

func (api MoneroAPI) Transfer(
	address string,
	subAccountIndex uint,
	amount uint,
) (*MoneroTransferResult, error) {
	client := wallet.NewClient(api.walletClient)

	// Create a transaction parameters JSON object
	params := struct {
		Destinations   []AmountRecipient `json:"destinations"`
		AccountIndex   uint              `json:"account_index"`
		SubaddrIndices []uint            `json:"subaddr_indices"`
		Priority       int64             `json:"priority"`
		RingSize       uint              `json:"ring_size"`
		GetTxKey       bool              `json:"get_tx_key"`
	}{
		Destinations: []AmountRecipient{
			{Amount: amount, Address: address},
		},
		AccountIndex:   0,
		SubaddrIndices: []uint{subAccountIndex},
		Priority:       0,
		RingSize:       7,
		GetTxKey:       true,
	}

	type Result struct {
		Id      int                  `json:"id"`
		JsonRpc string               `json:"jsonrpc"`
		Result  MoneroTransferResult `json:"result"`
	}
	result := Result{}

	ctx := context.TODO()
	err := client.Requester.JSONRPC(ctx, "transfer", &params, &result)

	return &result.Result, err
}

func (api MoneroAPI) SweepAll(
	address string,
	subAccountIndex uint,
) error {
	client := wallet.NewClient(api.walletClient)

	// Create a transaction parameters JSON object
	params := struct {
		Address        string `json:"address"`
		AccountIndex   uint   `json:"account_index"`
		SubaddrIndices []uint `json:"subaddr_indices"`
	}{
		Address:        address,
		AccountIndex:   0,
		SubaddrIndices: []uint{subAccountIndex},
	}

	type Result struct {
		Id      int    `json:"id"`
		JsonRpc string `json:"jsonrpc"`
	}
	result := Result{}

	ctx := context.TODO()
	err := client.Requester.JSONRPC(ctx, "sweep_all", &params, &result)

	return err
}
