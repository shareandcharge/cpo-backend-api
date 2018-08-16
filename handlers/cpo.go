package handlers

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/motionwerkGmbH/cpo-backend-api/configs"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
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

// creates a Reimbursement
func CpoCreateReimbursement(c *gin.Context) {

	//gets the MSP address from url
	mspAddress := c.Param("msp_address")
	tokenContract := c.DefaultQuery("tokenContract", "")

	config := configs.Load()
	cpoWallet := config.GetString("cpo.wallet_address")

	// get the history

	body := tools.GETRequest("http://localhost:3000/api/cdr/info")

	var cdrs []tools.CDR
	err := json.Unmarshal(body, &cdrs)
	if err != nil {
		log.Panic(err)
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
				cdrsOutput = append(cdrsOutput, cdr)
			} else {
				log.Warn("transaction with hash " + cdr.TransactionHash + " already present in a reimbursement")
			}
		} else {
			log.Warnf("tx doesn't have the required token contract")
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

	query := "INSERT INTO reimbursements ( msp_name, cpo_name, amount, currency, status, reimbursement_id, timestamp, cdr_records, txs_number, server_addr)" +
		"  VALUES ('%s','%s',%d,'%s','%s','%s',%d,'%s',%d,'%s')"
	command := fmt.Sprintf(query, mspAddress, cpoWallet, totalAmount, "Charge&Fuel Token", "pending", tools.GetSha1Hash(cdrsOutput), time.Time.Unix(time.Now()), string(cdrsOutputBytes), len(cdrsOutput), externalIp)
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

	//if we set it as complete, transfer the coins to the MSP
	if reimbursementStatus == "complete" {

		var reimb tools.Reimbursement
		err := tools.MDB.Select(&reimb, "SELECT * FROM reimbursements WHERE reimbursement_id = ?", reimbursementId)
		tools.ErrorCheck(err, "cpo.go", false)

		log.Warnf("should send now to the msp the ammount %d, modifty the code please", reimb.Amount)

		_, err = tools.POSTRequest("http://localhost:3000/api/token/transfer/0xf60b71a4d360a42ec9d4e7977d8d9928fd7c8365/10", nil)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": reimbursementStatus + " and 10 tokens transferred to the MSP address"})

}

// the records for the particular token
func CpoPaymentCDR(c *gin.Context) {

	tokenAddress := c.Param("token")
	config := configs.Load()
	cpoAddress := config.GetString("cpo.wallet_address")

	body := tools.GETRequest("http://localhost:3000/api/cdr/info") //+ ?tokenContract= tokenAddress

	var cdrs []tools.CDR
	err := json.Unmarshal(body, &cdrs)
	if err != nil {
		log.Panic(err)
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
		err = tools.MDB.QueryRowx("SELECT COUNT(*) as count FROM reimbursements WHERE cdr_records LIKE '%" + cdr.TransactionHash + "%'").StructScan(&count)
		tools.ErrorCheck(err, "cpo.go", false)

		if count == 0 {

			//todo: this should be removed once filtering is fixed
			if cdr.TokenContract == tokenAddress {
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
					cdrsOutput = append(cdrsOutput, cdr)
				}

				cdrsOutput = append(cdrsOutput, cdr)
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

//========= PDF Generation ========

func CpoReimbursementGenPdf(c *gin.Context) {
	reimbursementId := c.Param("reimbursement_id")

	var reimb tools.Reimbursement

	err := tools.MDB.QueryRowx("SELECT * FROM reimbursements WHERE reimbursement_id = ? LIMIT 1", reimbursementId).StructScan(&reimb)
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
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceDate}}", time.Now().Format(time.RFC822), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceNumber}}", randInt1, 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{clientReference}}", "S&C"+randInt2, 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{purchaseOrder}}", "S&C"+randInt3, 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToName}}", "Volkswagen Financial", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToAddress}}", "Services (UKJ) Limited", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToPerson}}", "Milton Keynes", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceToCode}}", "MK15 8HG", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{invoiceDueDate}}", "02 August 2018 ", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{description}}", "Sum of Tokens received through Share&Charge Network ", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{quantity}}", strconv.Itoa(reimb.Amount), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{unit}}", "Tokens", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{price}}", "£ 0,01", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{vat}}", "20%", 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{total}}", fmt.Sprintf("%.4f", float64(reimb.Amount)*0.001), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{subTotal}}", "£ "+strconv.Itoa(reimb.Amount), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{totalVat}}", "£ "+strconv.Itoa(int(float64(reimb.Amount)*float64(0.20))), 1)
	htmlTemplateRaw = strings.Replace(htmlTemplateRaw, "{{totalAmount}}", "£ "+strconv.Itoa(int(float64(reimb.Amount)*float64(1.20))), 1)
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
	log.Info("trying to convert it to pdf -> static/invoice_" + reimbursementId + ".html")
	err = tools.GeneratePdf("static/invoice_"+reimbursementId+".html", "static/invoice_"+reimbursementId+".pdf")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redirect": "http://{{server_addr}}:{{server_port}}/static/invoice_" + reimbursementId + ".pdf"})
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

//uploads new location
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
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

//uploads 1 location
func CpoPutLocation(c *gin.Context) {
	var stations tools.XLocation
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
		log.Panic(err)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "there aren't any tariffs registered with this CPO"})
		return
	}

	if len(tariffs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "there aren't any tariffs registered with this CPO"})
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
		log.Panic(err)
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
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

//deletes tariffs. All of them!
func CpoDeleteTariffs(c *gin.Context) {

	_, err := tools.DELETERequest("http://localhost:3000/api/store/tariffs")
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
