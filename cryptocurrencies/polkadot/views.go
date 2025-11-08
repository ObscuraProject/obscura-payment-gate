package polkadot

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gocraft/web"

	polkadotapi "github.com/ObscuraProject/obscura-payment-gate/apis/polkadot"
)

type ViewListPolkadotWalletsWithPositiveBalanceResponse struct {
	DisplaPolkadotWallets []DisplayPolkadotWallet `json:"wallets"`
	SummaryBalance        float64                 `json:"balance"`
}

func ViewCreatePolkadotWallet(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wallet, err := CreatePolkadotWallet()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsPolkadotWallet, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsPolkadotWallet)
}

func ViewShowPolkadotWallet(w web.ResponseWriter, r *web.Request) {
	wallet, err := FindPolkadotWalletByPublicKey(r.PathParams["address"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wb, err := wallet.UpdateBalance()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wb.PolkadotWallet = *wallet

	ds := wb.DisplayPolkadotWallet()
	ds.CreatedAt = wallet.CreatedAt
	jsPolkadotWallet, err := json.MarshalIndent(ds, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsPolkadotWallet)
}

func ViewListPolkadotWalletsWithPositiveBalance(w web.ResponseWriter, r *web.Request) {
	walletsResponse := ViewListPolkadotWalletsWithPositiveBalanceResponse{}
	wallets, err := FindPolkadotWalletsWithNonZeroBalance()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for _, wallet := range wallets {
		ds := wallet.DisplayPolkadotWallet()
		walletsResponse.DisplaPolkadotWallets = append(walletsResponse.DisplaPolkadotWallets, ds)
		walletsResponse.SummaryBalance += ds.FreeBalance
	}

	jsonPolkadotWallets, err := json.MarshalIndent(walletsResponse, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPolkadotWallets)
}

func ViewReleasePolkadotEscrow(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	escrowIdString := r.PathParams["escrow_id"]
	escrowId, err := strconv.ParseUint(escrowIdString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	transactionResult, err := polkadotapi.ReleasePolkadotEscrow(escrowId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsResp, err := json.MarshalIndent(transactionResult, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResp)
}

func ViewMintPolkadotEscrow(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var (
		escrowMintMetadata = DOTEscrowMintMetadata{}
		metadata           = r.FormValue("metadata")
		addressFrom        = r.PathParams["address"]
	)

	err = json.Unmarshal([]byte(metadata), &escrowMintMetadata)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wallet, err := FindPolkadotWalletByPublicKey(addressFrom)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	balance, err := wallet.UpdateBalance()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if balance.FreeBalance < escrowMintMetadata.Value {
		http.Error(w, "insufficient balance", 500)
		return
	}

	transactionResult, err := polkadotapi.MintEscrow(
		escrowMintMetadata.SellerAddress,
		escrowMintMetadata.Cid,
		escrowMintMetadata.Value,
		wallet.Mnemonic,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsResp, err := json.MarshalIndent(transactionResult, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResp)
}

func ViewCancelPolkadotEscrow(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	escrowIdString := r.PathParams["escrow_id"]
	escrowId, err := strconv.ParseUint(escrowIdString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	transactionResult, err := polkadotapi.CancelPolkadotEscrow(escrowId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsResp, err := json.MarshalIndent(transactionResult, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResp)
}

func ViewUpdatePolkadotWallets(w web.ResponseWriter, r *web.Request) {
	go UpdatePolkadotWalletBalances()
}

func ViewPolkadotEscrowInfo(w web.ResponseWriter, r *web.Request) {
	escrowId, err := strconv.ParseUint(r.PathParams["escrow_id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	escrowInfo, err := polkadotapi.EscrowInfo(escrowId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsEscrowInfo, err := json.MarshalIndent(escrowInfo, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsEscrowInfo)
}

func ViewListPolkadotEscrowsForAddress(w web.ResponseWriter, r *web.Request) {
	wallet, err := polkadotapi.EscrowsForAddress(r.PathParams["address"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsPolkadotWallet, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsPolkadotWallet)
}

func ViewListPolkadotDealsForAddress(w web.ResponseWriter, r *web.Request) {
	wallet, err := polkadotapi.DealsForAddress(r.PathParams["address"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsPolkadotWallet, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsPolkadotWallet)
}

func ViewSendFromPolkadotWallet(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var (
		addressFrom  = r.PathParams["address"]
		paymentsForm = r.FormValue("payments")
		payments     = []DOTPayment{}
	)

	err = json.Unmarshal([]byte(paymentsForm), &payments)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wallet, err := FindPolkadotWalletByPublicKey(addressFrom)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	results, err := wallet.TransferDot(payments)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsonResult, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

func ViewPolkadotLoanInfo(w web.ResponseWriter, r *web.Request) {
	escrowId, err := strconv.ParseUint(r.PathParams["loan_id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	escrowInfo, err := polkadotapi.CrowdloanInfo(escrowId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsEscrowInfo, err := json.MarshalIndent(escrowInfo, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsEscrowInfo)
}

type DOTFundCrowdloanMetadata struct {
	CrowdloanId uint64  `json:"crowdloan_id"`
	Value       float64 `json:"value"`
}

func ViewFundPolkadotLoan(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var (
		escrowFundMetadata = DOTFundCrowdloanMetadata{}
		metadata           = r.FormValue("metadata")
		addressFrom        = r.PathParams["address"]
	)

	err = json.Unmarshal([]byte(metadata), &escrowFundMetadata)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wallet, err := FindPolkadotWalletByPublicKey(addressFrom)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	transactionResult, err := polkadotapi.FundCrowdloan(
		escrowFundMetadata.CrowdloanId,
		escrowFundMetadata.Value,
		wallet.Mnemonic,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsResp, err := json.MarshalIndent(transactionResult, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResp)
}

type DOTWithdawCrowdloanMetadata struct {
	CrowdloanId uint64 `json:"crowdloan_id"`
}

func ViewWithdrawPolkadotLoan(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var (
		escrowFundMetadata = DOTWithdawCrowdloanMetadata{}
		metadata           = r.FormValue("metadata")
		addressFrom        = r.PathParams["address"]
	)

	err = json.Unmarshal([]byte(metadata), &escrowFundMetadata)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wallet, err := FindPolkadotWalletByPublicKey(addressFrom)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	transactionResult, err := polkadotapi.WithdrawCrowdloan(
		escrowFundMetadata.CrowdloanId,
		wallet.Mnemonic,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsResp, err := json.MarshalIndent(transactionResult, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResp)
}

type DOTPaybackCrowdloanMetadata struct {
	CrowdloanId uint64 `json:"crowdloan_id"`
}

type DOTPayoutCrowdloanMetadata struct {
	CrowdloanId uint64 `json:"crowdloan_id"`
}

func ViewPaybackPolkadotLoan(w web.ResponseWriter, r *web.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var (
		paybackMetadata = DOTPaybackCrowdloanMetadata{}
		metadata        = r.FormValue("metadata")
		addressFrom     = r.PathParams["address"]
	)

	err = json.Unmarshal([]byte(metadata), &paybackMetadata)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wallet, err := FindPolkadotWalletByPublicKey(addressFrom)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	transactionResult, err := polkadotapi.PaybackCrowdloan(
		paybackMetadata.CrowdloanId,
		wallet.Mnemonic,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsResp, err := json.MarshalIndent(transactionResult, "", "  ")
	if err != nil {
		println("err 5")
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResp)
}

func ViewPayoutPolkadotLoan(w web.ResponseWriter, r *web.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var (
		paybackMetadata = DOTPayoutCrowdloanMetadata{}
		metadata        = r.FormValue("metadata")
		addressFrom     = r.PathParams["address"]
	)

	err = json.Unmarshal([]byte(metadata), &paybackMetadata)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wallet, err := FindPolkadotWalletByPublicKey(addressFrom)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	transactionResult, err := polkadotapi.PayoutCrowdloan(
		paybackMetadata.CrowdloanId,
		wallet.Mnemonic,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsResp, err := json.MarshalIndent(transactionResult, "", "  ")
	if err != nil {
		println("err 5")
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResp)
}

func ViewMintPolkadotLoan(w web.ResponseWriter, r *web.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var (
		escrowMintMetadata = DOTCrowdloanMintMetadata{}
		metadata           = r.FormValue("metadata")
		addressFrom        = r.PathParams["address"]
	)

	err = json.Unmarshal([]byte(metadata), &escrowMintMetadata)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	wallet, err := FindPolkadotWalletByPublicKey(addressFrom)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	transactionResult, err := polkadotapi.MintCrowdloan(
		escrowMintMetadata.GoalValue,
		escrowMintMetadata.GoalDate,
		uint64(escrowMintMetadata.WeeklyInterest),
		wallet.Mnemonic,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsResp, err := json.MarshalIndent(transactionResult, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResp)
}

func ViewListPolkadotCrowdloansForAddress(w web.ResponseWriter, r *web.Request) {
	wallet, err := polkadotapi.CrowdloansForAddress(r.PathParams["address"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsPolkadotWallet, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsPolkadotWallet)
}

func ViewListPolkadotCrowdloansForAddressWithId(w web.ResponseWriter, r *web.Request) {
	wallet, err := polkadotapi.CrowdloanLendsForAddress(r.PathParams["address"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsPolkadotWallet, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsPolkadotWallet)
}

func ViewListPolkadotCrowdloanLendsForAddress(w web.ResponseWriter, r *web.Request) {

	wallet, err := polkadotapi.CrowdloanLendsForAddress(r.PathParams["address"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsPolkadotWallet, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsPolkadotWallet)
}
