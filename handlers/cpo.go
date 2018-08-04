package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"fmt"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"github.com/motionwerkGmbH/cpo-backend-api/configs"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"strings"
	"time"
	"strconv"
	"io/ioutil"
)

func CpoCreate(c *gin.Context) {

	type CpoInfo struct {
		Name        string `json:"name"`
		Address1    string `json:"address_1"`
		Address2    string `json:"address_2"`
		Town        string `json:"town"`
		Postcode    string `json:"postcode"`
		MailAddress string `json:"mail_address"`
		Website     string `json:"website"`
		VatNumber   string `json:"vat_number"`
	}
	var cpoInfo CpoInfo

	if err := c.MustBindWith(&cpoInfo, binding.JSON); err == nil {
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if there is already an CPO registered
	rows, err := tools.DB.Query("SELECT cpo_id FROM cpo")
	tools.ErrorCheck(err, "cpo.go", true)
	defer rows.Close()

	//check if we already have an CPO registered
	if rows.Next() {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "there's already an CPO registered on this backend"})
		return
	}

	//if not, insert a new one with ID = 1, unique.
	query := "INSERT INTO cpo (cpo_id, wallet, seed, name, address_1, address_2, town, postcode, mail_address, website, vat_number) VALUES (%d, '%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')"
	command := fmt.Sprintf(query, 1, "", "", cpoInfo.Name, cpoInfo.Address1, cpoInfo.Address2, cpoInfo.Town, cpoInfo.Postcode, cpoInfo.MailAddress, cpoInfo.Website, cpoInfo.VatNumber)
	tools.DB.MustExec(command)

	c.JSON(http.StatusOK, gin.H{"status": "created ok"})
}

//returns the info for the CPO
func CpoInfo(c *gin.Context) {

	rows, _ := tools.DB.Query("SELECT cpo_id FROM cpo")
	defer rows.Close()

	//check if we already have an CPO registered
	if rows.Next() == false {
		c.JSON(http.StatusNotFound, gin.H{"error": "we couldn't find any CPO registered in the database."})
		return
	}

	cpo := tools.CPO{}

	tools.DB.QueryRowx("SELECT * FROM cpo LIMIT 1").StructScan(&cpo)
	c.JSON(http.StatusOK, cpo)
}

// the main function for the wallets section of payment page
func CpoPaymentWallet(c *gin.Context) {

	config := configs.Load()
	cpoWallet := config.GetString("cpo.wallet_address")

	body := tools.GETRequest("http://localhost:3000/api/token/balance/" + cpoWallet)

	log.Printf("Balance is %s", body)
	balanceFloat, _ := strconv.ParseFloat(string(body), 64)

	type WalletRecord struct {
		MspName    string `json:"msp_name"`
		MspAddress string `json:"msp_address"`

		TotalTransactions int     `json:"total_transactions"`
		Amount            float64 `json:"amount"`
		Currency          string  `json:"currency"`
		TokenAddr         string  `json:"token_address"`
	}
	var walletRecords []WalletRecord

	//TODO: fix the total transactions count
	record := WalletRecord{MspName: "Charge & Fuel", MspAddress: "0xf60b71a4d360a42ec9d4e7977d8d9928fd7c8365", TotalTransactions: 1337, Amount: balanceFloat, Currency: "Charge & Fuel Token", TokenAddr: "0x682F10b5e35bA3157E644D9e7c7F3C107EB46305"}

	walletRecords = append(walletRecords, record)

	c.JSON(http.StatusOK, walletRecords)
}

//see a list of transactions from a particular msp
func CpoTransactionFromMsp(c *gin.Context) {

	type History struct {
		Id              int    `json:"id" db:"id"`
		Block           int    `json:"block" db:"block"`
		FromAddr        string `json:"from_addr" db:"from_addr"`
		ToAddr          string `json:"to_addr" db:"to_addr"`
		Amount          uint64 `json:"amount" db:"amount"`
		Currency        string `json:"currency" db:"currency"`
		GasUsed         uint64 `json:"gas_used" db:"gas_used"`
		GasPrice        uint64 `json:"gas_price" db:"gas_price"`
		CreatedAt       uint64 `json:"created_at" db:"created_at"`
		TransactionHash string `json:"transaction_hash" db:"transaction_hash"`
	}
	var histories []History

	config := configs.Load()
	cpoWallet := config.GetString("cpo.wallet_address")

	err := tools.MDB.Select(&histories, "SELECT * FROM ethtosql WHERE to_addr = ? AND currency = ? ORDER BY block DESC", cpoWallet, "Charge & Fuel Token")
	tools.ErrorCheck(err, "cpo.go", false)

	c.JSON(http.StatusOK, histories)
}

