package ethereum

import (
	"github.com/ObscuraProject/obscura-payment-gate/settings"
)

var (
	PAYMENT_GATE_SETTINGS = settings.LoadPaymentGateSettings()
	WEI_IN_ETH            = float64(1e18)
)
