package vehicles

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/addihorn/enode-gosdk/pkg/session"
)

/*
EVs provide charge, location, and odometer data.
Vehicles can be controlled either directly using the Control [ChargingAPI] endpoint, or through [Smart Charging] and [Schedules].

[ChargingAPI]: https://developers.enode.com/api/reference#postVehiclesVehicleidCharging
[Smart Charging]: https://developers.enode.com/docs/smart-charging/introduction
[Schedules]: https://developers.enode.com/docs/scheduling
*/
func ListVehicles(sess *session.Session) (map[string]*Vehicle, error) {

	url := fmt.Sprintf("%s/vehicles", sess.Authentication.Environment)
	req, err := http.NewRequest("GET", url, nil)
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
	default:
		return nil, errors.Join(fmt.Errorf(REST_VEHICLE_GENERAL_ERROR+"\n %+v", resp))
	case http.StatusUnauthorized:
		return nil, errors.Join(errors.New(REST_VEHICLE_UNAUTHORIZED_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusBadGateway:
		return nil, errors.Join(errors.New(REST_VEHICLE_TRANSFER_ERROR), fmt.Errorf("Get %s: Bad Gateway", url))
	case http.StatusInternalServerError:
		return nil, errors.Join(errors.New(REST_VEHICLE_GENERAL_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusOK:
		return readVehiclesPayload(resp)
	}
}

func readVehiclesPayload(resp *http.Response) (map[string]*Vehicle, error) {
	bodyPayload, err := getResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var vehiclesData Data
	if err := json.Unmarshal(bodyPayload, &vehiclesData); err != nil {
		return nil, errors.Join(errors.New(REST_VEHICLE_PARSE_ERROR), err)
	}

	vehicleCache := make(map[string]*Vehicle)
	for _, user := range vehiclesData.Data {
		vehicleCache[user.Id] = user
	}

	return vehicleCache, nil
}
