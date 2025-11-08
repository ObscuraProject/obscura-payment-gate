package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/liyue201/erc20-go/erc20"
	"github.com/onrik/ethrpc"
)

var (
	ethClient *ethclient.Client
	usdtToken *erc20.GGToken

	ethereumAPI = ethrpc.NewEthRPC(PAYMENT_GATE_SETTINGS.EthereumAPIURL)
)

func init() {
	var err error
	ethClient, err = ethclient.Dial(PAYMENT_GATE_SETTINGS.EthereumAPIURL)
	if err != nil {
		panic(err)
	}

	usdtToken, _ = erc20.NewGGToken(common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), ethClient)
}
