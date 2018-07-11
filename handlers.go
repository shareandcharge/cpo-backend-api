package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"fmt"
	"log"
)

func HandleIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Look! It's moving. It's alive. It's alive... It's alive, it's moving, it's alive, it's alive, it's alive, it's alive, IT'S ALIVE! (Frankenstein 1931)"})
}

// handling the wallet creation
func HandleCpoCreate(c *gin.Context) {

	t := "INSERT INTO cpo (cpo_id, public_addr, seed, email, password) VALUES (%d, '%s','%s', '%s', '%s')"
	command := fmt.Sprintf(t, 1, "0x123123123123", "word eye leg ...", "cpo@email.com", "hardpassword")
	tools.DB.MustExec(command)

	c.JSON(http.StatusOK, gin.H{"status": "wallet creation here."})
}

// handling the wallet info
func HandleWalletInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "wallet info here."})
}

func HandleCpoInfo(c *gin.Context) {

	cpo := tools.CPO{}
	err := tools.DB.QueryRowx("SELECT * FROM cpo").StructScan(&cpo)
	if err != nil {
		log.Panic(err)
	}
	c.JSON(http.StatusOK, cpo)
}

// this will TRUNCATE the database.
func HandleReinit(c *gin.Context) {

	var schema = `
		DROP TABLE IF EXISTS cpo;
		CREATE TABLE cpo (
			cpo_id    INTEGER PRIMARY KEY,
    		public_addr VARCHAR(80)  DEFAULT '',
    		seed  VARCHAR(250)  DEFAULT '',
			email      VARCHAR(250) DEFAULT '',
			password   VARCHAR(250) DEFAULT NULL
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
