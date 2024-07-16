package users

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/vendors"
)

/*
Disconnect a specific vendor type from the User's account. Assets of this type from that Vendor will be removed. If no other types from that vendor remain, all its stored data will be deleted.

Parameters:
  - sess: A pointer to a session.Session object containing the necessary authentication and environment information.
  - vendor: A string representing the vendor to be disconnected.
  - venType: A string representing the type of vendor to be disconnected.

Returns:
  - An error object if the request fails or encounters an unsuccessful status code.
    If the request is successful, the function returns nil.
*/
func (user *User) DisconnectVendortype(sess *session.Session, vendor string, venType vendors.VendorType) error {

	fmt.Printf("type of venType: %+v\n", reflect.TypeOf(venType))

	url := fmt.Sprintf("%s/users/%s/vendors/%s/%s", sess.Authentication.Environment, user.Id, vendor, venType)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sess.Authentication.Access_token))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(REST_USER_TRANSFER_ERROR)
		return errors.Join(errors.New(REST_USER_TRANSFER_ERROR), err)
	}

	switch resp.StatusCode {
	default:
		return errors.Join(fmt.Errorf(REST_USER_GENERAL_ERROR+"\n %+v", resp))
	case http.StatusUnauthorized:
		return errors.Join(errors.New(REST_USER_UNAUTHORIZED_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusNotFound:
		return errors.Join(errors.New(REST_USER_NO_USERS_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusInternalServerError:
		return errors.Join(errors.New(REST_USER_GENERAL_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusBadRequest:
		return errors.Join(errors.New(vendors.REST_VENDOR_NO_VENDOR_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusForbidden:
		return errors.Join(errors.New(REST_USER_CONNECTION_LIMIT_REACHED), fmt.Errorf("%+v", resp.Status))
	case http.StatusOK, http.StatusNoContent:
		return nil
	}

}
