package electrum

import (
	"context"

	el "github.com/checksum0/go-electrum/electrum"
	"github.com/jasonlvhit/gocron"
	"github.com/ObscuraProject/obscura-payment-gate/settings"
)

type ElectrumAPI struct {
	APIEndpoint string
	Client      *el.Client
}

var api = ElectrumAPI{
	APIEndpoint: settings.PAYMENT_GATE_SETTINGS.ElectrumAPIURL,
	Client:      nil,
}

func init() {
	client, err := el.NewClientTCP(context.TODO(), api.APIEndpoint)
	if err != nil {
		panic(err)
	}

	api.Client = client

	gocron.Every(1).Hour().Do(pingElectrum)
}

func GetClinet() (*el.Client, error) {
	return api.Client, nil
}

func pingElectrum() {
	err := api.Client.Ping(context.TODO())
	if err != nil {
		panic(err)
	}
}
