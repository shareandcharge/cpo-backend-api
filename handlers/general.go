package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"strconv"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Look! It's moving. It's alive. It's alive... It's alive, it's moving, it's alive, it's alive, it's alive, it's alive, IT'S ALIVE! (Frankenstein 1931)"})
}

//gets the balance for a wallet in Ether (the thing that pays for the gas)
func GetWalletBalance(c *gin.Context) {

	addr := c.Param("addr")

	type TBalance struct {
		Balance string `json:"balance"`
	}

	body := tools.GetRequest("http://localhost:3000/api/wallet/balance/" + addr)

	var tBalance = new(TBalance)
	err := json.Unmarshal(body, &tBalance)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	log.Printf("Balance is %s", tBalance.Balance)
	balanceFloat, _ := strconv.ParseFloat(string(tBalance.Balance), 64)


	c.JSON(http.StatusOK, gin.H{"balance": balanceFloat / 1000000000000000000, "currency": "ETH"})
}


//Returns a list of all drivers
func GetAllDrivers(c *gin.Context){

	driversList, err := tools.ReturnAllDrivers()
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	var mDriversList []tools.Driver
	for _, driver := range driversList   {
		driver.Token = "S&C Token" //TODO: attention, it's hardcoded

		body := tools.GetRequest("http://localhost:3000/api/token/balance/" + driver.Address)
		balanceFloat, _ := strconv.ParseFloat(string(body), 64)
		driver.Balance = balanceFloat

		mDriversList = append(mDriversList, driver)

	}

	c.JSON(http.StatusOK, mDriversList)
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

// getting the token info
func TokenBalance(c *gin.Context) {

	addr := c.Param("addr")
	log.Printf("getting token balance for %s", addr)

	body := tools.GetRequest("http://localhost:3000/api/token/balance/" + addr)

	log.Printf("Balance is %s", body)
	balanceFloat, _ := strconv.ParseFloat(string(body), 64)
	c.JSON(http.StatusOK, gin.H{"balance": balanceFloat})

}


// mint the tokens for the EV Driver
func TokenMint(c *gin.Context) {

	addr := c.Param("addr")
	amount := c.DefaultQuery("amount", "100")
	log.Printf("mint tokens for %s with the amount %s", addr, amount)

	amountFloat, _ := strconv.ParseFloat(string(amount), 64)

	if amountFloat == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "the amount doesn't make sense"})
		return
	}

	values := map[string]interface{}{"driver": addr, "amount": amountFloat}

	_, err := tools.PostJsonRequest("http://localhost:3000/api/token/mint", values)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
