package vehicles

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/users"
)

func ListUserVehicles(sess *session.Session, userId string) (map[string]*Vehicle, error) {
	url := fmt.Sprintf("%s/users/%s/vehicles", sess.Authentication.Environment, userId)
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
	case http.StatusNotFound:
		return nil, errors.Join(errors.New(users.REST_USER_NO_USERS_ERROR), fmt.Errorf("%+v", resp.Status))
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
