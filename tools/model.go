package tools

import "time"

type Location struct {
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
		StatusSchedule []interface{} `json:"status_schedule"`
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