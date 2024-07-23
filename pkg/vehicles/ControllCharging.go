package vehicles

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/addihorn/enode-gosdk/internal/actions"
	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/users"
)

type state string

const (
	on  state = "START"
	off state = "STOP"
)

type chargeAction struct {
	Action state `json:"action,omitempty"`
}

type VehicleAction struct {
	actions.DeviceAction
	Kind string `json:"kind"`
}

func (veh *Vehicle) StartCharging(sess *session.Session) (*VehicleAction, error) {
	startCharge := chargeAction{Action: on}
	return controllCharging(sess, veh, startCharge)

}

func (veh *Vehicle) StopCharging(sess *session.Session) (*VehicleAction, error) {
	stopCharge := chargeAction{Action: off}
	return controllCharging(sess, veh, stopCharge)
}

func controllCharging(sess *session.Session, veh *Vehicle, chargeStatus chargeAction) (*VehicleAction, error) {

	requestBody, err := json.Marshal(chargeStatus)

	if err != nil {
		return nil, errors.Join(errors.New("users: unable to create payload for link user service"), err)
	}

	url := fmt.Sprintf("%s/vehicles/%s/charging", sess.Authentication.Environment, veh.Id)
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sess.Authentication.Access_token))

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(REST_VEHICLE_TRANSFER_ERROR)
		return nil, errors.Join(errors.New(REST_VEHICLE_TRANSFER_ERROR), err)
	}

	switch resp.StatusCode {
	// TODO: Handle Response Codes 400 and 422
	default:
		return nil, errors.Join(fmt.Errorf(REST_VEHICLE_GENERAL_ERROR+"\n %+v", resp))
	case http.StatusNotFound:
		return nil, errors.Join(errors.New(users.REST_USER_NO_USERS_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusUnauthorized:
		return nil, errors.Join(errors.New(REST_VEHICLE_UNAUTHORIZED_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusBadGateway:
		return nil, errors.Join(errors.New(REST_VEHICLE_TRANSFER_ERROR), fmt.Errorf("Get %s: Bad Gateway", url))
	case http.StatusInternalServerError:
		return nil, errors.Join(errors.New(REST_VEHICLE_GENERAL_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusOK:
		return readVehicleActionPayload(resp)
	}

}

func readVehicleActionPayload(resp *http.Response) (*VehicleAction, error) {
	bodyPayload, err := getResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var action *VehicleAction
	if err := json.Unmarshal(bodyPayload, &action); err != nil {
		return nil, errors.Join(errors.New(actions.REST_ACTIONS_PARSE_ERROR), err)
	}
	return action, nil

}