// creates a Reimbursement
func CpoCreateReimbursement(c *gin.Context) {

	//gets the MSP address from url
	mspAddress := c.Param("msp_address")

	config := configs.Load()
	cpoWallet := config.GetString("cpo.wallet_address")

	//-------- gets the History of the account

	type Reimbursement struct {
		Id              int    `json:"id" db:"id"`
		MspName         string `json:"msp_name" db:"msp_name"`
		CpoName         string `json:"cpo_name" db:"cpo_name"`
		Amount          int    `json:"amount" db:"amount"`
		Currency        string `json:"currency" db:"currency"`
		Timestamp       int    `json:"timestamp" db:"timestamp"`
		Status          string `json:"status" db:"status"`
		ReimbursementId string `json:"reimbursement_id" db:"reimbursement_id"`
		History         string `json:"history" db:"history"`
	}

	// get the history

	type CDR struct {
		EvseID           string `json:"evseId"`
		ScID             string `json:"scId"`
		Controller       string `json:"controller"`
		Start            string `json:"start"`
		End              string `json:"end"`
		FinalPrice       string `json:"finalPrice"`
		TokenContract    string `json:"tokenContract"`
		ChargingContract string `json:"chargingContract"`
		TransactionHash  string `json:"transactionHash"`
		Currency         string `json:"currency"`
	}

	body := tools.GETRequest("http://localhost:3000/api/cdr/info") //+ ?tokenContract= tokenAddress

	var cdrs []CDR
	err := json.Unmarshal(body, &cdrs)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	drivers, err := tools.ReturnAllDrivers()
	tools.ErrorCheck(err, "cpo.go", false)

	var cdrsOutput []CDR

	for _, cdr := range cdrs {

		//map driver email to address
		for _, driver := range drivers {
			cdr.Currency = "Charge & Fuel Token"
			if driver.Address == cdr.Controller {
				cdr.Controller = driver.Email
				break
			}
		}

		//is this record already present in some reimbursement ?
		count := 0
		row := tools.MDB.QueryRow("SELECT COUNT(*) as count FROM reimbursements WHERE cdr_records LIKE '%" + cdr.TransactionHash + "%'")
		row.Scan(&count)

		if count == 0 {
			log.Info("we have an unprocessed transaction: " + cdr.TransactionHash)
			cdrsOutput = append(cdrsOutput, cdr)
		} else {
			log.Warn("transaction with hash " + cdr.TransactionHash + " already present in a reimbursement")
		}

	}

	if len(cdrsOutput) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "there are no more records to be used to create new reimbursement."})
		return
	}

	cdrsOutputBytes, err := json.Marshal(cdrsOutput)
	tools.ErrorCheck(err, "cpo.go", false)

	//calculate the amount of tokens CPO has
	var totalAmount uint64
	for _, h := range cdrsOutput {
		if h.Currency == "Charge & Fuel Token" {
			finalPriceInt, err := strconv.Atoi(h.FinalPrice)
			tools.ErrorCheck(err, "cpo.go", false)
			totalAmount = totalAmount + uint64(finalPriceInt)
		}
	}

	query := "INSERT INTO reimbursements ( msp_name, cpo_name, amount, currency, status, reimbursement_id, timestamp, cdr_records)" +
		"  VALUES ('%s','%s',%d,'%s','%s','%s',%d,'%s')"
	command := fmt.Sprintf(query, mspAddress, cpoWallet, totalAmount, "Charge&Fuel Token", "pending", tools.GetSha1Hash(cdrsOutput), time.Time.Unix(time.Now()), string(cdrsOutputBytes))
	_, err = tools.MDB.Exec(command)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "there's already a reimbursement issued for the current transactions."})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "reimbursement sent"})

}

// Lists all reimbursements
func CpoGetAllReimbursements(c *gin.Context) {

	status := c.Param("status")

	config := configs.Load()
	cpoWallet := config.GetString("cpo.wallet_address")

	type Reimbursement struct {
		Id              int    `json:"id" db:"id"`
		MspName         string `json:"msp_name" db:"msp_name"`
		CpoName         string `json:"cpo_name" db:"cpo_name"`
		Amount          int    `json:"amount" db:"amount"`
		Currency        string `json:"currency" db:"currency"`
		Timestamp       int    `json:"timestamp" db:"timestamp"`
		Status          string `json:"status" db:"status"`
		ReimbursementId string `json:"reimbursement_id" db:"reimbursement_id"`
		CdrRecords      string `json:"cdr_records" db:"cdr_records"`
	}
	var reimb []Reimbursement

	err := tools.MDB.Select(&reimb, "SELECT * FROM reimbursements WHERE cpo_name = ? AND status = ?", cpoWallet, status)
	tools.ErrorCheck(err, "cpo.go", false)

	if len(reimb) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	c.JSON(http.StatusOK, reimb)
}

