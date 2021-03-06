package handlers

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/motionwerkGmbH/cpo-backend-api/configs"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"net/http"
	"strconv"
	"strings"
	"time"
	"math/rand"
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

	config := configs.Load()
	addr := config.GetString("cpo.wallet_address")
	seed := config.GetString("cpo.wallet_seed")

	//if not, insert a new one with ID = 1, unique.
	query := "INSERT INTO cpo (cpo_id, wallet, seed, name, address_1, address_2, town, postcode, mail_address, website, vat_number) VALUES (%d, '%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')"
	command := fmt.Sprintf(query, 1, addr, seed, cpoInfo.Name, cpoInfo.Address1, cpoInfo.Address2, cpoInfo.Town, cpoInfo.Postcode, cpoInfo.MailAddress, cpoInfo.Website, cpoInfo.VatNumber)
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

	//get the total amount of transactions

	//-------- BEGIN HIStORY -------------
	log.Info("loading all cpo's locations, this might take some time...")
	locationBody := tools.GETRequest("http://localhost:3000/api/store/locations/" + cpoWallet)
	var locations []tools.XLocation
	err0 := json.Unmarshal(locationBody, &locations)
	if err0 != nil {
		log.Error(err0)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	var scIds []string

	for _, location := range locations {
		scIds = append(scIds, location.ScID)
	}

	log.Info("loading all cdrs, this might take some time...")
	body = tools.GETRequest("http://localhost:3000/api/cdr/info") //+ ?tokenContract= tokenAddress

	var cdrs []tools.CDR
	err := json.Unmarshal(body, &cdrs)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	var cdrsOutput []tools.CDR

	if len(cdrs) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	for _, cdr := range cdrs {

		cdr.Currency = "Charge & Fuel Token"

		var count int
		err = tools.MDB.QueryRowx("SELECT COUNT(*) as count FROM reimbursements WHERE cdr_records LIKE '%" + cdr.TransactionHash + "%'").Scan(&count)
		tools.ErrorCheck(err, "cpo.go", false)

		if count == 0 {
			log.Info("seems we have an unprocessed tx.")
			//todo: this should be removed once filtering is fixed
			if cdr.TokenContract == "0xAcD218713094a5F78Ea6d8D439DA22B5FCdb1A28" {
				log.Info("seems we have an unprocessed tx. that maches our token contract")
				isMyLocation := false
				for _, loc := range scIds {
					isMyLocation = loc == cdr.ScID
					if isMyLocation {
						//get the location name & address
						body = tools.GETRequest("http://localhost:3000/api/store/locations/" + cpoWallet + "/" + cdr.ScID)
						if body != nil {

							var loc tools.Location
							err := json.Unmarshal(body, &loc)
							if err != nil {
								log.Warnf(err.Error())
							} else {
								log.Info(loc)
								cdr.LocationName = loc.Name
								cdr.LocationAddress = loc.City + ", " + loc.Address + ", " + loc.Country
							}
						}

						cU, _ := strconv.ParseFloat(cdr.ChargedUnits, 64)
						cdr.ChargedUnits = fmt.Sprintf("%.3f", cU / 1000)
						cdrsOutput = append(cdrsOutput, cdr)
					}
				}
			}
		} else {
			log.Warn("transaction with hash " + cdr.TransactionHash + " already present in some reimbursement")
		}
	}

	//------------ END HISTORY

	sumTx := len(cdrsOutput)
	record := WalletRecord{MspName: "Charge & Fuel", MspAddress: "0xAcD218713094a5F78Ea6d8D439DA22B5FCdb1A28", TotalTransactions: sumTx, Amount: balanceFloat, Currency: "Charge & Fuel Token", TokenAddr: "0xAcD218713094a5F78Ea6d8D439DA22B5FCdb1A28"}

	walletRecords = append(walletRecords, record)

	c.JSON(http.StatusOK, walletRecords)
}

