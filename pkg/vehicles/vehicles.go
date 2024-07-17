package vehicles

import "time"

type Data struct {
	Data []*Vehicle `json:"data"`
}

type Vehicle struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userId"`
	Vendor      string    `json:"vendor"`
	IsReachable bool      `json:"isReachable"`
	LastSeen    time.Time `json:"lastSeen"`
	Information struct {
		VIN         string `json:"vin"`
		Brand       string `json:"brand"`
		Model       string `json:"model"`
		Year        int    `json:"year"`
		DisplayName string `json:"displayName"`
	} `json:"information"`
	ChargeState struct {
		BatteryLevel        float64   `json:"batteryLevel"`
		Range               float64   `json:"range"`
		IsPluggedIn         bool      `json:"isPluggedIn"`
		IsCharging          bool      `json:"isCharging"`
		IsFullyCharged      bool      `json:"isFullyCharged"`
		BatteryCapacity     float64   `json:"batteryCapacity"`
		ChargeLimit         float64   `json:"chargeLimit"`
		ChargeRate          float64   `json:"chargeRate"`
		ChargeTimeRemaining float64   `json:"chargeTimeRemaining"`
		LastUpdated         time.Time `json:"lastUpdated"`
		MaxCurrent          float64   `json:"maxCurrent"`
		PowerDeliveryState  string    `json:"powerDeliveryState"`
	} `json:"chargeState"`
	Odometer struct {
		Distance    float64   `json:"distance"`
		LastUpdated time.Time `json:"lastUpdated"`
	} `json:"odometer"`
	Location struct {
		Longitude   float64   `json:"longitude"`
		Latitude    float64   `json:"latitude"`
		LastUpdated time.Time `json:"lastUpdated"`
	} `json:"location"`
}

const (
	REST_VEHICLE_TRANSFER_ERROR     string = "vehicles: could not read vehicles"
	REST_VEHICLE_READ_ERROR         string = "vehicles: could not read response body"
	REST_VEHICLE_PARSE_ERROR        string = "vehicles: unable to parse vehicle data"
	REST_VEHICLE_UNAUTHORIZED_ERROR string = "vehicles: unauthorized access"
	REST_VEHICLE_GENERAL_ERROR      string = "vehicles: some kind of error occured"
	REST_VEHICLE_NO_USERS_ERROR     string = "vehicles: no vehicles with this id found"
	REST_VEHICLE_VALLIDATION_ERROR  string = "vehicles: invalid request payload input"
)
