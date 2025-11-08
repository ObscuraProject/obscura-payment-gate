package polkadot

import (
	"errors"

	polkadotapi "github.com/ObscuraProject/obscura-payment-gate/apis/polkadot"
)

func (w PolkadotWallet) TransferDot(payments []DOTPayment) (DOTPaymentResult, error) {
	_, err := w.UpdateBalance()
	if err != nil {
		return DOTPaymentResult{}, err
	}

	if len(payments) != 1 {
		return DOTPaymentResult{}, errors.New("payments for should contain only one item")
	}

	payment := payments[0]

	var amount float64
	if payment.Amount > 0 {
		amount = float64(payment.Amount)
	} else if payment.Percent > 0 {
		balance, err := w.CurrentBalance()
		if err == nil {
			amount = float64(balance.FreeBalance * payment.Percent)
		}
	}

	transactionResult, err := polkadotapi.SendDot(payment.Address, amount, w.Mnemonic)
	if err != nil {
		return DOTPaymentResult{}, err
	}

	return DOTPaymentResult{
		Hash:       transactionResult.Hash,
		WalletFrom: w.PublicKey,
		WalletTo:   payment.Address,
		Amount:     amount,
	}, nil
}
