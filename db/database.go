package db

import (
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/ObscuraProject/obscura-payment-gate/settings"
)

var (
	Database  *gorm.DB
	TDatabase *gorm.DB
)

/*
	Entry point
*/

func init() {
	var err error
	Database, err = gorm.Open("postgres", settings.PAYMENT_GATE_SETTINGS.PostgresConnectionString)
	if err != nil {
		panic(err)
	}

	if settings.APPLICATION_SETTINGS.Debug {
		Database.LogMode(true)
	}

	Database.DB().SetMaxOpenConns(8)
}
