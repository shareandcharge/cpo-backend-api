package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func StationsInfo(c *gin.Context) {

	type StationInfo struct {
		ID          string `json:"id"`
		Type        string `json:"type"`
		Name        string `json:"name"`
		Address     string `json:"address"`
		Opened      string `json:"opened"`
		City        string `json:"city"`
		PostalCode  string `json:"postal_code"`
		Country     string `json:"country"`
		Coordinates struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"coordinates"`
	}

	var stationInfo StationInfo
	stationInfo.ID = "SC EXAMPLE 1"
	stationInfo.Type = "ON_STREET"
	stationInfo.Name = "MW1"
	stationInfo.Address = "Ruettenscheiderstr. 120"
	stationInfo.Opened = "24/7"
	stationInfo.City = "Essen"
	stationInfo.PostalCode = "1334 GE"
	stationInfo.Country = "Germany"
	stationInfo.Coordinates.Latitude = "51.432870"
	stationInfo.Coordinates.Longitude = "7.004115"


	var stationInfo2 StationInfo
	stationInfo2.ID = "SC EXAMPLE 2"
	stationInfo2.Type = "ON_STREET"
	stationInfo2.Name = "MW2"
	stationInfo2.Address = "Ruettenscheiderstr. 121"
	stationInfo2.Opened = "24/7"
	stationInfo2.City = "Essen"
	stationInfo2.PostalCode = "1334 GE"
	stationInfo2.Country = "Germany"
	stationInfo2.Coordinates.Latitude = "51.433870"
	stationInfo2.Coordinates.Longitude = "7.003115"

	var stationsInfos []StationInfo
	stationsInfos = append(stationsInfos, stationInfo, stationInfo2)

	c.JSON(http.StatusOK, stationsInfos)
}
