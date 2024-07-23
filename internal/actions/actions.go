package actions

import "time"

type DeviceAction struct {
	Id            string      `json:"id"`
	UserId        string      `json:"userId"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	CompletedAt   time.Time   `json:"completedAt,omitempty"`
	State         actionState `json:"state"`
	DeviceId      string      `json:"targetId"`
	DeviceKind    string      `json:"targetKind"`
	FailureReason struct {
		Type   string `json:"type"`
		Detail string `json:"detail,omitempty"`
	} `json:"failureReason,omitempty"`
}

type actionState string

const (
	DEVICE_ACTION_PENDING   actionState = "PENDING"
	DEVICE_ACTION_CONFIRMED actionState = "CONFIRMED"
	DEVICE_ACTION_FAILED    actionState = "FAILED"
	DEVICE_ACTION_CANCELLED actionState = "CANCELLED"
)

const (
	REST_ACTIONS_PARSE_ERROR string = "action: unable to parse actions data"
)
