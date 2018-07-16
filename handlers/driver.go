package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func DriversList(c *gin.Context) {

	type Driver struct {
		Name string `json:"name"`
		Phone string `json:"phone"`
		Addr string `json:"addr"`
		Balance float64 `json:"balance"`
		Currency string `json:"currency"`
		TokenAddr string `json:"token_address"`
	}

	var drivers []Driver
	drivers = append(drivers, Driver{Name:"Nemanja", Phone:"3933892822",Addr:"0x0322323232",Balance:0.02203, Currency:"S&C Token", TokenAddr:"0x03232432423"})
	drivers = append(drivers, Driver{Name:"Andy", Phone:"393892822",Addr:"0x032323232",Balance:0.00003, Currency:"S&C Token", TokenAddr:"0x03232432423"})

	c.JSON(http.StatusOK, drivers)
}
