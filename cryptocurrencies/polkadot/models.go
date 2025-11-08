package polkadot

import (
	"errors"
	"time"

	polkadotapi "github.com/ObscuraProject/obscura-payment-gate/apis/polkadot"
	"github.com/ObscuraProject/obscura-payment-gate/db"
)

type PolkadotWallet struct {
	PublicKey string `json:"public_key" gorm:"primary_key"`
	Mnemonic  string `json:"mnemonic"`

	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"index"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"index"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type PolkadotWallets []PolkadotWallet

type PolkadotWalletBalance struct {
	ID        int    `json:"id" gorm:"primary_key"`
	PublicKey string `json:"address" gorm:"index"`

	FreeBalance     float64 `json:"free_balance"`
	ReservedBalance float64 `json:"reserved_balance"`
	FrozenBalance   float64 `json:"frozen_balance"`

	CreatedAt *time.Time `json:"created_at" gorm:"index"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"index"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	PolkadotWallet PolkadotWallet `json:"-" gorm:"ForeignKey:PublicKey;AssociationForeignKey:PublicKey"`
}

type DisplayPolkadotWallet struct {
	PublicKey       string     `json:"address"`
	Mnemonic        string     `json:"mnemonic"`
	FreeBalance     float64    `json:"free_balance"`
	ReservedBalance float64    `json:"reserved_balance"`
	FrozenBalance   float64    `json:"frozen_balance"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

func (w PolkadotWalletBalance) DisplayPolkadotWallet() DisplayPolkadotWallet {
	ds := DisplayPolkadotWallet{
		PublicKey:       w.PublicKey,
		UpdatedAt:       w.UpdatedAt,
		CreatedAt:       w.PolkadotWallet.CreatedAt,
		FreeBalance:     w.FreeBalance,
		ReservedBalance: w.ReservedBalance,
		FrozenBalance:   w.FrozenBalance,
		Mnemonic:        w.PolkadotWallet.Mnemonic,
	}
	return ds
}

/*
	Request / Response Models
*/

type DOTPayment struct {
	Address string  `json:"address"`
	Percent float64 `json:"percent"`
	Amount  float64 `json:"amount"`
}

type DOTPaymentResult struct {
	Hash       string  `json:"hash"`
	WalletFrom string  `json:"wallet_from"`
	WalletTo   string  `json:"wallet_to"`
	Amount     float64 `json:"amount"`
}

type DOTEscrowMintMetadata struct {
	SellerAddress string  `json:"seller_address"`
	Cid           string  `json:"cid"`
	Value         float64 `json:"value"`
}

type DOTCrowdloanMintMetadata struct {
	GoalValue      float64 `json:"goal_value"`
	GoalDate       uint64  `json:"goal_date"`
	WeeklyInterest uint16  `json:"weekly_interest"`
}

/*
	DB Models
*/

func CreatePolkadotWallet() (PolkadotWallet, error) {
	polkaWallet, err := polkadotapi.CreatePolkadotWallet()
	if err != nil {
		return PolkadotWallet{}, err
	}

	now := time.Now()
	polkadotWallet := PolkadotWallet{
		PublicKey: polkaWallet.PublicKey,
		Mnemonic:  polkaWallet.Mnemonic,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	err = polkadotWallet.Save()
	return polkadotWallet, err
}

func (w PolkadotWallet) Validate() error {
	if w.PublicKey == "" {
		return errors.New("Public key is not valid")
	}

	if w.Mnemonic == "" {
		return errors.New("Mnemonic is not valid")
	}

	return nil
}

func (w PolkadotWallet) Save() error {
	err := w.Validate()
	if err != nil {
		return err
	}
	return w.SaveToDatabase()
}

/*
	External APIs
*/

func GetPolkadotWalletBalance(w PolkadotWallet) (*PolkadotWalletBalance, error) {
	balance, err := polkadotapi.GetWalletBalance(w.PublicKey)
	if err != nil {
		return &PolkadotWalletBalance{}, err
	}

	return &PolkadotWalletBalance{
		PublicKey:       w.PublicKey,
		FreeBalance:     *&balance.Free,
		ReservedBalance: *&balance.Reserved,
		FrozenBalance:   *&balance.Frozen,
	}, nil
}

/*
	Balance Methods
*/

func (w PolkadotWallet) UpdateBalance() (*PolkadotWalletBalance, error) {
	wb, err := GetPolkadotWalletBalance(w)
	if err != nil {
		return nil, err
	}

	currentBalance, _ := w.CurrentBalance()

	// First-time balance update
	if currentBalance == nil {
		return wb, wb.Save()
	} else if (currentBalance.FreeBalance != wb.FreeBalance) ||
		(currentBalance.FrozenBalance != wb.FrozenBalance) ||
		(currentBalance.ReservedBalance != wb.ReservedBalance) {
		return wb, wb.Save()
	}

	return currentBalance, nil
}

func (w PolkadotWallet) CurrentBalance() (*PolkadotWalletBalance, error) {
	var wb PolkadotWalletBalance

	err := db.Database.
		Where(&PolkadotWalletBalance{PublicKey: w.PublicKey}).
		Order("ID desc").
		First(&wb).
		Error
	if err != nil {
		return nil, err
	}

	return &wb, err
}

/*
	Database
*/

func (w PolkadotWalletBalance) Save() error {
	return db.Database.Create(&w).Error
}

func (w PolkadotWallet) SaveToDatabase() error {
	if existing, _ := FindPolkadotWalletByPublicKey(w.PublicKey); existing == nil {
		return db.Database.Create(&w).Error
	}
	return db.Database.Save(&w).Error
}

func FindPolkadotWalletByPublicKey(publicKey string) (*PolkadotWallet, error) {
	var w PolkadotWallet

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

func CountPolkadotWallets() uint {
	var count uint
	err := db.Database.
		Model(&PolkadotWallet{}).
		Count(&count).
		Error
	if err != nil {
		return 0
	}

	return count
}

func FindPolkadotWalletsWithNonZeroBalance() ([]PolkadotWalletBalance, error) {
	PolkadotWalletBalances := []PolkadotWalletBalance{}

	err := db.Database.
		Table("v_polkadot_wallet_balances").
		Model(PolkadotWalletBalance{}).
		Preload("PolkadotWallet").
		Where("free_balance > ? or reserved_balance > ? or reserved_balance > ?", 0, 0, 0).
		Find(&PolkadotWalletBalances).
		Error

	return PolkadotWalletBalances, err
}

func FindAllPolkadotWallets() ([]PolkadotWallet, error) {
	PolkadotWallets := []PolkadotWallet{}

	err := db.Database.
		Order("created_at desc").
		Find(&PolkadotWallets).
		Error

	return PolkadotWallets, err
}

/*
	Database Views
*/

func SetupViews() {
	db.Database.Exec("DROP VIEW IF EXISTS v_polkadot_wallet_balances")
	db.Database.Exec(`
		CREATE VIEW v_polkadot_wallet_balances AS
			WITH balances AS (
				SELECT *
				FROM (
					SELECT DISTINCT ON (public_key) *
					FROM polkadot_wallet_balances
					ORDER BY public_key, updated_at DESC, ID
				)  bcs
			)

			SELECT
				balances.id,
				balances.public_key,
				balances.frozen_balance,
				balances.free_balance,
				balances.reserved_balance,
				polkadot_wallets.mnemonic,
				polkadot_wallets.created_at,
				balances.updated_at,
				polkadot_wallets.deleted_at FROM balances
			JOIN polkadot_wallets on polkadot_wallets.public_key=balances.public_key
			ORDER BY type, free_balance DESC, reserved_balance DESC, frozen_balance DESC
	`)
}
