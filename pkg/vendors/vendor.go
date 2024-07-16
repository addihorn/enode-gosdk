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

const (
	REST_VENDOR_TRANSFER_ERROR     string = "vendors: could not read vendors"
	REST_VENDOR_READ_ERROR         string = "vendors: could not read response body"
	REST_VENDOR_PARSE_ERROR        string = "vendors: unable to parse vendor data"
	REST_VENDOR_UNAUTHORIZED_ERROR string = "vendors: unauthorized access"
	REST_VENDOR_GENERAL_ERROR      string = "vendors: some kind of error occured"
	REST_VENDOR_NO_VENDOR_ERROR    string = "vendors: no vendor with this id found"
	REST_VENDOR_VALLIDATION_ERROR  string = "vendors: invalid request payload input"
)
