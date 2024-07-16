package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/addihorn/enode-gosdk/pkg/session"
)

func readUsersPayload(resp *http.Response) (map[string]*User, error) {

	if resp.ContentLength == 0 {
		return nil,
			errors.Join(errors.New(REST_USER_READ_ERROR), io.EOF)
	}

	var bodyPayload []byte
	if resBody, err := io.ReadAll(resp.Body); err != nil {
		fmt.Printf("%s: %s\n", REST_USER_READ_ERROR, err)
		return nil, errors.Join(errors.New(REST_USER_READ_ERROR), err)
	} else {
		bodyPayload = resBody
	}

	fmt.Printf("%s\n", bodyPayload)

	var userData Data
	if err := json.Unmarshal(bodyPayload, &userData); err != nil {
		return nil, errors.Join(errors.New(REST_USER_PARSE_ERROR), err)
	}

	userCache := make(map[string]*User)
	for _, user := range userData.Data {
		userCache[user.Id] = user
	}

	return userCache, nil
}

/*
Returns a paginated list of all users.

Parameters:
  - sess: A pointer to a session.Session object containing the authentication details and environment URL.

Returns:
  - A map of user IDs to User structs, or nil if an error occurs.
  - An error, or nil if the operation is successful.
*/
func ListUsers(sess *session.Session) (map[string]*User, error) {

	url := fmt.Sprintf("%s/users", sess.Authentication.Environment)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sess.Authentication.Access_token))

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(REST_USER_TRANSFER_ERROR)
		return nil, errors.Join(errors.New(REST_USER_TRANSFER_ERROR), err)
	}

	switch resp.StatusCode {
	default:
		return nil, errors.Join(fmt.Errorf(REST_USER_GENERAL_ERROR+"\n %+v", resp))
	case http.StatusUnauthorized:
		return nil, errors.Join(errors.New(REST_USER_UNAUTHORIZED_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusBadGateway:
		return nil, errors.Join(errors.New(REST_USER_TRANSFER_ERROR), fmt.Errorf("Get %s: Bad Gateway", url))
	case http.StatusInternalServerError:
		return nil, errors.Join(errors.New(REST_USER_GENERAL_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusOK:
		return readUsersPayload(resp)
	}

}
