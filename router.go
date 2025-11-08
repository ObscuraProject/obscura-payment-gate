package main

import (
	"github.com/gocraft/web"

	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies"
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/bitcoin"
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/ethereum"
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/monero"
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/polkadot"
)

func ConfigureRouter(router *web.Router) *web.Router {
	bitcoinRouter := router.Subrouter(Context{}, "/bitcoin")
	bitcoin.ConfigureRouter(bitcoinRouter)

	ethereumRouter := router.Subrouter(Context{}, "/ethereum")
	ethereum.ConfigureRouter(ethereumRouter)

	moneroRouter := router.Subrouter(Context{}, "/monero")
	monero.ConfigureRouter(moneroRouter)

	polkadotRouer := router.Subrouter(Context{}, "/polkadot")
	polkadot.ConfigureRouter(polkadotRouer)

	// Currency
	router.Get("/currency/:base_currency", cryptocurrencies.ViewShowCurrency)
	return router
}
