package monero

import (
	"github.com/gocraft/web"
)

func ConfigureRouter(router *web.Router) {
	router.Get("/update", ViewUpdateMoneroWallets)
	router.Get("/wallets", ViewListMoneroWalletsWithPositiveBalance)
	router.Get("/wallets/:address", ViewShowMoneroWallet)
	router.Post("/wallets/:address_from/sweep", ViewSweepMoneroWallet)
	router.Post("/wallets/:address_from/send", ViewSendMonero)
	router.Post("/wallets/new", ViewCreateMoneroWallet)
}
