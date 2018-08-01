package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"fmt"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"github.com/motionwerkGmbH/cpo-backend-api/configs"
	"math/rand"
	"math"
	"strconv"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"net/url"
	"github.com/davecgh/go-spew/spew"
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

	type History struct {
		From      string  `json:"from"`
		Amount    float64 `json:"amount"`
		Currency  string  `json:"currency"`
		Timestamp int64   `json:"timestamp"`
	}
	//

	var cDb *tools.CouchDB
	cDb, err := tools.Database("18.197.172.83", 5984)
	tools.ErrorCheck(err, "general.go", false)

	err = cDb.SelectDb("blockchain", "admin", "hardpassword1")
	tools.ErrorCheck(err, "general.go", false)

	type FindResponse struct {
		TotalRows int `json:"total_rows"`
		Offset    int `json:"offset"`
		Rows []struct {
			ID    string   `json:"id"`
			Key   []string `json:"key"`
			Value string   `json:"value"`
			Doc struct {
				ID        string  `json:"_id"`
				Rev       string  `json:"_rev"`
				Block     int     `json:"block"`
				From      string  `json:"from"`
				To        string  `json:"to"`
				Amount    float64 `json:"amount"`
				Currency  string  `json:"currency"`
				GasUsed   string  `json:"gas_used"`
				GasPrice  string  `json:"gas_price"`
				Timestamp string  `json:"timestamp"`
			} `json:"doc"`
		} `json:"rows"`
	}

	var findResult FindResponse
	params := url.Values{}
	var addrx []string
	config := configs.Load()
	addrx = append(addrx, config.GetString("cpo.wallet_address"))
	data, _ := json.Marshal(addrx)
	params.Set("key", string(data))
	params.Set("include_docs", "true")
	params.Set("descending", "true")
	err = cDb.Db.GetView("doc", "history_of_account", &findResult, &params)

	tools.ErrorCheck(err, "general.go", false)

	if len(findResult.Rows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no transactions found for this address"})
		return
	}

	var histories []History
	for _, row := range findResult.Rows {
		if configs.AddressToName(row.Doc.To) != config.GetString("cpo.wallet_address") {

			if row.Doc.Amount >= 1000000000000000000 {
				row.Doc.Amount = row.Doc.Amount / 1000000000000000000
				row.Doc.Currency = "ETH"
			}
			n := History{From: configs.AddressToName(row.Doc.From), Amount: row.Doc.Amount, Currency: row.Doc.Currency, Timestamp: tools.HexToInt(row.Doc.Timestamp)}
			histories = append(histories, n)
		}

	}

	spew.Dump(histories)

	type WalletRecord struct {
		MspName           string `json:"msp_name"`
		TotalTransactions int    `json:"total_transactions"`
		Amount            int64  `json:"amount"`
		Currency          string `json:"currency"`
		TokenAddr         string `json:"token_address"`
	}
	var walletRecords []WalletRecord

	//total transaction

	record := WalletRecord{MspName: "Charge & Fuel", TotalTransactions: 16, Amount: 52, Currency: "C&F Tokens", TokenAddr: "0xA39A488BEf3EC11be06AA3B24fe8a51c9F899205"}
	walletRecords = append(walletRecords, record)

	c.JSON(http.StatusOK, walletRecords)
}

