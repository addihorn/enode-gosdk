package vendors

type Vendor struct {
	Vendor  string     `json:"vendor,omitempty"`
	Type    VendorType `json:"vendorType,omitempty"`
	IsValid bool       `json:"isValid,omitempty"`
}

type VendorType string

const (
	VEHICLE  VendorType = "vehicle"
	CHARGER             = "charger"
	HVAC                = "hvac"
	INVERTER            = "inverter"
	BATTERY             = "battery"
	METER               = "meter"
)