// creates a Reimbursement
func CpoCreateReimbursement(c *gin.Context) {

	//gets the MSP address from url
	mspAddress := c.Param("msp_address")
	tokenContract := c.DefaultQuery("tokenContract", "")

	config := configs.Load()
	cpoWallet := config.GetString("cpo.wallet_address")

	//get all ScIDs belonging to the cpo address
	locationBody := tools.GETRequest("http://localhost:3000/api/store/locations/" + cpoWallet)
	var locations []tools.XLocation
	err0 := json.Unmarshal(locationBody, &locations)
	if err0 != nil {
		log.Error(err0)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}
	var scIds []string
	for _, location := range locations {
		scIds = append(scIds, location.ScID)
	}

	// get the history
	body := tools.GETRequest("http://localhost:3000/api/cdr/info")

	var cdrs []tools.CDR
	err := json.Unmarshal(body, &cdrs)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	var cdrsOutput []tools.CDR

	for _, cdr := range cdrs {

		cdr.Currency = "Charge & Fuel Token"

		//TODO: removeme when fixing filtering by token contract is fixed
		if cdr.TokenContract == tokenContract {

			//is this record already present in some reimbursement ?
			count := 0
			row := tools.MDB.QueryRow("SELECT COUNT(*) as count FROM reimbursements WHERE cdr_records LIKE '%" + cdr.TransactionHash + "%'")
			row.Scan(&count)

			if count == 0 {
				log.Info("we have an unprocessed transaction: " + cdr.TransactionHash)

				isMyLocation := false
				for _, loc := range scIds {
					isMyLocation = loc == cdr.ScID
					if isMyLocation {

						//get the location name & address
						body = tools.GETRequest("http://localhost:3000/api/store/locations/" + cpoWallet + "/" + cdr.ScID)
						if body != nil {

							var loc tools.Location
							err := json.Unmarshal(body, &loc)
							if err != nil {
								log.Warnf(err.Error())
							} else {
								log.Info(loc)
								cdr.LocationName = loc.Name
								cdr.LocationAddress = loc.City + ", " + loc.Address + ", " + loc.Country
							}
						}

						cU, _ := strconv.Atoi(cdr.ChargedUnits)
						cdr.ChargedUnits = fmt.Sprintf("%d ", cU/1000)
						cdrsOutput = append(cdrsOutput, cdr)
					}
				}
			}
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

	//get the external ip of the server
	externalIp := "http://" + string(tools.GetExternalIp()) + ":9090"

	query := "INSERT INTO reimbursements ( msp_name, cpo_name, amount, currency, status, reimbursement_id, timestamp, cdr_records, server_addr, txs_number, token_address)" +
		"  VALUES ('%s','%s',%d,'%s','%s','%s',%d,'%s','%s', %d, '%s')"
	command := fmt.Sprintf(query, mspAddress, cpoWallet, totalAmount, "Charge and Fuel Token", "pending", tools.GetSha1Hash(cdrsOutput), time.Time.Unix(time.Now()), string(cdrsOutputBytes), externalIp, len(cdrsOutput), tokenContract)
	_, err = tools.MDB.Exec(command)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "there's already a reimbursement issued for the current transactions."})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Generate the PDFs

	reimbursementId := tools.GetSha1Hash(cdrsOutput)
	log.Info(">>> Generating PDF for " + reimbursementId)

	var reimb tools.Reimbursement

	err = tools.MDB.QueryRowx("SELECT * FROM reimbursements WHERE reimbursement_id = ? LIMIT 1", reimbursementId).StructScan(&reimb)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//load the template file
	b, err := ioutil.ReadFile("configs/invoice_template.html")
	tools.ErrorCheck(err, "cpo.go", false)

	cpo := tools.CPO{}
	tools.DB.QueryRowx("SELECT * FROM cpo LIMIT 1").StructScan(&cpo)

	rand.Seed(time.Now().UnixNano())
	randInt1 := strconv.Itoa(rand.Intn(900000))
	randInt2 := strconv.Itoa(rand.Intn(500000))
	randInt3 := strconv.Itoa(rand.Intn(200000))

	htmlTemplateRaw := string(b)

	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromName}}", cpo.Name, 2)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromAddress}}", cpo.Address1+" "+cpo.Address2+" "+cpo.Town, 2)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceDate}}", time.Now().Add(time.Hour * 1).Format(time.RFC1123), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceNumber}}", randInt1, 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{clientReference}}", "S&C"+randInt2, 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{purchaseOrder}}", "S&C"+randInt3, 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToName}}", "Volkswagen Financial", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToAddress}}", "Services (UKJ) Limited", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToPerson}}", "Milton Keynes", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToCode}}", "MK15 8HG", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceDueDate}}", time.Now().Add(time.Hour * 24 * 7 * time.Duration(2)).Format(time.RFC1123), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{description}}", "Sum of Tokens received through Share&Charge Network ", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{quantity}}", strconv.Itoa(reimb.Amount), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{unit}}", "Tokens", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{price}}", "£ 0,01", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{vat}}", "20%", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{total}}", fmt.Sprintf("£ %.2f", float64(reimb.Amount)*0.01), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{subTotal}}", fmt.Sprintf("£ %.2f", float64(reimb.Amount)*0.01), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{totalVat}}", fmt.Sprintf("£ %.2f", float64(reimb.Amount)*0.01*0.2), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{totalAmount}}", fmt.Sprintf("£ %.2f", float64(reimb.Amount)*0.01*1.2), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromPhone}}", "0014 882 739 2282", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromMail}}", cpo.MailAddr, 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromWebsite}}", cpo.Website, 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromBankName}}", "UNICREDIT", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromSortCode}}", "12312", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceFromAccountNumber}}", "123 345 532", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{vatNumber}}", "321ADF23", 1)

	//write it to a file
	ioutil.WriteFile("static/invoice_"+reimbursementId+".html", []byte(htmlTemplateRaw), 0644)
	time.Sleep(time.Second * 1)

	//convert it to pdf
	log.Info("Trying to convert it to pdf using chrome headless -> static/invoice_" + reimbursementId + ".html")
	err = tools.GeneratePdf("static/invoice_"+reimbursementId+".html", "static/invoice_"+reimbursementId+".pdf")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "reimbursement created", "pdf_location": reimb.ServerAddr + "/static/invoice_" + reimbursementId + ".pdf"})

}

