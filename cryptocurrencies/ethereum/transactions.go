package ethereum

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/onrik/ethrpc"

	"github.com/ObscuraProject/obscura-payment-gate/ethereum/contracts/golang/registry"
)

/*
	Globals
*/

func (w EthereumWallet) TransferEth(payments []ETHPayment) (ETHPaymentResult, error) {
	_, err := w.UpdateBalance()
	if err != nil {
		return ETHPaymentResult{}, err
	}

	// send amount to a single wallet -- basic transfer in Ethereum
	if len(payments) == 1 {
		var amount uint64
		if payments[0].Amount > 0 {
			amount = EtherToWei(payments[0].Amount)
		} else if payments[0].Percent > 0 {
			balance, err := w.CurrentBalance()
			if err == nil {
				amount = uint64(balance.BalanceETH * WEI_IN_ETH * payments[0].Percent)
			}
		}
		return w.transfer_eth(payments[0].Address, amount)
	}

	// send amount to multiple wallets -- via a deployed smart contract
	var (
		_addresses = []common.Address{}
		_percents  = []*big.Int{}
	)

	for _, payment := range payments {
		address := common.HexToAddress(payment.Address)
		percent := big.NewInt(int64(payment.Percent * 100))

		_addresses = append(_addresses, address)
		_percents = append(_percents, percent)
	}

	return w.contractMultitransferETH(_addresses, _percents)
}

func EtherToWei(eth float64) uint64 {
	bfWeiInEth := big.NewFloat(0.0).SetPrec(64).SetFloat64(WEI_IN_ETH)
	bfBalance := big.NewFloat(0.0).SetPrec(64).SetFloat64(eth)
	bfTerm1 := big.NewFloat(0.0).Mul(bfBalance, bfWeiInEth)
	// precision problems when doing otherwise
	ff, _ := strconv.ParseFloat(bfTerm1.String(), 64)
	return uint64(ff)
}

func WeiToEther(wei *big.Int) float64 {
	bfWeiInEth := big.NewFloat(WEI_IN_ETH)
	bfWei := big.NewFloat(0.0).SetInt(wei)
	f64, _ := new(big.Float).SetPrec(128).Quo(bfWei, bfWeiInEth).Float64()
	f64 = math.Floor(f64/0.000001) * 0.000001
	return f64
}

/*
	Ethereum Basic Transfer
*/

func (w EthereumWallet) transfer_eth(to string, value uint64) (ETHPaymentResult, error) {
	var (
		buff          bytes.Buffer
		chainId       = big.NewInt(1)
		recipientAddr = common.HexToAddress(to)
		gasLimit      = uint64(21000)
		etx           = ETHPaymentResult{
			WalletFrom: w.PublicKey,
			WalletTo:   to,
			Amount:     value,
		}
	)

	// hacky way to convert uint64 to big.Int since straight type conv adds unnessesary sign
	biAmount := new(big.Int).SetUint64(value)

	senderPrivKey, err := crypto.HexToECDSA(w.PrivateKey[2:])
	if err != nil {
		return etx, err
	}

	tcount, err := ethereumAPI.EthGetTransactionCount(w.PublicKey, "latest")
	if err != nil {
		return etx, err
	}

	gasPrice, err := ethereumAPI.EthGasPrice()
	if err != nil {
		return etx, err
	}

	valueForEst := value / 1000 * 999
	biAmountEst := new(big.Int).SetUint64(valueForEst)
	gas, err := ethereumAPI.EthEstimateGas(ethrpc.T{
		From:     w.PublicKey,
		To:       to,
		Value:    biAmountEst,
		GasPrice: &gasPrice,
		Nonce:    tcount,
	})
	if err != nil {
		return etx, err
	}
	biGas := big.NewInt(int64(gas * 2))

	var (
		nonce             = uint64(tcount)
		biTransactionCost = big.NewInt(0).Mul(&gasPrice, biGas)
		txAmount          = big.NewInt(0).Sub(biAmount, biTransactionCost)
		signer            = types.NewEIP155Signer(chainId)
	)

	tx := types.NewTransaction(nonce, recipientAddr, txAmount, gasLimit, &gasPrice, nil)

	signedTx, err := types.SignTx(tx, signer, senderPrivKey)
	if err != nil {
		return etx, err
	}

	err = signedTx.EncodeRLP(&buff)
	if err != nil {
		return etx, err
	}

	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	etx.RawTx = fmt.Sprintf("0x%x", buff.Bytes())
	hash := signedTx.Hash().Hex()

	etx.Amount = txAmount.Uint64()
	etx.Hash = hash
	etx.Cost = biTransactionCost.Uint64()

	return etx, err
}

/*
	Smart Contract Multitransfer
*/

func (w EthereumWallet) contractMultitransferETH(_addresses []common.Address, _percents []*big.Int) (ETHPaymentResult, error) {
	ethPaymentResult := ETHPaymentResult{
		WalletFrom:     w.PublicKey,
		WalletTo:       registry.TOCHKA_ESCROW_PAYMENT_ADDRESS,
		Contract:       "TochkaEscrowPayment",
		ContractMethod: "multitransfer",
	}

	if len(_addresses) != len(_percents) {
		return ethPaymentResult, errors.New("len(_addresses) should be equal to len(_percents)")
	}

	cost := uint64(100000)

	// Smart Contract
	address := common.HexToAddress(registry.TOCHKA_ESCROW_PAYMENT_ADDRESS)
	tochkaEscrowPayment, err := NewTochkaEscrowPayment(address, ethClient)
	if err != nil {
		return ethPaymentResult, err
	}

	// Transactor
	privateKey, err := crypto.HexToECDSA(w.PrivateKey[2:])
	if err != nil {
		return ethPaymentResult, err
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return ethPaymentResult, err
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return ethPaymentResult, err
	}

	gasPrice.Mul(gasPrice, big.NewInt(2))

	balance, err := w.CurrentBalance()
	if err != nil {
		return ethPaymentResult, err
	}

	var (
		bfCost  = big.NewInt(int64(cost))
		bfTerm1 = big.NewInt(0).SetUint64(EtherToWei(balance.BalanceETH))
		bfTerm2 = big.NewInt(0).Mul(bfCost, gasPrice)
		bfValue = big.NewInt(0).Sub(bfTerm1, bfTerm2)
	)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
	session := &TochkaEscrowPaymentSession{
		Contract: tochkaEscrowPayment,
		CallOpts: bind.CallOpts{
			Pending: false,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			Value:    bfValue,
			GasLimit: uint64(cost),
			Nonce:    big.NewInt(int64(nonce)),
		},
	}

	tx, err := session.Multitransfer(_addresses, _percents)
	if err != nil {
		return ethPaymentResult, err
	}

	ethPaymentResult.Cost = tx.Cost().Uint64()
	ethPaymentResult.Amount = tx.Value().Uint64()
	ethPaymentResult.Hash = fmt.Sprintf("0x%x", tx.Hash())
	ethPaymentResult.RawTx = string(tx.Data())

	if !strings.Contains(w.Type, "_dust") {
		w.Type += "_dust"
	}
	return ethPaymentResult, w.Save()
}
