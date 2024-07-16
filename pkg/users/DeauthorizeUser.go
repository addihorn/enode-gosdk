package users

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/addihorn/enode-gosdk/pkg/session"
)

/*
Deletes the User's stored vendor authorizations and credentials, invalidates any associated sessions, authorization codes, and access/refresh tokens.

All other User data is retained, and if the User is sent through the Link User flow in the future, their account will be just as they left it.

No webhook events will be generated for a deauthorized user.

Parameters:
  - sess: A pointer to a session.Session object containing the necessary authentication and environment information.

Returns:
  - An error if the request fails or the status code indicates an error.
    If the request is successful, it returns nil.
*/
func (user *User) Deauthorize(sess *session.Session) error {
	url := fmt.Sprintf("%s/users/%s/authorization", sess.Authentication.Environment, user.Id)

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
	case http.StatusOK, http.StatusNoContent:
		return nil
	}

}
