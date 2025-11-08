package main

import (
	"os"

	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/bitcoin"
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/ethereum"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 1 && argsWithoutProg[0] == "sync" {
		SyncDatabase()
		// -------
		// Bitcoin
		// -------
	} else if len(argsWithoutProg) == 2 && argsWithoutProg[0] == "bitcoin" && argsWithoutProg[1] == "update" {
		bitcoin.UpdateBitcoinWalletBalances()
	} else if len(argsWithoutProg) == 4 && argsWithoutProg[0] == "bitcoin" && argsWithoutProg[1] == "sweep" {
		bitcoin.SweepBitcoinWallet(argsWithoutProg[2], argsWithoutProg[3])
	} else if len(argsWithoutProg) == 3 && argsWithoutProg[0] == "bitcoin" && argsWithoutProg[1] == "sweep-commission-premium" {
		bitcoin.SweepCommissionAndPremiumBitcoinWallets(argsWithoutProg[2])
	} else if len(argsWithoutProg) == 5 && argsWithoutProg[0] == "bitcoin" && argsWithoutProg[1] == "consolidate" {
		inputAddress := argsWithoutProg[2]
		inputPk := argsWithoutProg[3]
		outputAddress := argsWithoutProg[4]
		bitcoin.SanitizeBitcoinAddress(inputAddress, inputPk, outputAddress)
		// --------
		// Ethereum
		// --------
	} else if len(argsWithoutProg) == 2 && argsWithoutProg[0] == "ethereum" && argsWithoutProg[1] == "update" {
		ethereum.UpdateEthereumWalletBalances()
	} else if len(argsWithoutProg) == 4 && argsWithoutProg[0] == "ethereum" && argsWithoutProg[1] == "sweep" {
		ethereum.SweepEthereumWallet(argsWithoutProg[2], argsWithoutProg[3])
		// } else if len(argsWithoutProg) == 3 && argsWithoutProg[0] == "bitcoin" && argsWithoutProg[1] == "sweep-commission-premium" {
		// SweepCommissionAndPremiumBitcoinWallets(argsWithoutProg[2])

		// -----------
		// Export Data
		// -----------
	} else if len(argsWithoutProg) == 1 && argsWithoutProg[0] == "export-data" {
		readDumpAndPrivateKeysAndEncrypt()
		// -------
		// Default
		// -------
	} else { // default
		RunServer()
	}
}
