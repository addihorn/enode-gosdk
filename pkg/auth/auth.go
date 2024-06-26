package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/addihorn/enode-gosdk/pkg/environments"
)

type Authentication struct {
	Access_token string
	Scope        string
	Token_type   string
}

func NewAuthentication(client_id, client_secret, environment string, automaticTokenRefresh bool) (*Authentication, error) {

	authData, err := authenticate(client_id, client_secret, environment)

	if err != nil {
		return nil, errors.Join(errors.New("authentication: could not get a new authentication session"), err)
	}

	auth := &Authentication{}

	if err = json.Unmarshal(authData, auth); err != nil {
		return nil, errors.Join(errors.New("authentication: error, while trying to unmarshal response from auth service"), err)
	}

	if automaticTokenRefresh {
		var f map[string]interface{}

		if err = json.Unmarshal(authData, &f); err != nil {
			fmt.Println(errors.Join(errors.New("authentication:unable to parse expiration-timer"), err))
		}
		expires_in, _ := strconv.Atoi(fmt.Sprint(f["expires_in"]))

		time.AfterFunc(
			time.Duration(expires_in-30)*time.Second,
			func() { auth.refreshToken(client_id, client_secret, environments.SANDBOX) },
		)
	}

	return auth, nil
}

func (sess *Authentication) refreshToken(client_id, client_secret, environment string) {

	// fmt.Printf("Timer has been fired with token %s \n Refreshing Token...\n", sess.Access_token)

	authData, err := authenticate(client_id, client_secret, environment)
	if err != nil {
		fmt.Println(errors.Join(errors.New("authentication: could not get a new authentication session"), err))
	}

	var f map[string]interface{}
	if err = json.Unmarshal(authData, &f); err != nil {
		fmt.Println(errors.Join(errors.New("authentication: error, while trying to unmarshal response from auth service"), err))
	}

	if err = json.Unmarshal(authData, sess); err != nil {
		fmt.Println(errors.Join(errors.New("authentication:unable to parse expiration-timer"), err))
	}

	expires_in, _ := strconv.Atoi(fmt.Sprint(f["expires_in"]))

	time.AfterFunc(
		time.Duration(expires_in-30)*time.Second,
		func() { sess.refreshToken(client_id, client_secret, environment) },
	)
}

func authenticate(client_id, client_secret, environment string) ([]byte, error) {

	authUrl := fmt.Sprintf("https://oauth.%s.enode.io/oauth2/token", environment)

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authUrl, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(client_id, client_secret)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("client: could not execute authentication request: \n")
		return nil, errors.Join(errors.New("client: could not execute authentication request"), err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return nil, errors.Join(errors.New("client: could not read response body: \n"), err)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	return resBody, nil
}
