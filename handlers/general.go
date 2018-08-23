package handlers

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"net/http"
	"strconv"
	"bytes"
	"encoding/csv"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Look! It's moving. It's alive. It's alive... It's alive, it's moving, it's alive, it's alive, it's alive, it's alive, IT'S ALIVE! (Frankenstein 1931)"})
}

//gets the balance for a wallet in Ether (EV Coin) (the thing that pays for the gas)
func GetWalletBalance(c *gin.Context) {

	addr := c.Param("addr")

	type TBalance struct {
		Balance string `json:"balance"`
	}

	body := tools.GETRequest("http://localhost:3000/api/wallet/balance/" + addr)

	var tBalance TBalance
	err := json.Unmarshal(body, &tBalance)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	balanceFloat, _ := strconv.ParseFloat(string(tBalance.Balance), 64)
	c.JSON(http.StatusOK, gin.H{"balance": balanceFloat / 1000000000000000000, "currency": "EV Coin"})
}

// get the history of transaction for ETH (EV Coin)
//TODO: refactor this not to use mysql, but the FATDB
func GetWalletHistoryEVCoin(c *gin.Context) {
	addr := c.Param("addr")

	type History struct {
		Block           int    `json:"block" db:"block"`
		FromAddr        string `json:"from_addr" db:"from_addr"`
		ToAddr          string `json:"to_addr" db:"to_addr"`
		Amount          uint64 `json:"amount" db:"amount"`
		Currency        string `json:"currency" db:"currency"`
		CreatedAt       uint64 `json:"created_at" db:"created_at"`
		TransactionHash string `json:"transaction_hash" db:"transaction_hash"`
	}
	var histories []History

	var transactions []tools.TxTransaction
	err := tools.MDB.Select(&transactions, "SELECT * FROM transactions WHERE (to_addr = ? OR from_addr = ?) ORDER BY blockNumber DESC", addr, addr)
	tools.ErrorCheck(err, "cpo.go", false)

	for _, tx := range transactions {
		if tx.Value == "0x0" {
			//we have a contract tx

			var txResponse tools.TxReceiptResponse
			err := tools.MDB.QueryRowx("SELECT * FROM transaction_receipts WHERE transactionHash = ?", tx.Hash).StructScan(&txResponse)
			tools.ErrorCheck(err, "cpo.go", false)
			calculatedGas := tools.HexToUInt(txResponse.GasUsed) * tools.HexToUInt(tx.GasPrice)
			//fake it
			calculatedGas = calculatedGas * 10
			log.Info("calculated gas is: %d" , calculatedGas)
			if calculatedGas > 1000000000 {
				histories = append(histories, History{Block: tx.BlockNumber, FromAddr: tx.From, ToAddr: tx.To, Amount: calculatedGas, Currency: "wei", CreatedAt: tx.Timestamp, TransactionHash: tx.Hash})
			}

		} else {
			//we have eth transfer
			histories = append(histories, History{Block: tx.BlockNumber, FromAddr: tx.From, ToAddr: tx.To, Amount: tools.HexToUInt(tx.Value), Currency: "wei", CreatedAt: tx.Timestamp, TransactionHash: tx.Hash})
		}
	}

	c.JSON(http.StatusOK, histories)
}

// getting the token info
func TokenInfo(c *gin.Context) {
	type TokenInfo struct {
		Name    string `json:"name"`
		Symbol  string `json:"symbol"`
		Address string `json:"address"`
		Owner   string `json:"owner"`
	}
	body := tools.GETRequest("http://localhost:3000/api/token/info")

	var tokenInfo = new(TokenInfo)
	//var tokenInfo TokenInfo  //TODO: check this one
	err := json.Unmarshal(body, &tokenInfo)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}
	c.JSON(http.StatusOK, tokenInfo)
}

// getting the token balance for an /:addr
func TokenBalance(c *gin.Context) {

	addr := c.Param("addr")
	log.Printf("getting token balance for %s", addr)

	body := tools.GETRequest("http://localhost:3000/api/token/balance/" + addr)

	log.Info("Balance for %s is %s", addr, body)
	balanceFloat, _ := strconv.ParseFloat(string(body), 64)
	c.JSON(http.StatusOK, gin.H{"balance": balanceFloat})

}

// mint the tokens for the EV Driver /:addr?amount=100
func TokenMint(c *gin.Context) {

	addr := c.Param("addr")
	amount := c.DefaultQuery("amount", "100")
	log.Printf("mint tokens for %s with the amount %s", addr, amount)

	amountFloat, _ := strconv.ParseFloat(string(amount), 64)
	if amountFloat == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "the amount is incorrect"})
		return
	}

	values := map[string]interface{}{"driver": addr, "amount": amountFloat}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = tools.POSTRequest("http://localhost:3000/api/token/mint", jsonValue)
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
`

	tools.DB.MustExec(schema)

	c.JSON(http.StatusOK, gin.H{"status": "database truncated."})
}



//shows the CDR records of a reimbursement
func ViewCDRs(c *gin.Context) {

	reimbursementId := c.Param("reimbursement_id")



	type Reimbursement struct {
		CdrRecords string `json:"cdr_records" db:"cdr_records"`
	}
	var reimbursement Reimbursement

	err := tools.MDB.QueryRowx("SELECT cdr_records FROM reimbursements WHERE reimbursement_id = ?", reimbursementId).StructScan(&reimbursement)
	tools.ErrorCheck(err, "cpo.go", false)

	if reimbursement.CdrRecords == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "no reimbursements found"})
		return
	}

	var cdrs []tools.CDR
	err = json.Unmarshal([]byte(reimbursement.CdrRecords), &cdrs)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}


	b := &bytes.Buffer{} // creates IO Writer
	wr := csv.NewWriter(b) // creates a csv writer that uses the io buffer.


	wr.Write([]string{"locationName", "locationAddress", "evseId", "scId","controller","start","end","finalPrice","tokenContract","tariff","chargedUnits","chargingContract","transactionHash","currency"})
	for _, cdr := range cdrs {
		wr.Write([]string{cdr.LocationName, cdr.LocationAddress, cdr.EvseID, cdr.ScID, cdr.Controller, cdr.Start, cdr.End, cdr.FinalPrice, cdr.TokenContract, cdr.Tariff, cdr.ChargedUnits, cdr.ChargingContract, cdr.TransactionHash, cdr.Currency})
	}
	wr.Flush()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=cdrs_"+reimbursementId+".csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())

}
