package users

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/addihorn/enode-gosdk/pkg/session"
)

/*
Deletes a User and all of their data permanently and invalidates any associated sessions, authorization codes, and access/refresh tokens.

Parameters:
- sess: A pointer to a session.Session object representing the user's session.

Returns:
  - An error if the request fails or the status code indicates an error.
    If the request is successful, it returns nil.
*/
func (user *User) Unlink(sess *session.Session) error {
	url := fmt.Sprintf("%s/users/%s", sess.Authentication.Environment, user.Id)

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
		return errors.Join(errors.New(REST_USER_VALLIDATION_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusForbidden:
		return errors.Join(errors.New(REST_USER_CONNECTION_LIMIT_REACHED), fmt.Errorf("%+v", resp.Status))
	case http.StatusOK, http.StatusNoContent:
		return nil
	}

}
