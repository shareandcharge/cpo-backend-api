package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"math/rand"
	"strconv"
	"math"
)

func AccountInfo(c *gin.Context) {

	type AccountInfo struct {
		Name        string   `json:"name"`
		Address1 string `json:"address_1"`
		Address2 string `json:"address_2"`
		Town string `json:"town"`
		Postcode string `json:"postcode"`
		MailAddress string `json:"mail_address"`
		Website string `json:"website"`
		VatNumber string `json:"vat_number"`
	}
	var accountInfo AccountInfo
	accountInfo.Name = "Ecotricity Group Limited"
	accountInfo.Address1 = "Lion House"
	accountInfo.Address2 = "2 Rowcroft"
	accountInfo.Town = "Stroud"
	accountInfo.Postcode = "GL5 3BY"
	accountInfo.MailAddress = "info@ecotricity.co.uk"
	accountInfo.Website = "https://www.ecotricity.co.uk"
	accountInfo.VatNumber = "1234567890"

	c.JSON(http.StatusOK,  accountInfo)
}


func WalletInfo(c *gin.Context) {

	type WalletInfo struct {
		Addr        string   `json:"addr"`
		Balance string `json:"balance"`
	}
	var walletInfo WalletInfo
	walletInfo.Addr = "0x003bd62e1aa057884CeAdB6eb773CaEaF87F5EBc"
	walletInfo.Balance = "1024 ETH"

	c.JSON(http.StatusOK,  walletInfo)
}

func AccountHistory(c *gin.Context) {


	type History struct {
		Amount float64      `json:"amount"`
		Currency string `json:"currency"`
		Timestamp string `json:"timestamp"`
	}

	s1 := rand.NewSource(1337)
	r1 := rand.New(s1)

	var histories []History
	for i := 0; i<100 ;i++ {
		n := History{Amount:  math.Floor(r1.Float64() * 10000) / 10000, Currency: "MSP Tokens", Timestamp:  "01.04.2018 "+strconv.Itoa(10+r1.Intn(23))+":"+strconv.Itoa(10+r1.Intn(49))+":" + strconv.Itoa(10+r1.Intn(49))}
		histories = append(histories,n)
	}



	c.JSON(http.StatusOK, histories)
}

func AccountMnemonic(c *gin.Context){
	type Mnemonic struct {
		Seed        string   `json:"seed"`
	}
	var mnemonic Mnemonic
	mnemonic.Seed = "health salt town tiger vintage trend cart nation grace mechanic long dial"

	c.JSON(http.StatusOK,  mnemonic)
}