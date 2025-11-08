package monero

import (
	"errors"
	"time"

	moneroapi "github.com/ObscuraProject/obscura-payment-gate/apis/monero"
	"github.com/ObscuraProject/obscura-payment-gate/db"
)

type MoneroWallet struct {
	PublicKey    string `json:"address" gorm:"primary_key"`
	AccountIndex uint   `json:"account_index" gorm:"index"`
	Type         string `json:"type"`

	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"index"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"index"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type MoneroWallets []MoneroWallet

type MoneroWalletBalance struct {
	ID              int     `json:"id" gorm:"primary_key"`
	PublicKey       string  `json:"address" gorm:"index"`
	Balance         float64 `json:"balance"`
	UnlockedBalance float64 `json:"unlocked_balance"`

	CreatedAt *time.Time `json:"created_at" gorm:"index"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"index"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	MoneroWallet MoneroWallet `json:"-" gorm:"ForeignKey:PublicKey;AssociationForeignKey:PublicKey"`
}

type DisplayMoneroWallet struct {
	PublicKey       string     `json:"address"`
	AccountIndex    uint       `json:"account_index"`
	Balance         float64    `json:"balance"`
	UnlockedBalance float64    `json:"unlocked_balance"`
	Type            string     `json:"type"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

func (w MoneroWalletBalance) DisplayMoneroWallet() DisplayMoneroWallet {
	ds := DisplayMoneroWallet{
		PublicKey:       w.PublicKey,
		UpdatedAt:       w.UpdatedAt,
		CreatedAt:       w.MoneroWallet.CreatedAt,
		Balance:         w.Balance,
		UnlockedBalance: w.UnlockedBalance,
		Type:            w.MoneroWallet.Type,
		AccountIndex:    w.MoneroWallet.AccountIndex,
	}

	return ds
}

/*
	Response Models
*/

type XMRPayment struct {
	Address string  `json:"address"`
	Percent float64 `json:"percent"`
	Amount  float64 `json:"amount"`
}

/*
	DB Models
*/

func CreateMoneroWallet(tp string) (MoneroWallet, error) {
	api := moneroapi.API

	numWallets := CountMoneroWallets()
	publicKey, addressIndex, err := api.CreateWallet(tp, numWallets+1)
	if err != nil {
		return MoneroWallet{}, err
	}

	now := time.Now()
	moneroWallet := MoneroWallet{
		PublicKey:    *publicKey,
		AccountIndex: *addressIndex,
		Type:         tp,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	err = moneroWallet.Save()
	return moneroWallet, err
}

func (w MoneroWallet) Validate() error {
	if w.PublicKey == "" {
		return errors.New("Public key is not valid")
	}

	return nil
}

func (w MoneroWallet) Save() error {
	err := w.Validate()
	if err != nil {
		return err
	}
	return w.SaveToDatabase()
}

/*
	External APIs
*/

func GetMoneroWalletBalance(w MoneroWallet) (*MoneroWalletBalance, error) {
	api := moneroapi.API
	unlockedBalance, balance, err := api.Balance(w.AccountIndex)
	if err != nil {
		return &MoneroWalletBalance{}, err
	}

	return &MoneroWalletBalance{
		PublicKey:       w.PublicKey,
		Balance:         float64(*balance) / 100e9,
		UnlockedBalance: float64(*unlockedBalance) / 100e9,
	}, nil
}

/*
	Balance Methods
*/

func (w MoneroWallet) UpdateBalance() (*MoneroWalletBalance, error) {
	wb, err := GetMoneroWalletBalance(w)
	if err != nil {
		return nil, err
	}

	currentBalance, _ := w.CurrentBalance()

	// First-time balance update
	if currentBalance == nil {
		return wb, wb.Save()
	} else if (currentBalance.Balance != wb.Balance) || (currentBalance.UnlockedBalance != wb.UnlockedBalance) {
		return wb, wb.Save()
	}

	return currentBalance, nil
}

func (w MoneroWallet) CurrentBalance() (*MoneroWalletBalance, error) {
	var wb MoneroWalletBalance

	err := db.Database.
		Where(&MoneroWalletBalance{PublicKey: w.PublicKey}).
		Order("ID desc").
		First(&wb).
		Error
	if err != nil {
		return nil, err
	}

	return &wb, err
}

/*
	Financial Methods
*/

func (w MoneroWallet) SendWithPercentSplit(payments []XMRPayment) ([]moneroapi.MoneroTransferResult, error) {
	balance, err := w.CurrentBalance()
	if err != nil {
		return nil, err
	}

	balanceSat := uint(balance.UnlockedBalance * 100e9)
	onePercent := uint(balanceSat / 100)

	paymentsResults := make([]moneroapi.MoneroTransferResult, len(payments))

	for i, payment := range payments {

		percentToSend := uint(payment.Percent * 100.0)
		amountToSend := uint(onePercent * percentToSend)

		result, err := moneroapi.API.Transfer(
			payment.Address,
			w.AccountIndex,
			amountToSend,
		)
		if err != nil {
			return paymentsResults, err
		}

		paymentsResults[i] = *result

	}

	return paymentsResults, nil
}

func (w MoneroWallet) SendWithAmountSplit(payments []XMRPayment) ([]moneroapi.MoneroTransferResult, error) {
	paymentsResults := make([]moneroapi.MoneroTransferResult, len(payments))

	for i, payment := range payments {

		amountToSend := uint(payment.Amount * 100e9)

		result, err := moneroapi.API.Transfer(
			payment.Address,
			w.AccountIndex,
			amountToSend,
		)
		if err != nil {
			return paymentsResults, err
		}

		paymentsResults[i] = *result

	}

	return paymentsResults, nil
}

func (w MoneroWallet) Sweep(addressTo string) error {
	return moneroapi.API.SweepAll(
		addressTo,
		w.AccountIndex,
	)
}

/*
	Database
*/

func (w MoneroWalletBalance) Save() error {
	return db.Database.Create(&w).Error
}

func (w MoneroWallet) SaveToDatabase() error {
	if existing, _ := FindMoneroWalletByPublicKey(w.PublicKey); existing == nil {
		return db.Database.Create(&w).Error
	}
	return db.Database.Save(&w).Error
}

func FindMoneroWalletByPublicKey(publicKey string) (*MoneroWallet, error) {
	var w MoneroWallet

	err := db.Database.
		Where(
			"public_key=?",
			publicKey,
		).
		First(&w).
		Error
	if err != nil {
		return nil, err
	}

	return &w, err
}

func CountMoneroWallets() uint {
	var count uint
	err := db.Database.
		Model(&MoneroWallet{}).
		Count(&count).
		Error
	if err != nil {
		return 0
	}

	return count
}

func FindMoneroWalletsWithNonZeroBalance() ([]MoneroWalletBalance, error) {
	moneroWalletBalances := []MoneroWalletBalance{}

	err := db.Database.
		Table("v_monero_wallet_balances").
		Model(MoneroWalletBalance{}).
		Preload("MoneroWallet").
		Where("balance > ? or unlocked_balance > ?", 0, 0).
		Find(&moneroWalletBalances).
		Error

	return moneroWalletBalances, err
}

func FindAllMoneroWallets() ([]MoneroWallet, error) {
	moneroWallets := []MoneroWallet{}

	err := db.Database.
		Order("created_at desc").
		Find(&moneroWallets).
		Error

	return moneroWallets, err
}

/*
	Database Views
*/

func SetupViews() {
	db.Database.Exec("DROP VIEW IF EXISTS v_monero_wallet_balances")
	db.Database.Exec(`
		CREATE VIEW v_monero_wallet_balances AS
			WITH balances AS (
				SELECT * 
				FROM (
					SELECT DISTINCT ON (public_key) *
					FROM monero_wallet_balances
					ORDER BY public_key, updated_at DESC, ID
				)  bcs
			)

			SELECT 
				balances.id, 
				balances.public_key, 
				balances.balance, 
				balances.unlocked_balance, 
				monero_wallets.account_index, 
				monero_wallets.type, 
				monero_wallets.created_at, 
				balances.updated_at, 
				balances.deleted_at FROM balances
			JOIN monero_wallets on monero_wallets.public_key=balances.public_key
			ORDER BY type, balance DESC
	`)
}