// marks the reimbursement as complete
func CpoSetReimbursementComplete(c *gin.Context) {

	reimbursementId := c.Param("reimbursement_id")

	rows, err := tools.MDB.Query("SELECT id FROM reimbursements WHERE reimbursement_id = ?", reimbursementId)
	tools.ErrorCheck(err, "cpo.go", false)
	defer rows.Close()

	//check if we have a reimbursement with this id present
	if !rows.Next() {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "there's isn't any reimbursement with this id present"})
		return
	}

	query := "UPDATE reimbursements SET status='%s' WHERE reimbursement_id = '%s'"
	command := fmt.Sprintf(query, "complete", reimbursementId)
	_, err = tools.MDB.Exec(command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "complete"})

}

// the records for the particular token
func CpoPaymentCDR(c *gin.Context) {

	//tokenAddress := c.Param("token")

	type CDR struct {
		EvseID           string `json:"evseId"`
		ScID             string `json:"scId"`
		Controller       string `json:"controller"`
		Start            string `json:"start"`
		End              string `json:"end"`
		FinalPrice       string `json:"finalPrice"`
		TokenContract    string `json:"tokenContract"`
		ChargingContract string `json:"chargingContract"`
		TransactionHash  string `json:"transactionHash"`
		Currency         string `json:"currency"`
	}

	body := tools.GETRequest("http://localhost:3000/api/cdr/info") //+ ?tokenContract= tokenAddress

	var cdrs []CDR
	err := json.Unmarshal(body, &cdrs)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	drivers, err := tools.ReturnAllDrivers()
	tools.ErrorCheck(err, "cpo.go", false)

	var cdrsOutput []CDR

	for _, cdr := range cdrs {

		//map driver email to address
		for _, driver := range drivers {
			cdr.Currency = "Charge & Fuel Token"
			if driver.Address == cdr.Controller {
				cdr.Controller = driver.Email
				break
			}
		}

		count := 0
		row := tools.MDB.QueryRow("SELECT COUNT(*) as count FROM reimbursements WHERE cdr_records LIKE '%" + cdr.TransactionHash + "%'")
		row.Scan(&count)

		if count == 0 {
			log.Info("we have an unprocessed transaction hash " + cdr.TransactionHash)
			cdrsOutput = append(cdrsOutput, cdr)
		} else {
			log.Warn("transaction with hash " + cdr.TransactionHash + " already present in some reimbursement")
		}

	}

	if len(cdrsOutput) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	c.JSON(http.StatusOK, cdrsOutput)

}

//generates a new wallet for the cpo
func CpoGenerateWallet(c *gin.Context) {

	type WalletInfo struct {
		Seed string `json:"seed"`
		Addr string `json:"address"`
	}
	var walletInfo WalletInfo

	//Leave this commented code please
	//body := tools.GETRequest("http://localhost:3000/api/wallet/create")
	//log.Printf("<- %s", string(body))
	//err := json.Unmarshal(body, &walletInfo)
	//if err != nil {
	//	log.Panic(err)
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
	//	return
	//}

	config := configs.Load()
	walletInfo.Addr = config.GetString("cpo.wallet_address")
	walletInfo.Seed = config.GetString("cpo.wallet_seed")

	//update the db for CPO

	query := "UPDATE cpo SET wallet='%s', seed='%s' WHERE cpo_id = 1"
	command := fmt.Sprintf(query, walletInfo.Addr, walletInfo.Seed)
	tools.DB.MustExec(command)

	//update the ~/.sharecharge/config.json
	configs.UpdateBaseAccountSeedInSCConfig(walletInfo.Seed)

	c.JSON(http.StatusOK, walletInfo)
}

//returns the info for the CPO
func CpoGetSeed(c *gin.Context) {

	cpo := tools.CPO{}
	tools.DB.QueryRowx("SELECT * FROM cpo LIMIT 1").StructScan(&cpo)

	if cpo.Seed == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "there isn't any seed in the cpo account. Maybe you need to create the wallet first ?."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"seed": cpo.Seed})
}

//Gets the history for the CPO
func CpoHistory(c *gin.Context) {

	type History struct {
		Id              int    `json:"id" db:"id"`
		Block           int    `json:"block" db:"block"`
		FromAddr        string `json:"from_addr" db:"from_addr"`
		ToAddr          string `json:"to_addr" db:"to_addr"`
		Amount          uint64 `json:"amount" db:"amount"`
		Currency        string `json:"currency" db:"currency"`
		GasUsed         uint64 `json:"gas_used" db:"gas_used"`
		GasPrice        uint64 `json:"gas_price" db:"gas_price"`
		CreatedAt       uint64 `json:"created_at" db:"created_at"`
		TransactionHash string `json:"transaction_hash" db:"transaction_hash"`
	}
	var histories []History

	config := configs.Load()
	cpoWallet := config.GetString("cpo.wallet_address")

	err := tools.MDB.Select(&histories, "SELECT * FROM ethtosql WHERE to_addr = ? ORDER BY block DESC", cpoWallet)
	tools.ErrorCheck(err, "cpo.go", false)

	c.JSON(http.StatusOK, histories)
}

