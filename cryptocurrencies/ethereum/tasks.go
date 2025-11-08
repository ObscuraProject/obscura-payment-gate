package ethereum

import (
	"time"

	"github.com/ObscuraProject/obscura-payment-gate/settings"
)

func TaskUpdateEthereumWalletBalances() {
	c := time.Tick(24 * time.Hour)
	for range c {
		UpdateEthereumWalletBalances()
	}
}

func init() {
	if !settings.APPLICATION_SETTINGS.Debug {
		go TaskUpdateEthereumWalletBalances()
	}
}
