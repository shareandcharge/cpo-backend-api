package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"fmt"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
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
	}

	//check if there is already an msp registered
	rows, err := tools.DB.Query("SELECT msp_id FROM msp")
	tools.ErrorCheck(err, "msp.go", true)

	//check if we already have an MSP registered
	if rows.Next() {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "there's already an MSP registered on this backend"})
	}

	//if not, insert a new one with ID = 1, unique.
	query := "INSERT INTO msp (msp_id, wallet, seed, name, address_1, address_2, town, postcode, mail_address, website, vat_number) VALUES (%d, '%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')"
	command := fmt.Sprintf(query, 1, "", "", mspInfo.Name, mspInfo.Address1, mspInfo.Address2, mspInfo.Town, mspInfo.Postcode, mspInfo.MailAddress, mspInfo.Website, mspInfo.VatNumber)
	tools.DB.MustExec(command)

	c.JSON(http.StatusOK, gin.H{"status": "created ok"})
}

//returns the info for the MSP
func MspInfo(c *gin.Context) {


	rows, err := tools.DB.Query("SELECT msp_id FROM msp")
	tools.ErrorCheck(err, "msp.go", false)

	//check if we already have an MSP registered
	if rows.Next() {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "there's already an MSP registered on this backend"})
	}

	msp := tools.MSP{}

	err = tools.DB.QueryRowx("SELECT * FROM msp LIMIT 1").StructScan(&msp)
	tools.ErrorCheck(err,"msp.go", false)
	c.JSON(http.StatusOK, msp)
}

//
//// handling the wallet creation
//func CpoCreate(c *gin.Context) {
//
//	t := "INSERT INTO cpo (cpo_id, public_addr, seed, email, password) VALUES (%d, '%s','%s', '%s', '%s')"
//	command := fmt.Sprintf(t, 1, "0x123123123123", "word eye leg ...", "cpo@email.com", "hardpassword")
//	tools.DB.MustExec(command)
//
//	c.JSON(http.StatusOK, gin.H{"status": "wallet creation here."})
//}
//
//
//
//func CpoInfo(c *gin.Context) {
//
//	cpo := tools.CPO{}
//	err := tools.DB.QueryRowx("SELECT * FROM cpo").StructScan(&cpo)
//	if err != nil {
//		log.Panic(err)
//	}
//	c.JSON(http.StatusOK, cpo)
//}
//
//// getting the token info
//func TokenInfo(c *gin.Context) {
//	type TokenInfo struct {
//		Name    string `json:"name"`
//		Symbol  string `json:"symbol"`
//		Address string `json:"address"`
//		Owner   string `json:"owner"`
//	}
//	body := tools.GetRequest("http://localhost:3000/api/token/info")
//
//	var tokenInfo = new(TokenInfo)
//	err := json.Unmarshal(body, &tokenInfo)
//	if err != nil {
//		log.Panic(err)
//	}
//	c.JSON(http.StatusOK, tokenInfo)
//}
//
//// this will TRUNCATE the database.
//func Reinit(c *gin.Context) {
//
//	var schema = `
//		DROP TABLE IF EXISTS cpo;
//		CREATE TABLE cpo (
//			cpo_id    INTEGER PRIMARY KEY,
//    		public_addr VARCHAR(80)  DEFAULT '',
//    		seed  VARCHAR(250)  DEFAULT '',
//			email      VARCHAR(250) DEFAULT '',
//			password   VARCHAR(250) DEFAULT NULL
//		);
//`
//
//	tools.DB.MustExec(schema)
//
//	c.JSON(http.StatusOK, gin.H{"status": "database truncated."})
//}
//
////func HandleWriteBlock(c *gin.Context) {
////	var mess Message
////	if err := c.MustBindWith(&mess, binding.JSON); err == nil {
////
////		//validate your data here
////		if err := validateTheData(mess.TheData); err != nil {
////			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
////			return
////		}
////
////		newBlock, er := generateBlock(Blockchain[len(Blockchain)-1], mess.TheData)
////		if er != nil {
////			c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
////			return
////		}
////		if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
////			newBlockchain := append(Blockchain, newBlock)
////			replaceChain(newBlockchain)
////			spew.Dump(Blockchain)
////		}
////		c.JSON(http.StatusCreated, gin.H{"status": "block " + strconv.Itoa(newBlock.Index) + " added"})
////
////
////
////
////	} else {
////		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
////	}
////}
//
////func validateTheData(theData string) error {
////	if theData == "" {
////		return errors.New("invalid data")
////	}
////	return nil
////}
