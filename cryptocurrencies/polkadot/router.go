package polkadot

import (
	"github.com/gocraft/web"
)

func ConfigureRouter(router *web.Router) {
	router.Get("/wallets/update", ViewUpdatePolkadotWallets)
	router.Get("/wallets", ViewListPolkadotWalletsWithPositiveBalance)
	router.Post("/wallets/:address/send_dot", ViewSendFromPolkadotWallet)
	router.Get("/wallets/:address", ViewShowPolkadotWallet)
	router.Post("/wallets/new", ViewCreatePolkadotWallet)
	// Escrows
	router.Post("/wallets/:address/escrow_mint", ViewMintPolkadotEscrow)
	router.Get("/escrows/:escrow_id", ViewPolkadotEscrowInfo)
	router.Post("/escrows/:escrow_id/release", ViewReleasePolkadotEscrow)
	router.Post("/escrows/:escrow_id/cancel", ViewCancelPolkadotEscrow)
	router.Get("/wallets/:address/escrows", ViewListPolkadotEscrowsForAddress)
	router.Get("/wallets/:address/deals", ViewListPolkadotDealsForAddress)
	// Crowdloans
	router.Get("/loans/:loan_id", ViewPolkadotLoanInfo)
	router.Get("/wallets/:address/loans", ViewListPolkadotCrowdloansForAddress)
	router.Get("/wallets/:address/loans/:loan_id", ViewListPolkadotCrowdloansForAddressWithId)
	router.Get("/wallets/:address/lends", ViewListPolkadotCrowdloanLendsForAddress)
	router.Post("/wallets/:address/crowdloan_mint", ViewMintPolkadotLoan)
	router.Post("/wallets/:address/crowdloan_fund", ViewFundPolkadotLoan)
	router.Post("/wallets/:address/crowdloan_withdraw", ViewWithdrawPolkadotLoan)
	router.Post("/wallets/:address/crowdloan_payback", ViewPaybackPolkadotLoan)
	router.Post("/wallets/:address/crowdloan_payout", ViewPayoutPolkadotLoan)
}
