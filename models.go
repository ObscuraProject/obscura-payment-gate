package main

import (
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/bitcoin"
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/ethereum"
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/monero"
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/polkadot"
	"github.com/ObscuraProject/obscura-payment-gate/db"
)

func SyncDatabase() {
	db.Database.AutoMigrate(
		&bitcoin.BitcoinWallet{},
		&bitcoin.BitcoinWalletBalance{},

		&ethereum.EthereumWallet{},
		&ethereum.EthereumWalletBalance{},

		&monero.MoneroWallet{},
		&monero.MoneroWalletBalance{},

		&polkadot.PolkadotWallet{},
		&polkadot.PolkadotWalletBalance{},
	)

	ethereum.SetupViews()
	bitcoin.SetupViews()
	monero.SetupViews()
	polkadot.SetupViews()
}
