package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/addihorn/enode-gosdk/pkg/enums/languages"
	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/vendors"
)

type Data struct {
	Data []*User `json:"data"`
}

type User struct {
	Id            string           `json:"id"`
	CreatedAt     time.Time        `json:"createdAt"`
	LinkedVendors []vendors.Vendor `json:"linkedVendors,omitempty"`
}

type LinkAccess struct {
	LinkUrl   string `json:"linkUrl"`
	LinkToken string `json:"linkToken"`
}

type LinkData struct {
	Vendor         string             `json:"vendor,omitempty"`
	Type           vendors.VendorType `json:"vendorType"`
	Language       languages.Language `json:"language"`
	Scopes         []string           `json:"scopes"`
	RedirectUri    string             `json:"redirectUri"`
	ColorScheme    string             `json:"colorScheme,omitempty"`
	LinkAccessData LinkAccess         `json:"-"`
}

const (
	REST_USER_TRANSFER_ERROR     string = "users: could not read users"
	REST_USER_READ_ERROR         string = "users: could not read response body"
	REST_USER_PARSE_ERROR        string = "users: unable to parse user data"
	REST_USER_UNAUTHORIZED_ERROR string = "users: unauthorized access"
	REST_USER_GENERAL_ERROR      string = "users: some kind of error occured"
	DATA_USER_NO_USERS_ERROR     string = "users: no users with this id found"
)

func (user *User) LinkUser(sess *session.Session, data *LinkData) error {
	url := fmt.Sprintf("%s/users/%s/link", sess.Authentication.Environment, user.Id)

	requestBody, err := json.Marshal(data)

	if err != nil {
		return errors.Join(errors.New("users: unable to create payload for link user service"), err)
	}
	fmt.Printf("%s\n", requestBody)

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
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
	case http.StatusUnauthorized:
		return errors.Join(fmt.Errorf(REST_USER_UNAUTHORIZED_ERROR+"\n %+v", resp.Status))
	case http.StatusOK:
		json.Unmarshal(resBody, &data.LinkAccessData)
		return nil
	default:
		return errors.Join(fmt.Errorf(REST_USER_GENERAL_ERROR+"\n %+v", resp))
	}

}
