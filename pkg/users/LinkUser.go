package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/addihorn/enode-gosdk/pkg/session"
)

/*
Creates a short-lived (24 hours), single-use device linking session.
Use the returned linkUrl to present Link UI to your user via [mobile in-app browsers] or [web redirects], or use the linkToken to present Link UI via the [Link SDKs].

Parameters:
  - sess: A pointer to the session object containing authentication and environment details.
  - data: A pointer to the LinkData object containing the necessary data for linking. If no error occured, the object data will be updated with the generated linkUrl and linkToken.

Returns:
  - An error if any occurred during the linking process. If no error occurred, it returns nil.

[mobile in-app browsers]: https://developers.enode.com/docs/link-ui#mobile-in-app-browsers
[web redirects]: https://developers.enode.com/docs/link-ui#web-redirects
[Link SDKs]: https://developers.enode.com/docs/link-ui#mobile-sd-ks
*/
func (user *User) Link(sess *session.Session, data *LinkData) error {
	url := fmt.Sprintf("%s/users/%s/link", sess.Authentication.Environment, user.Id)

	requestBody, err := json.Marshal(data)

	if err != nil {
		return errors.Join(errors.New("users: unable to create payload for link user service"), err)
	}
	fmt.Printf("%s\n", requestBody)

	req, _ := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sess.Authentication.Access_token))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(REST_USER_TRANSFER_ERROR)
		return errors.Join(errors.New(REST_USER_TRANSFER_ERROR), err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s: %s\n", REST_USER_READ_ERROR, err)
		return errors.Join(errors.New(REST_USER_READ_ERROR), err)
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
	case http.StatusOK:
		if err := json.Unmarshal(resBody, &data.LinkAccessData); err != nil {
			fmt.Printf("%s: %s\n", REST_USER_PARSE_ERROR, err)
			return errors.Join(errors.New(REST_USER_PARSE_ERROR), err)
		}
		return nil
	}

}
