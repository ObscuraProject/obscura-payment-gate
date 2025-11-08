package polkadot

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type PolkadotWallet struct {
	PublicKey string `json:"public_key"`
	Mnemonic  string `json:"mnemonic"`
}

func CreatePolkadotWallet() (PolkadotWallet, error) {
	polkadotWallet := PolkadotWallet{}
	output, err := exec.Command("node", "--no-warnings", "nodejs/polkadot-escrow.mjs", "--mode", "create_wallet").Output()
	if err != nil {
		return polkadotWallet, err
	}

	err = json.Unmarshal(output, &polkadotWallet)
	return polkadotWallet, err
}

type PolkadotWalletBalance struct {
	Free     float64 `json:"free"`
	Reserved float64 `json:"reserved"`
	Frozen   float64 `json:"frozen"`
}

func GetWalletBalance(publicKey string) (PolkadotWalletBalance, error) {
	polkadotWalletBalance := PolkadotWalletBalance{}
	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-escrow.mjs",
		"--mode", "get_balance",
		"--wallet", publicKey,
	).Output()
	if err != nil {
		return polkadotWalletBalance, err
	}

	err = json.Unmarshal([]byte(output), &polkadotWalletBalance)
	return polkadotWalletBalance, err
}

type PolkadotTransactionResult struct {
	Hash    string `json:"hash"`
	Address string `json:"address"`
}

func ReleasePolkadotEscrow(escrowId uint64) (PolkadotTransactionResult, error) {
	polkadotTransactionResult := PolkadotTransactionResult{}

	cmd := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-escrow.mjs",
		"--mode", "release_escrow",
		"--escrow_id", fmt.Sprintf("%d", escrowId),
	)

	output, err := cmd.Output()
	if err != nil {
		return polkadotTransactionResult, err
	}

	err = json.Unmarshal([]byte(output), &polkadotTransactionResult)
	return polkadotTransactionResult, err
}

func CancelPolkadotEscrow(escrowId uint64) (PolkadotTransactionResult, error) {
	polkadotTransactionResult := PolkadotTransactionResult{}

	cmd := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-escrow.mjs",
		"--mode", "cancel_escrow",
		"--escrow_id", fmt.Sprintf("%d", escrowId),
	)

	output, err := cmd.Output()
	if err != nil {
		return polkadotTransactionResult, err
	}

	err = json.Unmarshal([]byte(output), &polkadotTransactionResult)
	return polkadotTransactionResult, err
}

type PolkadotEscrowMintResult struct {
	Hash    string `json:"hash"`
	Address string `json:"address"`
	Id      uint64 `json:"escrow_id"`
}

func MintEscrow(sellerAddress string, cid string, value float64, userWalletMnemonic string) (PolkadotEscrowMintResult, error) {
	polkadotMintTransactionResult := PolkadotEscrowMintResult{}
	plankValue := uint64(value * 1e10)

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-escrow.mjs",
		"--mode", "mint_escrow",
		"--seller_address", sellerAddress,
		"--cid", cid,
		"--value", fmt.Sprintf("%d", plankValue),
		"--user_wallet_mnemonic", userWalletMnemonic,
	).Output()
	if err != nil {
		return polkadotMintTransactionResult, err
	}

	err = json.Unmarshal(output, &polkadotMintTransactionResult)
	return polkadotMintTransactionResult, err
}

type PolkadotCrowdloanMintResult struct {
	Hash    string `json:"hash"`
	Address string `json:"address"`
	Id      uint64 `json:"crowdloan_id"`
}

func MintCrowdloan(goalValue float64, goalDate, weeklyInterest uint64, userWalletMnemonic string) (PolkadotCrowdloanMintResult, error) {
	polkadotMintTransactionResult := PolkadotCrowdloanMintResult{}
	plankValue := uint64(goalValue * 1e10)

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-crowdloan.mjs",
		"--mode", "mint_crowdloan",
		"--goal_value", fmt.Sprintf("%d", plankValue),
		"--goal_date", fmt.Sprintf("%d", goalDate*1000),
		"--weekly_interest", fmt.Sprintf("%d", weeklyInterest),
		"--user_wallet_mnemonic", userWalletMnemonic,
	).Output()
	if err != nil {
		return polkadotMintTransactionResult, err
	}

	err = json.Unmarshal(output, &polkadotMintTransactionResult)
	return polkadotMintTransactionResult, err
}

type PolkadotEscrowStatus struct {
	Status string `json:"status"`
	Time   uint64 `json:"timestamp"`
}

type PolkadotEscrow struct {
	EscrowStatus   PolkadotEscrowStatus `json:"escrow_status"`
	Id             uint64               `json:"id"`
	MarketplaceId  uint64               `json:"marketplace_id"`
	Owner          string               `json:"owner"`
	SellerAddress  string               `json:"seller_address"`
	ShippingStatus PolkadotEscrowStatus `json:"shipping_status"`
	Uri            string               `json:"uri"`
	Value          uint64               `json:"value"`
}

func EscrowInfo(escrowId uint64) (PolkadotEscrow, error) {
	polkadotEscrowInfoResult := PolkadotEscrow{}

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-escrow.mjs",
		"--mode", "escrow_info",
		"--escrow_id", fmt.Sprintf("%d", escrowId),
	).Output()
	if err != nil {
		return polkadotEscrowInfoResult, err
	}

	err = json.Unmarshal(output, &polkadotEscrowInfoResult)
	return polkadotEscrowInfoResult, err
}

