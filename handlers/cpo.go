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
	"log"
	"encoding/json"
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

//generates a new wallet for the cpo
func CpoGenerateWallet(c *gin.Context){

	type WalletInfo struct {
		Seed   string `json:"seed"`
		Addr string `json:"address"`
	}
	var walletInfo WalletInfo

	//body := tools.GetRequest("http://localhost:3000/api/wallet/create")
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
		c.JSON(http.StatusNotFound, gin.H{"error": "there isn't any seed in the msp account. Maybe you need to create the wallet first ?."})
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
		n := History{Amount: math.Floor(r1.Float64()*10000) / 10000, Currency: "MSP Tokens", Timestamp: "01.04.2018 " + strconv.Itoa(10+r1.Intn(23)) + ":" + strconv.Itoa(10+r1.Intn(49)) + ":" + strconv.Itoa(10+r1.Intn(49))}
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
	body := tools.GetRequest("http://localhost:3000/api/store/locations/"+cpoAddress)

	var locations []tools.Location
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
func CpoPutLocations(c *gin.Context){
	var stations []tools.Location

	if err := c.MustBindWith(&stations, binding.JSON); err == nil {
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stations)
}

//uploads new location
func CpoPostLocation(c *gin.Context){
	var station tools.Location

	if err := c.MustBindWith(&station, binding.JSON); err == nil {
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, station)
}

//deletes a location
func CpoDeleteLocation(c *gin.Context) {

	locationid := c.Param("locationid")

	c.JSON(http.StatusBadRequest, gin.H{"error": locationid})
}


//uploads new location
func CpoPostEvse(c *gin.Context){
	var evse tools.Evse

	if err := c.MustBindWith(&evse, binding.JSON); err == nil {
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, evse)
}