//=================================
//========= PDF Generation ========
//=================================

func CpoReimbursementGenPdf(c *gin.Context) {
	reimbursement_id := c.Param("reimbursement_id")

	type Reimbursement struct {
		Id              int    `json:"id" db:"id"`
		MspName         string `json:"msp_name" db:"msp_name"`
		CpoName         string `json:"cpo_name" db:"cpo_name"`
		Amount          int    `json:"amount" db:"amount"`
		Currency        string `json:"currency" db:"currency"`
		Timestamp       int    `json:"timestamp" db:"timestamp"`
		Status          string `json:"status" db:"status"`
		ReimbursementId string `json:"reimbursement_id" db:"reimbursement_id"`
		CdrRecords      string `json:"cdr_records" db:"cdr_records"`
	}
	var reimb Reimbursement

	err := tools.MDB.QueryRowx("SELECT * FROM reimbursements WHERE reimbursement_id = ? LIMIT 1", reimbursement_id).StructScan(&reimb)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	//load the template file
	b, err := ioutil.ReadFile("configs/invoice_template.html")
	tools.ErrorCheck(err, "cpo.go", false)

	htmlTemplateRaw := string(b)

	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromName}}", "ThePhoenixWorks", 2)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromAddress}}", "59-62R Springfield Centre LS28 5LY", 2)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceDate}}", "19 July 2018 ", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceNumber}}", "000000001", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{clientReference}}", " S&C000001 ", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{purchaseOrder}}", " S&C000001 ", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToName}}", "Volkswagen Financial", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToAddress}}", "Services (UKJ) Limited", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToPerson}}", "Milton Keynes", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToCode}}", "MK15 8HG", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceDueDate}}", "02 August 2018 ", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{description}}", "Sum of Tokens received through Share&Charge network ", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{quantity}}", strconv.Itoa(reimb.Amount), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{unit}}", "Tokens", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{price}}", "£ 0,01", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{vat}}", "20%", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{total}}", strconv.Itoa(reimb.Amount), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{subTotal}}", "£ "+strconv.Itoa(reimb.Amount), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{totalVat}}", "£ "+strconv.Itoa(int(float64(reimb.Amount)*float64(0.20))), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{totalAmount}}", "£ "+strconv.Itoa(int(float64(reimb.Amount)*float64(1.20))), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromPhone}}", "0014 882 739 2282", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromMail}}", "accounting@invoice.com", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromWebsite}}", "http://yourwebiste.com", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromBankName}}", "UNICREDIT", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromSortCode}}", "12312", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromAccountNumber}}", "123 345 532", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{vatNumber}}", "321ADF23", 1)

	//write it to a file
	ioutil.WriteFile("static/invoice_1.html", []byte(htmlTemplateRaw), 0644)

	//convert it to pdf
	err = tools.GeneratePdf("static/invoice_1.html", "static/invoice_1.pdf")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redirect": "http://{{server_addr}}:{{server_port}}/static/invoice_1.pdf"})
}

//=================================
//=========== LOCATIONS ===========
//=================================

//gets all locations of this CPO
func CpoGetLocations(c *gin.Context) {

	config := configs.Load()
	cpoAddress := config.GetString("cpo.wallet_address")
	body := tools.GETRequest("http://localhost:3000/api/store/locations/" + cpoAddress)

	var locations []tools.XLocation
	err := json.Unmarshal(body, &locations)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	if len(locations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "there aren't any locations registered with this CPO"})
		return
	}

	c.JSON(http.StatusOK, locations)

}

//uploads new locations and re-writes if they already are present
func CpoPutLocation(c *gin.Context) {
	var stations []tools.XLocation

	if err := c.MustBindWith(&stations, binding.JSON); err == nil {
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonValue, err := json.Marshal(stations)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = tools.PUTRequest("http://localhost:3000/api/store/locations", jsonValue)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

//uploads new location
func CpoPostLocation(c *gin.Context) {
	var stations []tools.Location

	if err := c.MustBindWith(&stations, binding.JSON); err == nil {
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonValue, err := json.Marshal(stations)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = tools.POSTRequest("http://localhost:3000/api/store/locations", jsonValue)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

//deletes a location
func CpoDeleteLocation(c *gin.Context) {

	locationid := c.Param("locationid")

	_, err := tools.DELETERequest("http://localhost:3000/api/store/locations/" + locationid)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})

}

//uploads new location
func CpoPostEvse(c *gin.Context) {
	var evse tools.Evse

	if err := c.MustBindWith(&evse, binding.JSON); err == nil {
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, evse)
}
