package users

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/vendors"
)

/*
Disconnect a single Vendor from the User's account.

All stored data about their Vendor account will be deleted, and any assets that were provided by that Vendor will disappear from the system.

Parameters:
  - sess: A pointer to a session.Session object containing the necessary authentication and environment information.
  - vendor: A string representing the vendor to be disconnected.

Returns:
  - An error object if the request fails or encounters an unsuccessful status code.
    If the request is successful, the function returns nil.
*/
func (user *User) DisconnectVendor(sess *session.Session, vendor string) error {
	url := fmt.Sprintf("%s/users/%s/vendors/%s", sess.Authentication.Environment, user.Id, vendor)

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
