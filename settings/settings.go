package settings

import (
	"encoding/json"
	"os"
)

type AppSettings struct {
	Host  string
	Port  int
	Debug bool
}

type PaymentGateSettings struct {
	BitcoinInsightAPIURL       string `json:"bitcoin_insight_api_url"`
	EthereumAPIURL             string `json:"ethereum_api_url"`
	BTCCommissionWallet        string `json:"btc_commission_wallet"`
	XMRCommissionWallet        string `json:"xmr_commission_wallet"`
	PostgresConnectionString   string `json:"pg_connection_string"`
	TochkaEscrowPaymentAddress string `json:"tochka_escrow_payment_address"`
	MoneroRPCURL               string `json:"monero_rpc_url"`
	MoneroWalletRPCURL         string `json:"monero_wallet_rpc_url"`
	MoneroWalletRPCUser        string `json:"monero_wallet_rpc_user"`
	MoneroWalletRPCPassword    string `json:"monero_wallet_rpc_password"`
	EtherscanAPIKey            string `json:"etherscan_api_key"`
	ElectrumAPIURL             string `json:"electrum_api_url"`
	BitcoinRPCExplorerURL      string `json:"bre_url"`
	Host                       string `json:"host"`
	Port                       int    `json:"port"`
	CryptocompareToken         string `json:"cryptocompare_token"`
}

/*
	Utility Functions
*/

func LoadPaymentGateSettings() PaymentGateSettings {
	var settings PaymentGateSettings

	file, err := os.ReadFile("settings.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &settings)
	if err != nil {
		panic(err)
	}

	return settings
}

/*
	Globals
*/

var (
	PAYMENT_GATE_SETTINGS = LoadPaymentGateSettings()

	APPLICATION_SETTINGS = AppSettings{
		Host:  PAYMENT_GATE_SETTINGS.Host,
		Port:  PAYMENT_GATE_SETTINGS.Port,
		Debug: os.Getenv("PAYAKA_DEBUG") == "1",
	}
)
