package vehicles

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/addihorn/enode-gosdk/internal/devices"
)

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
		devices.BasicInformation
		VIN         string `json:"vin"`
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
	Location     devices.Location              `json:"location"`
	Capabilities map[string]devices.Capability `json:"capabilities"`
}

const (
	REST_VEHICLE_TRANSFER_ERROR     string = "vehicles: could not read vehicles"
	REST_VEHICLE_READ_ERROR         string = "vehicles: could not read response body"
	REST_VEHICLE_PARSE_ERROR        string = "vehicles: unable to parse vehicle data"
	REST_VEHICLE_UNAUTHORIZED_ERROR string = "vehicles: unauthorized access"
	REST_VEHICLE_GENERAL_ERROR      string = "vehicles: some kind of error occured"
	REST_VEHICLE_VALLIDATION_ERROR  string = "vehicles: invalid request payload input"
	REST_VEHICLE_NO_VEHICLE_ERROR   string = "vehciles: no vehicle with this id found"
)

func getResponseBody(resp *http.Response) ([]byte, error) {
	if resp.ContentLength == 0 {
		return nil,
			errors.Join(errors.New(REST_VEHICLE_READ_ERROR), io.EOF)
	}

	if resBody, err := io.ReadAll(resp.Body); err != nil {
		fmt.Printf("%s: %s\n", REST_VEHICLE_READ_ERROR, err)
		return nil, errors.Join(errors.New(REST_VEHICLE_READ_ERROR), err)
	} else {
		return resBody, nil
	}
}
