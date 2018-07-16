package tools

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type MSP struct {
	MspId       int    `db:"msp_id" json:"msp_id"`
	Wallet  string `db:"wallet" json:"wallet"`
	Seed string `db:"seed" json:"seed"`
	Name       string `db:"name" json:"name"`
	Address1    string `db:"address_1" json:"address_1"`
	Address2    string `db:"address_2" json:"address_2"`
	Town    string `db:"town" json:"town"`
	Postcode    string `db:"postcode" json:"postcode"`
	MailAddr    string `db:"mail_address" json:"mail_addr"`
	Website    string `db:"website" json:"website"`
	VatNumber    string `db:"vat_number" json:"vat_number"`

}

type CPO struct {
	CpoId       int    `db:"cpo_id" json:"cpo_id"`
	Wallet  string `db:"wallet" json:"wallet"`
	Seed string `db:"seed" json:"seed"`
	Name       string `db:"name" json:"name"`
	Address1    string `db:"address_1" json:"address_1"`
	Address2    string `db:"address_2" json:"address_2"`
	Town    string `db:"town" json:"town"`
	Postcode    string `db:"postcode" json:"postcode"`
	MailAddr    string `db:"mail_address" json:"mail_addr"`
	Website    string `db:"website" json:"website"`
	VatNumber    string `db:"vat_number" json:"vat_number"`

}

var DB *sqlx.DB

func Connect(dbName string) {
	DB = sqlx.MustConnect("sqlite3", dbName)

	//some benchmark should be done here
	DB.SetMaxOpenConns(300)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(10 * time.Second)
}
