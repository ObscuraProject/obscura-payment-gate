package monero

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gocraft/web"
	"github.com/ObscuraProject/obscura-payment-gate/apis/monero"
	"github.com/ObscuraProject/obscura-payment-gate/settings"
)

type ViewListMoneroWalletsWithPositiveBalanceResponse struct {
	DisplaMoneroWallets []DisplayMoneroWallet `json:"wallets"`
	SummaryBalance      float64               `json:"balance"`
}

func ViewCreateMoneroWallet(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	walletType := "escrow"
	if r.FormValue("type") != "" {
		walletType = r.FormValue("type")
	}

	wallet, err := CreateMoneroWallet(walletType)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsBitcoinWallet, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsBitcoinWallet)
}

func ViewShowMoneroWallet(w web.ResponseWriter, r *web.Request) {
	wallet, err := FindMoneroWalletByPublicKey(r.PathParams["address"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wb, err := wallet.UpdateBalance()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wb.MoneroWallet = *wallet

	ds := wb.DisplayMoneroWallet()
	ds.Type = wallet.Type
	ds.CreatedAt = wallet.CreatedAt
	jsMoneroWallet, err := json.MarshalIndent(ds, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsMoneroWallet)
}

func ViewListMoneroWalletsWithPositiveBalance(w web.ResponseWriter, r *web.Request) {
	walletsResponse := ViewListMoneroWalletsWithPositiveBalanceResponse{}
	wallets, err := FindMoneroWalletsWithNonZeroBalance()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for _, wallet := range wallets {
		ds := wallet.DisplayMoneroWallet()
		walletsResponse.DisplaMoneroWallets = append(walletsResponse.DisplaMoneroWallets, ds)
		walletsResponse.SummaryBalance += ds.Balance
	}

	jsonBitcoinWallets, err := json.MarshalIndent(walletsResponse, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBitcoinWallets)
}

func ViewSendMonero(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var (
		addressFrom  = r.PathParams["address_from"]
		paymentsForm = r.FormValue("payments")
		payments     = []XMRPayment{}
	)

	err = json.Unmarshal([]byte(paymentsForm), &payments)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wallet, err := FindMoneroWalletByPublicKey(addressFrom)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if len(payments) == 0 {
		http.Error(w, "empty payments list", 500)
		return
	}

	var result []monero.MoneroTransferResult

	if payments[0].Amount == 0.0 && payments[0].Percent != 0.0 {
		result, err = wallet.SendWithPercentSplit(payments)
	} else if payments[0].Amount != 0.0 && payments[0].Percent == 0.0 {
		result, err = wallet.SendWithAmountSplit(payments)
	} else {
		err = errors.New("Either both amount are percent are 0 or not 0")
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsResp, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResp)
}

func ViewUpdateMoneroWallets(w web.ResponseWriter, r *web.Request) {
	go UpdateMoneroWalletBalances()
}

func ViewSweepMoneroWallet(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	addressFrom := r.PathParams["address_from"]

	wallet, err := FindMoneroWalletByPublicKey(addressFrom)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = wallet.Sweep(settings.PAYMENT_GATE_SETTINGS.XMRCommissionWallet)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
