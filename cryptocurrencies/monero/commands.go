package monero

import (
	_ "github.com/helloyi/go-waitgroup"
)

// todo: make parallel
func UpdateMoneroWalletBalances() {
	println("UpdateMoneroWalletBalances start")

	wallets, err := FindAllMoneroWallets()
	if err != nil {
		panic(err)
	}
	for _, wallet := range wallets {
		wallet.UpdateBalance()
	}

	println("UpdateMoneroWalletBalances done")
}
