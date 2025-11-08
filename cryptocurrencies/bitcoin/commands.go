package bitcoin

import (
	"github.com/helloyi/go-waitgroup"
)

func UpdateBitcoinWalletBalances() {
	println("UpdateBitcoinWalletBalances start")

	wg := waitgroup.Create(32)
	wallets, err := FindAllBitcoinWallets()
	if err != nil {
		panic(err)
	}
	for _, wallet := range wallets {
		wg.BlockAdd()
		go func(wallet BitcoinWallet) {
			wallet.UpdateBalance()
			wg.Done()
		}(wallet)
	}
	wg.Wait()

	println("UpdateBitcoinWalletBalances done")
}

func UpdatePremiumBitcoinWalletBalances() {
	println("UpdatePremiumBitcoinWalletBalances start")

	wg := waitgroup.Create(64)
	wallets, err := FindPremiumBitcoinWallets()
	if err != nil {
		panic(err)
	}
	for _, wallet := range wallets {
		wg.BlockAdd()
		go func(wallet BitcoinWallet) {
			defer wg.Done()
			wallet.UpdateBalance()
		}(wallet)
	}
	wg.Wait()

	println("UpdatePremiumBitcoinWalletBalances done")
}

func SweepBitcoinWallet(addressFrom string, addressTo string) {
	w, err := FindBitcoinWalletByPublicKey(addressFrom)
	if err != nil {
		panic(err)
	}

	payments := []BTCPayment{
		{Address: addressTo, Percent: 1.0},
	}

	_, err = w.SendWithPercentSplit(payments)
	if err != nil {
		panic(err)
	}
}

func SweepCommissionAndPremiumBitcoinWallets(addressTo string) {
	wp, err := FindBitcoinWalletsWithNonZeroBalanceByType("premium")
	if err != nil {
		panic(err)
	}

	wa, err := FindBitcoinWalletsWithNonZeroBalanceByType("advertising")
	if err != nil {
		panic(err)
	}

	wallets := BitcoinWallets{}
	for _, wb := range wp {
		if wb.PublicKey != "" {
			wallets = append(wallets, wb.BitcoinWallet)
		}
	}

	for _, wb := range wa {
		if wb.PublicKey != "" {
			wallets = append(wallets, wb.BitcoinWallet)
		}
	}

	payments := []BTCPayment{
		{Address: addressTo, Percent: 1.0},
	}

	_, err = wallets.SendWithPercentSplit(payments)
	if err != nil {
		println(err.Error())
	}
}
