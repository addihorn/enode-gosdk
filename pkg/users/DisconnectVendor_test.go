package users_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/addihorn/enode-gosdk/pkg/auth"
	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/users"
	"github.com/addihorn/enode-gosdk/pkg/vendors"
)

func TestDisconnectVendor_StatusOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}
	err := user.DisconnectVendor(sess, "CUPRA")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

}

func TestDisconnectVendor_StatusNoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}
	err := user.DisconnectVendor(sess, "TESLA")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDisconnectVendor_StatusNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}
	user := &users.User{Id: "unknownUser", CreatedAt: time.Now()}
	err := user.DisconnectVendor(sess, "HUAWEI")

	expectedError := errors.Join(errors.New(users.REST_USER_NO_USERS_ERROR), fmt.Errorf("%s", "404 Not Found"))
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error\n%v, \ngot\n%v", expectedError, err)
	}
}

func TestDisconnectVendor_StatusBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}
	err := user.DisconnectVendor(sess, "ZZZ_WRONG_VENDOR_ZZZ")

	expectedError := errors.Join(errors.New(vendors.REST_VENDOR_NO_VENDOR_ERROR), fmt.Errorf("%s", "400 Bad Request"))
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error\n%v, \ngot\n%v", expectedError, err)
	}
}

func TestDisconnectVendor_StatusUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "ZZZ_WRONG_TOKEN_ZZZ",
		},
	}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}
	err := user.DisconnectVendor(sess, "TESLA")

	expectedError := errors.Join(errors.New(users.REST_USER_UNAUTHORIZED_ERROR), fmt.Errorf("%s", "401 Unauthorized"))
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error\n%v foo, \ngot\n%v foo", expectedError, err)
	}
}
