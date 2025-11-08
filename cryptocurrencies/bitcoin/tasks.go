package bitcoin

import (
	"time"

	"github.com/ObscuraProject/obscura-payment-gate/settings"
)

func TaskUpdateBitcoinWalletBalances() {
	c := time.Tick(24 * time.Hour)
	for range c {
		UpdateBitcoinWalletBalances()
	}
}

func TaskUpdateBitcoinTransactionFee() {
	UpdateBTCTransactionFee()
	c := time.Tick(1 * time.Minute)
	for range c {
		go UpdateBTCTransactionFee()
	}
}

func init() {
	if !settings.APPLICATION_SETTINGS.Debug {
		go TaskUpdateBitcoinWalletBalances()
		go TaskUpdateBitcoinTransactionFee()
	}
}
