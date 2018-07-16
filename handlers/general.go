package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"fmt"
	"log"
	"encoding/json"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Look! It's moving. It's alive. It's alive... It's alive, it's moving, it's alive, it's alive, it's alive, it's alive, IT'S ALIVE! (Frankenstein 1931)"})
}

// handling the wallet creation
func CpoCreate(c *gin.Context) {

	t := "INSERT INTO cpo (cpo_id, public_addr, seed, email, password) VALUES (%d, '%s','%s', '%s', '%s')"
	command := fmt.Sprintf(t, 1, "0x123123123123", "word eye leg ...", "cpo@email.com", "hardpassword")
	tools.DB.MustExec(command)

	c.JSON(http.StatusOK, gin.H{"status": "wallet creation here."})
}



func CpoInfo(c *gin.Context) {

	cpo := tools.CPO{}
	err := tools.DB.QueryRowx("SELECT * FROM cpo").StructScan(&cpo)
	if err != nil {
		log.Panic(err)
	}
	c.JSON(http.StatusOK, cpo)
}

// getting the token info
func TokenInfo(c *gin.Context) {
	type TokenInfo struct {
		Name    string `json:"name"`
		Symbol  string `json:"symbol"`
		Address string `json:"address"`
		Owner   string `json:"owner"`
	}
	body := tools.GetRequest("http://localhost:3000/api/token/info")

	var tokenInfo = new(TokenInfo)
	err := json.Unmarshal(body, &tokenInfo)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}
	c.JSON(http.StatusOK, tokenInfo)
}

// this will TRUNCATE the database.
func Reinit(c *gin.Context) {

	var schema = `
		DROP TABLE IF EXISTS cpo;
		CREATE TABLE cpo (
			cpo_id    INTEGER PRIMARY KEY,
    		wallet VARCHAR(80)  DEFAULT '',
    		seed  VARCHAR(250)  DEFAULT '',
			name      VARCHAR(250) DEFAULT '',
			address_1      VARCHAR(250) DEFAULT '',
			address_2      VARCHAR(250) DEFAULT '',
			town      VARCHAR(250) DEFAULT '',
			postcode      VARCHAR(250) DEFAULT '',
			mail_address      VARCHAR(250) DEFAULT '',
			website      VARCHAR(250) DEFAULT '',
			vat_number      VARCHAR(250) DEFAULT ''
		);
	DROP TABLE IF EXISTS msp;
	CREATE TABLE msp (
			msp_id    INTEGER PRIMARY KEY,
    		wallet VARCHAR(80)  DEFAULT '',
    		seed  VARCHAR(250)  DEFAULT '',
			name      VARCHAR(250) DEFAULT '',
			address_1      VARCHAR(250) DEFAULT '',
			address_2      VARCHAR(250) DEFAULT '',
			town      VARCHAR(250) DEFAULT '',
			postcode      VARCHAR(250) DEFAULT '',
			mail_address      VARCHAR(250) DEFAULT '',
			website      VARCHAR(250) DEFAULT '',
			vat_number      VARCHAR(250) DEFAULT ''
		);
`

	tools.DB.MustExec(schema)

	c.JSON(http.StatusOK, gin.H{"status": "database truncated."})
}

//func HandleWriteBlock(c *gin.Context) {
//	var mess Message
//	if err := c.MustBindWith(&mess, binding.JSON); err == nil {
//
//		//validate your data here
//		if err := validateTheData(mess.TheData); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		newBlock, er := generateBlock(Blockchain[len(Blockchain)-1], mess.TheData)
//		if er != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
//			return
//		}
//		if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
//			newBlockchain := append(Blockchain, newBlock)
//			replaceChain(newBlockchain)
//			spew.Dump(Blockchain)
//		}
//		c.JSON(http.StatusCreated, gin.H{"status": "block " + strconv.Itoa(newBlock.Index) + " added"})
//
//
//
//
//	} else {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//	}
//}

//func validateTheData(theData string) error {
//	if theData == "" {
//		return errors.New("invalid data")
//	}
//	return nil
//}