// Lists all reimbursements
func CpoGetAllReimbursements(c *gin.Context) {

	status := c.Param("status")

	config := configs.Load()
	cpoWallet := config.GetString("cpo.wallet_address")

	var reimb []tools.Reimbursement

	err := tools.MDB.Select(&reimb, "SELECT * FROM reimbursements WHERE cpo_name = ? AND status = ?", cpoWallet, status)
	tools.ErrorCheck(err, "cpo.go", false)

	var output []tools.Reimbursement
	for k, reim := range reimb {
		reim.Index = k
		output = append(output, reim)
	}

	if len(reimb) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	c.JSON(http.StatusOK, output)
}

// marks the reimbursement as complete
func CpoSetReimbursementStatus(c *gin.Context) {

	reimbursementId := c.Param("reimbursement_id")
	reimbursementStatus := c.Param("status")

	rows, err := tools.MDB.Query("SELECT id FROM reimbursements WHERE reimbursement_id = ?", reimbursementId)
	tools.ErrorCheck(err, "cpo.go", false)
	defer rows.Close()

	//check if we have a reimbursement with this id present
	if !rows.Next() {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "there's isn't any reimbursement with this id present"})
		return
	}

	query := "UPDATE reimbursements SET status='%s' WHERE reimbursement_id = '%s'"
	command := fmt.Sprintf(query, reimbursementStatus, reimbursementId)
	_, err = tools.MDB.Exec(command)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if reimbursementStatus == "complete" {
		reimbursementId := c.Param("reimbursement_id")

		var reimb tools.Reimbursement
		err = tools.MDB.QueryRowx("SELECT * FROM reimbursements WHERE reimbursement_id =  \"" + reimbursementId + "\"").StructScan(&reimb)
		tools.ErrorCheck(err, "cpo.go", false)

		//get current token balance of the account
		config := configs.Load()
		cpoWallet := config.GetString("cpo.wallet_address")
		body := tools.GETRequest("http://localhost:3000/api/token/balance/" + cpoWallet)
		tokenBalanceFloat, _ := strconv.ParseFloat(string(body), 64)

		if tokenBalanceFloat < float64(reimb.Amount) {
			log.Error(err)
			c.JSON(http.StatusNotAcceptable, gin.H{"error": fmt.Sprintf("you are trying to send %d while you have only %f", reimb.Amount, tokenBalanceFloat)})
			return
		}

		log.Info(reimb)
		log.Warnf("sending now to CPO (hardcoded address) %d", reimb.Amount)

		_, err = tools.POSTRequest("http://localhost:3000/api/token/transfer/0xf60b71a4d360a42ec9d4e7977d8d9928fd7c8365/"+strconv.Itoa(reimb.Amount), nil)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": reimbursementId + " sent " + strconv.Itoa(reimb.Amount) + " tokens transferred to the MSP address"})
	}

	c.JSON(http.StatusOK, gin.H{"status": reimbursementStatus})

}

