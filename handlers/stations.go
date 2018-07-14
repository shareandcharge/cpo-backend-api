package handlers

import (
	"github.com/gin-gonic/gin"
	"time"
	"github.com/motionwerkGmbH/cpo-backend-api/configs"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"encoding/json"
	"log"
	"net/http"
)


func StationsInfo(c *gin.Context) {

	type StationInfo []struct {
		ID          string `json:"id"`
		Type        string `json:"type"`
		Name        string `json:"name"`
		Address     string `json:"address"`
		City        string `json:"city"`
		PostalCode  string `json:"postal_code"`
		Country     string `json:"country"`
		Coordinates struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"coordinates"`
		Evses []struct {
			UID            string        `json:"uid"`
			EvseID         string        `json:"evse_id"`
			Status         string        `json:"status"`
			StatusSchedule []interface{} `json:"status_schedule,omitempty"`
			Capabilities   []interface{} `json:"capabilities"`
			Connectors     []struct {
				ID          string    `json:"id"`
				Standard    string    `json:"standard"`
				Format      string    `json:"format"`
				PowerType   string    `json:"power_type"`
				Voltage     int       `json:"voltage"`
				Amperage    int       `json:"amperage"`
				TariffID    string    `json:"tariff_id"`
				LastUpdated time.Time `json:"last_updated"`
			} `json:"connectors"`
			PhysicalReference string    `json:"physical_reference"`
			FloorLevel        string    `json:"floor_level"`
			LastUpdated       time.Time `json:"last_updated"`
		} `json:"evses"`
		Operator struct {
			Name string `json:"name"`
		} `json:"operator"`
		LastUpdated time.Time `json:"last_updated"`
	}

	Config := configs.Load()
	cpoAddress := Config.GetString("cpo_address")

	body := tools.GetRequest("http://localhost:3000/api/store/locations/"+ cpoAddress)

	var stationInfo StationInfo
	err := json.Unmarshal(body, &stationInfo)
	if err != nil {
		log.Panic(err)
	}
	c.JSON(http.StatusOK, stationInfo)

	//b, err := ioutil.ReadFile("./configs/sc_configs/locations.json")
	//if err != nil {
	//	fmt.Print(err)
	//}
	//
	//locationsString := string(b)
	//
	//var dat StationInfo
	//
	//json.Unmarshal([]byte(locationsString), &dat)
	//
	//c.JSON(http.StatusOK, dat)
}