//reimbursements
// creates a Reimbursement
func CpoCreateReimbursement(c *gin.Context) {

	//gets the MSP address from url
	mspAddress := c.Param("msp_address")

	type History struct {
		From      string  `json:"from"`
		Amount    float64 `json:"amount"`
		Currency  string  `json:"currency"`
		Timestamp int64   `json:"timestamp"`
	}

	type Reinbursment struct {
		ReimbursementId string    `json:"reimbursement_id"`
		From            string    `json:"from"`
		To              string    `json:"to"`
		Amount          int64     `json:"amount"`
		Currency        string    `json:"currency"`
		CreatedAt       int64     `json:"created_at"`
		Status          string    `json:"status"`
		History         []History `json:"history"`
	}

	var reimbursement Reinbursment
	config := configs.Load()

	//-------- gets the History of the account

	var cDb *tools.CouchDB
	cDb, err := tools.Database("18.197.172.83", 5984)
	tools.ErrorCheck(err, "general.go", false)

	err = cDb.SelectDb("blockchain", "admin", "hardpassword1")
	tools.ErrorCheck(err, "general.go", false)

	type FindResponse struct {
		TotalRows int `json:"total_rows"`
		Offset    int `json:"offset"`
		Rows []struct {
			ID    string   `json:"id"`
			Key   []string `json:"key"`
			Value string   `json:"value"`
			Doc struct {
				ID        string  `json:"_id"`
				Rev       string  `json:"_rev"`
				Block     int     `json:"block"`
				From      string  `json:"from"`
				To        string  `json:"to"`
				Amount    float64 `json:"amount"`
				Currency  string  `json:"currency"`
				GasUsed   string  `json:"gas_used"`
				GasPrice  string  `json:"gas_price"`
				Timestamp string  `json:"timestamp"`
			} `json:"doc"`
		} `json:"rows"`
	}

	var findResult FindResponse
	params := url.Values{}
	var addrx []string
	addrx = append(addrx, config.GetString("cpo.wallet_address"))
	data, _ := json.Marshal(addrx)
	params.Set("key", string(data))
	params.Set("include_docs", "true")
	params.Set("descending", "true")
	err = cDb.Db.GetView("doc", "history_of_account", &findResult, &params)

	tools.ErrorCheck(err, "cpo.go", false)

	if len(findResult.Rows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no transactions found for this address"})
		return
	}

	var histories []History
	for _, row := range findResult.Rows {
		if configs.AddressToName(row.Doc.From) != "MSP" {

			if row.Doc.Amount >= 1000000000000000000 {
				row.Doc.Amount = row.Doc.Amount / 1000000000000000000
				row.Doc.Currency = "ETH"
			}
			n := History{From: configs.AddressToName(row.Doc.From), Amount: row.Doc.Amount, Currency: row.Doc.Currency, Timestamp: tools.HexToInt(row.Doc.Timestamp)}
			histories = append(histories, n)
		}


	}
	//================= HISTORY ENDS ==========


	err = cDb.SelectDb("reimbursements", "admin", "hardpassword1")
	tools.ErrorCheck(err, "general.go", false)

	reimbursement.From = config.GetString("cpo.wallet_address")
	reimbursement.To = mspAddress
	reimbursement.Amount = 32
	reimbursement.Currency = "Charge & Fuel Token"
	reimbursement.CreatedAt = time.Time.Unix(time.Now())
	reimbursement.Status = "pending"
	reimbursement.History = histories
	reimbursement.ReimbursementId = tools.GetSha1Hash(histories)

	// Check if this reimbursement is already present

	type XResponse struct {
		TotalRows int `json:"total_rows"`
		Offset    int `json:"offset"`
		Rows []struct {} `json:"rows"`
	}

	//calls the unique view in couchdb
	var xResult XResponse
	var addry []string
	xparams := url.Values{}
	addry = append(addry, reimbursement.ReimbursementId)
	data, _ = json.Marshal(addry)
	xparams.Set("key", string(data))
	err = cDb.Db.GetView("doc", "unique", &xResult, &xparams)


	if len(xResult.Rows) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "there's already a reimbursement issued for the current transactions."})
		return
	}


	revId, err := cDb.Insert(reimbursement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"_revId": revId})
}

// the records for the particular token
func CpoPaymentCDR(c *gin.Context) {

	type CDRRecord struct {
		Date       string `json:"date"`
		DriverName string `json:"driver_name"`
		Amount     int64  `json:"amount"`
		Currency   string `json:"currency"`
		EvseId     string `json:"evseid"`
	}
	var cdrRecords []CDRRecord

	record := CDRRecord{Date: "2018/03/18 18:22:33", DriverName: "Joh Lewis", Amount: 52, Currency: "C&F Tokens", EvseId: "213kdfs93"}
	cdrRecords = append(cdrRecords, record)

	record = CDRRecord{Date: "2018/03/18 18:32:54", DriverName: "Joh Lewis", Amount: 12, Currency: "C&F Tokens", EvseId: "2321213"}
	cdrRecords = append(cdrRecords, record)

	c.JSON(http.StatusOK, cdrRecords)
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
		Amount    float64 `json:"amount"`
		Currency  string  `json:"currency"`
		Timestamp string  `json:"timestamp"`
	}

	s1 := rand.NewSource(1337)
	r1 := rand.New(s1)

	var histories []History
	for i := 0; i < 100; i++ {
		n := History{Amount: math.Floor(r1.Float64()*10000) / 10000, Currency: "CPO Tokens", Timestamp: "01.04.2018 " + strconv.Itoa(10+r1.Intn(23)) + ":" + strconv.Itoa(10+r1.Intn(49)) + ":" + strconv.Itoa(10+r1.Intn(49))}
		histories = append(histories, n)
	}

	c.JSON(http.StatusOK, histories)
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