// marks the reimbursement as complete
func CpoSendTokensToMsp(c *gin.Context) {

	reimbursementId := c.Param("reimbursement_id")

	rows, err := tools.MDB.Query("SELECT id FROM reimbursements WHERE reimbursement_id = ?", reimbursementId)
	tools.ErrorCheck(err, "cpo.go", false)
	defer rows.Close()

	//check if we have a reimbursement with this id present
	if !rows.Next() {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "there's isn't any reimbursement with this id present"})
		return
	}

	var reimb tools.Reimbursement
	err = tools.MDB.QueryRowx("SELECT * FROM reimbursements WHERE reimbursement_id =  \"" + reimbursementId + "\"").StructScan(&reimb)
	tools.ErrorCheck(err, "cpo.go", false)

	//get current token balance of the account
	config := configs.Load()
	cpoWallet := config.GetString("cpo.wallet_address")
	body := tools.GETRequest("http://localhost:3000/api/token/balance/" + cpoWallet)
	tokenBalanceFloat, _ := strconv.ParseFloat(string(body), 64)

	if tokenBalanceFloat < float64(reimb.Amount) {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"error": fmt.Sprintf("you are trying to send %d while you have only %f", reimb.Amount, tokenBalanceFloat)})
		return
	}

	log.Info(reimb)
	log.Warnf("sending now to CPO (hardcoded address) %d", reimb.Amount)

	_, err = tools.POSTRequest("http://localhost:3000/api/token/transfer/0xf60b71a4d360a42ec9d4e7977d8d9928fd7c8365/"+strconv.Itoa(reimb.Amount), nil)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": reimbursementId + " sent " + strconv.Itoa(reimb.Amount) + " tokens transferred to the MSP address"})

}

