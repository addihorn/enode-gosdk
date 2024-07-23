package devices

import "time"

type Capability struct {
	InterventionIds []string `json:"interventionIds"`
	IsCapable       bool     `json:"isCapable"`
}

type Location struct {
	Longitude   float64   `json:"longitude"`
	Latitude    float64   `json:"latitude"`
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
}

type BasicInformation struct {
	Brand string `json:"brand"`
	Model string `json:"model"`
}
