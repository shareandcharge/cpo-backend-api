package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"fmt"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"github.com/motionwerkGmbH/cpo-backend-api/configs"
	"encoding/json"
	"log"
)

func MspCreate(c *gin.Context) {

	type MspInfo struct {
		Name        string `json:"name"`
		Address1    string `json:"address_1"`
		Address2    string `json:"address_2"`
		Town        string `json:"town"`
		Postcode    string `json:"postcode"`
		MailAddress string `json:"mail_address"`
		Website     string `json:"website"`
		VatNumber   string `json:"vat_number"`
	}
	var mspInfo MspInfo

	if err := c.MustBindWith(&mspInfo, binding.JSON); err == nil {
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if there is already an msp registered
	rows, err := tools.DB.Query("SELECT msp_id FROM msp")
	tools.ErrorCheck(err, "msp.go", true)

	//check if we already have an MSP registered
	if rows.Next() {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "there's already an MSP registered on this backend"})
		return
	}

	//if not, insert a new one with ID = 1, unique.
	query := "INSERT INTO msp (msp_id, wallet, seed, name, address_1, address_2, town, postcode, mail_address, website, vat_number) VALUES (%d, '%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')"
	command := fmt.Sprintf(query, 1, "", "", mspInfo.Name, mspInfo.Address1, mspInfo.Address2, mspInfo.Town, mspInfo.Postcode, mspInfo.MailAddress, mspInfo.Website, mspInfo.VatNumber)
	tools.DB.MustExec(command)

	c.JSON(http.StatusOK, gin.H{"status": "created ok"})
}

//returns the info for the MSP
func MspInfo(c *gin.Context) {


	rows, _ := tools.DB.Query("SELECT msp_id FROM msp")

	//check if we already have an MSP registered
	if rows.Next() == false {
		c.JSON(http.StatusNotFound, gin.H{"error": "we couldn't find any MPS registered in the database."})
		return
	}

	msp := tools.MSP{}

	tools.DB.QueryRowx("SELECT * FROM msp LIMIT 1").StructScan(&msp)
	c.JSON(http.StatusOK, msp)
}


//returns the info for the MSP
func MspGetSeed(c *gin.Context) {


	msp := tools.MSP{}
	tools.DB.QueryRowx("SELECT * FROM msp LIMIT 1").StructScan(&msp)

	if msp.Seed == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "there isn't any seed in the msp account. Maybe you need to create the wallet first ?."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"seed": msp.Seed})
}

//generates a new wallet for the msp
func MspGenerateWallet(c *gin.Context){

	type WalletInfo struct {
		Seed   string `json:"seed"`
		PubKey string `json:"pubKey"`
	}

	body := tools.GetRequest("http://localhost:3000/api/wallet/create")
	log.Printf("<- %s", string(body))

	var walletInfo WalletInfo
	err := json.Unmarshal(body, &walletInfo)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	//update the db for MSP
	query := "UPDATE msp SET wallet='%s', seed='%s' WHERE msp_id = 1"
	command := fmt.Sprintf(query, walletInfo.PubKey, walletInfo.Seed)
	tools.DB.MustExec(command)

	//update the ~/.sharecharge/config.json
	configs.UpdateBaseAccountSeedInSCConfig(walletInfo.Seed)

	c.JSON(http.StatusOK, walletInfo)
}