// the records for the particular token
func CpoPaymentCDR(c *gin.Context) {

	tokenAddress := c.Param("token")
	config := configs.Load()
	cpoAddress := config.GetString("cpo.wallet_address")

	log.Info("loading all cpo's locations, this might take some time...")
	locationBody := tools.GETRequest("http://localhost:3000/api/store/locations/" + cpoAddress)
	var locations []tools.XLocation
	err0 := json.Unmarshal(locationBody, &locations)
	if err0 != nil {
		log.Error(err0)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	var scIds []string

	for _, location := range locations {
		scIds = append(scIds, location.ScID)
	}

	log.Info("loading all cdrs, this might take some time...")
	body := tools.GETRequest("http://localhost:3000/api/cdr/info") //+ ?tokenContract= tokenAddress

	var cdrs []tools.CDR
	err := json.Unmarshal(body, &cdrs)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	var cdrsOutput []tools.CDR

	if len(cdrs) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	for _, cdr := range cdrs {

		cdr.Currency = "Charge & Fuel Token"

		var count int
		err = tools.MDB.QueryRowx("SELECT COUNT(*) as count FROM reimbursements WHERE cdr_records LIKE '%" + cdr.TransactionHash + "%'").Scan(&count)
		tools.ErrorCheck(err, "cpo.go", false)

		if count == 0 {
			log.Info("seems we have an unprocessed tx.")
			//todo: this should be removed once filtering is fixed
			if cdr.TokenContract == tokenAddress {
				log.Info("seems we have an unprocessed tx. that maches our token contract")
				isMyLocation := false
				for _, loc := range scIds {
					isMyLocation = loc == cdr.ScID
					if isMyLocation {
						//get the location name & address
						body = tools.GETRequest("http://localhost:3000/api/store/locations/" + cpoAddress + "/" + cdr.ScID)
						if body != nil {

							var loc tools.Location
							err := json.Unmarshal(body, &loc)
							if err != nil {
								log.Warnf(err.Error())
							} else {
								log.Info(loc)
								cdr.LocationName = loc.Name
								cdr.LocationAddress = loc.City + ", " + loc.Address + ", " + loc.Country
							}
						}

						cU, _ := strconv.ParseFloat(cdr.ChargedUnits, 64)
						cdr.ChargedUnits = fmt.Sprintf("%.3f", cU / 1000)
						cdrsOutput = append(cdrsOutput, cdr)
					}
				}
			}
		} else {
			log.Warn("transaction with hash " + cdr.TransactionHash + " already present in some reimbursement")
		}

	}

	if len(cdrsOutput) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	//reverse the cdrOutput to get the latest on top
	for i, j := 0, len(cdrsOutput)-1; i < j; i, j = i+1, j-1 {
		cdrsOutput[i], cdrsOutput[j] = cdrsOutput[j], cdrsOutput[i]
	}

	c.JSON(http.StatusOK, cdrsOutput)

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

//=================================
//=========== LOCATIONS ===========
//=================================

//gets all locations of this CPO
func CpoGetLocations(c *gin.Context) {

	config := configs.Load()
	cpoAddress := config.GetString("cpo.wallet_address")
	body := tools.GETRequest("http://localhost:3000/api/store/locations/" + cpoAddress)
	if body == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	var locations []tools.XLocation
	err := json.Unmarshal(body, &locations)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, []string{})
		return
	}

	if len(locations) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	c.JSON(http.StatusOK, locations)

}

//uploads new locations
func CpoPostLocations(c *gin.Context) {
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
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

//uploads 1 location
func CpoPostLocation(c *gin.Context) {
	var stations tools.Location

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

	_, err = tools.POSTRequest("http://localhost:3000/api/store/location", jsonValue)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

//uploads 1 location
func CpoPutLocation(c *gin.Context) {
	var stations tools.Location
	scId := c.Param("scid")

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

	_, err = tools.PUTRequest("http://localhost:3000/api/store/location/"+scId, jsonValue)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

//deletes a location
func CpoDeleteLocation(c *gin.Context) {

	locationid := c.Param("locationid")

	_, err := tools.DELETERequest("http://localhost:3000/api/store/location/" + locationid)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})

}

//=================================
//=========== TARIFFS =============
//=================================

//gets all tariffs of this CPO
func CpoGetTariffs(c *gin.Context) {

	config := configs.Load()
	cpoAddress := config.GetString("cpo.wallet_address")
	body := tools.GETRequest("http://localhost:3000/api/store/tariffs/" + cpoAddress + "?raw=true")

	var tariffs []tools.Tariff
	err := json.Unmarshal(body, &tariffs)
	if err != nil {
		log.Warn(err)
		c.JSON(http.StatusOK, []string{})
		return
	}

	if len(tariffs) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	c.JSON(http.StatusOK, tariffs)

}

//uploads new tariffs and re-writes if they already are present
func CpoPutTariff(c *gin.Context) {
	var tariffs []tools.Tariff

	if err := c.MustBindWith(&tariffs, binding.JSON); err == nil {
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonValue, err := json.Marshal(tariffs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = tools.PUTRequest("http://localhost:3000/api/store/tariffs", jsonValue)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

//uploads new tariff
func CpoPostTariff(c *gin.Context) {
	var stations []tools.Tariff

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

	_, err = tools.POSTRequest("http://localhost:3000/api/store/tariffs?raw=true", jsonValue)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

//deletes tariffs. All of them!
func CpoDeleteTariffs(c *gin.Context) {

	_, err := tools.DELETERequest("http://localhost:3000/api/store/tariffs")
	if err != nil {
		log.Error(err)
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
