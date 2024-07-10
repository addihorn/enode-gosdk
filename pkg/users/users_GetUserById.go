package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/addihorn/enode-gosdk/pkg/session"
)

func readUserByIdPayload(resp *http.Response) (*User, error) {

	if resp.ContentLength == 0 {
		return nil, errors.Join(errors.New(REST_USER_READ_ERROR), io.EOF)
	}

	var bodyPayload []byte
	if resBody, err := io.ReadAll(resp.Body); err != nil {
		return nil, errors.Join(errors.New(REST_USER_READ_ERROR), err)
	} else {
		bodyPayload = resBody
	}

	fmt.Printf("%s\n", bodyPayload)

	var userData *User
	if err := json.Unmarshal(bodyPayload, &userData); err != nil {
		return nil, errors.Join(errors.New(REST_USER_PARSE_ERROR), err)
	}

	return userData, nil
}

func GetUserById(sess *session.Session, userId string) (*User, error) {

	url := fmt.Sprintf("%s/users/%s", sess.Authentication.Environment, userId)
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
	case http.StatusBadGateway:
		return nil, errors.Join(errors.New(REST_USER_TRANSFER_ERROR), fmt.Errorf("Get %s: Bad Gateway", url))
	case http.StatusUnauthorized:
		return nil, errors.Join(errors.New(REST_USER_UNAUTHORIZED_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusInternalServerError:
		return nil, errors.Join(errors.New(REST_USER_GENERAL_ERROR), fmt.Errorf("%+v", resp.Status))
	case http.StatusOK:
		return readUserByIdPayload(resp)
	}

}
