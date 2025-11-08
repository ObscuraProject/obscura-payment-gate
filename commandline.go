package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/bitcoin"
	"github.com/ObscuraProject/obscura-payment-gate/cryptocurrencies/ethereum"
)

type DumpUser struct {
	Username           string `json:"username"`
	PGPPublicKey       string `json:"pgp"`
	BitcoinPublicKey   string `json:"bitcoin_public_key"`
	EthereumPublicKey  string `json:"ethereum_public_key"`
	BitcoinPrivateKey  string `json:"bitcoin_private_key"`
	EthereumPrivateKey string `json:"ethereum_private_key"`
}

func encryptText(encryptionText, publicKey string) (string, error) {
	entitylist, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(publicKey))
	if err != nil {
		return "", err
	}

	encbuf := bytes.NewBuffer(nil)
	w, err := armor.Encode(encbuf, "PGP MESSAGE", nil)
	if err != nil {
		return "", err
	}

	plaintext, err := openpgp.Encrypt(w, entitylist, nil, nil, nil)
	if err != nil {
		return "", err
	}

	message := []byte(encryptionText)
	_, err = plaintext.Write(message)

	plaintext.Close()
	w.Close()

	return fmt.Sprintf("%s", encbuf), nil
}

func readDumpAndPrivateKeysAndEncrypt() {
	// Open our jsonFile
	jsonFile, err := os.Open("all_users.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var dumpUsers []DumpUser
	json.Unmarshal([]byte(byteValue), &dumpUsers)

	for _, user := range dumpUsers {
		if user.BitcoinPublicKey == "" || user.EthereumPublicKey == "" {
			continue
		}
		wB, _ := bitcoin.FindBitcoinWalletByPublicKey(user.BitcoinPublicKey)
		wE, _ := ethereum.FindEthereumWalletByPublicKey(user.EthereumPublicKey)
		user.EthereumPrivateKey = wE.PrivateKey
		user.BitcoinPrivateKey = wB.WIFC

		byteJson, _ := json.Marshal(user)
		jsonUser := string(byteJson)
		encryptedText, _ := encryptText(jsonUser, user.PGPPublicKey)

		data := []byte(user.Username)
		hash := sha256.Sum256(data)
		usernameHash := fmt.Sprintf("%x", hash[:])
		f, err := os.OpenFile("credentials/"+usernameHash, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
		if err != nil {
			panic(err)
		}

		f.WriteString(encryptedText)
	}
}
