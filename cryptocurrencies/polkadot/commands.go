package polkadot

import (
	_ "github.com/helloyi/go-waitgroup"
)

// todo: make parallel
func UpdatePolkadotWalletBalances() {
	println("UpdatePolkadotWalletBalances start")

	wallets, err := FindAllPolkadotWallets()
	if err != nil {
		panic(err)
	}
	for _, wallet := range wallets {
		wallet.UpdateBalance()
	}

	println("UpdatePolkadotWalletBalances done")
}
