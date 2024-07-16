package users

import (
	"time"

	"github.com/addihorn/enode-gosdk/pkg/enums/languages"
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
	REST_USER_TRANSFER_ERROR           string = "users: could not read users"
	REST_USER_READ_ERROR               string = "users: could not read response body"
	REST_USER_PARSE_ERROR              string = "users: unable to parse user data"
	REST_USER_UNAUTHORIZED_ERROR       string = "users: unauthorized access"
	REST_USER_GENERAL_ERROR            string = "users: some kind of error occured"
	REST_USER_NO_USERS_ERROR           string = "users: no users with this id found"
	REST_USER_VALLIDATION_ERROR        string = "users: invalid request payload input"
	REST_USER_CONNECTION_LIMIT_REACHED string = "users: connection limit for this user exceeded"
)