type PolkdadotLoan struct {
	CloseDate       *uint                `json:"close_date"`
	ClosingAmount   uint                 `json:"closing_amount"`
	Value           float64              `json:"value"`
	EscrowStatus    PolkadotEscrowStatus `json:"status"`
	GoalValue       float64              `json:"goal_value"`
	GoalDate        uint64               `json:"goal_date"`
	Id              uint64               `json:"id"`
	NumberOfLenders uint                 `json:"n_lenders"`
	Owner           string               `json:"owner"`
	WeeklyIterest   uint64               `json:"weekly_interest"`
}

func CrowdloanInfo(escrowId uint64) (PolkdadotLoan, error) {
	polkadotCrowdloanInfoResult := PolkdadotLoan{}

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-crowdloan.mjs",
		"--mode", "crowdloan_info",
		"--crowdloan_id", fmt.Sprintf("%d", escrowId),
	).Output()
	if err != nil {
		return polkadotCrowdloanInfoResult, err
	}

	err = json.Unmarshal(output, &polkadotCrowdloanInfoResult)
	return polkadotCrowdloanInfoResult, err
}

func EscrowsForAddress(address string) ([]PolkadotEscrow, error) {
	polkadotEscrowInfoResult := []PolkadotEscrow{}

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-escrow.mjs",
		"--mode", "list_wallet_escrows",
		"--wallet", address,
	).Output()
	if err != nil {
		return polkadotEscrowInfoResult, err
	}

	err = json.Unmarshal(output, &polkadotEscrowInfoResult)
	return polkadotEscrowInfoResult, err
}

func DealsForAddress(address string) ([]PolkadotEscrow, error) {
	polkadotEscrowInfoResult := []PolkadotEscrow{}

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-escrow.mjs",
		"--mode", "list_wallet_deals",
		"--wallet", address,
	).Output()
	if err != nil {
		return polkadotEscrowInfoResult, err
	}

	err = json.Unmarshal(output, &polkadotEscrowInfoResult)
	return polkadotEscrowInfoResult, err
}

func SendDot(addressTo string, value float64, userWalletMnemonic string) (PolkadotTransactionResult, error) {
	polkadotTransactionResult := PolkadotTransactionResult{}
	plankValue := uint64(value * 1e10)

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-escrow.mjs",
		"--mode", "send_dot",
		"--user_wallet_mnemonic", userWalletMnemonic,
		"--amount", fmt.Sprintf("%d", plankValue),
		"--to", addressTo,
	).Output()
	if err != nil {
		return polkadotTransactionResult, err
	}

	err = json.Unmarshal([]byte(output), &polkadotTransactionResult)
	return polkadotTransactionResult, err
}

func FundCrowdloan(crowdloanId uint64, value float64, userWalletMnemonic string) (PolkadotTransactionResult, error) {
	result := PolkadotTransactionResult{}
	plankValue := uint64(value * 1e10)

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-crowdloan.mjs",
		"--mode", "fund_crowdloan",
		"--crowdloan_id", fmt.Sprintf("%d", crowdloanId),
		"--value", fmt.Sprintf("%d", plankValue),
		"--user_wallet_mnemonic", userWalletMnemonic,
	).Output()
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(output, &result)
	return result, err
}

func WithdrawCrowdloan(crowdloanId uint64, userWalletMnemonic string) (PolkadotTransactionResult, error) {
	result := PolkadotTransactionResult{}

	println(
		"node", "--no-warnings", "nodejs/polkadot-crowdloan.mjs",
		"--mode", "withdraw_crowdloan",
		"--crowdloan_id", fmt.Sprintf("%d", crowdloanId),
		"--user_wallet_mnemonic", userWalletMnemonic,
	)

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-crowdloan.mjs",
		"--mode", "withdraw_crowdloan",
		"--crowdloan_id", fmt.Sprintf("%d", crowdloanId),
		"--user_wallet_mnemonic", userWalletMnemonic,
	).Output()
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(output, &result)
	return result, err
}

func PaybackCrowdloan(crowdloanId uint64, userWalletMnemonic string) (PolkadotTransactionResult, error) {
	result := PolkadotTransactionResult{}

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-crowdloan.mjs",
		"--mode", "payback_crowdloan",
		"--crowdloan_id", fmt.Sprintf("%d", crowdloanId),
		"--user_wallet_mnemonic", userWalletMnemonic,
	).Output()
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(output, &result)
	return result, err
}

func PayoutCrowdloan(crowdloanId uint64, userWalletMnemonic string) (PolkadotTransactionResult, error) {
	result := PolkadotTransactionResult{}

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-crowdloan.mjs",
		"--mode", "payout_crowdloan",
		"--crowdloan_id", fmt.Sprintf("%d", crowdloanId),
		"--user_wallet_mnemonic", userWalletMnemonic,
	).Output()
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(output, &result)
	return result, err
}

func CrowdloansForAddress(address string) ([]PolkdadotLoan, error) {
	loansResults := []PolkdadotLoan{}

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-crowdloan.mjs",
		"--mode", "list_wallet_crowdloans",
		"--wallet", address,
	).Output()
	if err != nil {
		return loansResults, err
	}

	err = json.Unmarshal(output, &loansResults)
	return loansResults, err
}

type PolkdadotLoanLend struct {
	Id          uint64        `json:"id"`
	CrowdloanId uint64        `json:"crowdloan_id"`
	LendValue   float64       `json:"lend_value"`
	Index       uint64        `json:"index"`
	Loan        PolkdadotLoan `json:"loan"`
}

func CrowdloanLendsForAddress(address string) ([]PolkdadotLoanLend, error) {
	loansResults := []PolkdadotLoanLend{}

	output, err := exec.Command(
		"node", "--no-warnings", "nodejs/polkadot-crowdloan.mjs",
		"--mode", "list_wallet_lends",
		"--wallet", address,
	).Output()
	if err != nil {
		return loansResults, err
	}

	err = json.Unmarshal(output, &loansResults)
	return loansResults, err
}
