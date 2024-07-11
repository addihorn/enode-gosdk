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
)

func TestLinkUser_StatusOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"linkUrl": "https://localhost/link-token", "linkToken": "abc"}`)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}

	data := &users.LinkData{}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}
	err := user.LinkUser(sess, data)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if data.LinkAccessData.LinkUrl != "https://localhost/link-token" || data.LinkAccessData.LinkToken != "abc" {
		t.Errorf("expected LinkAccessData to be set, got %+v", data.LinkAccessData)
	}
}

func TestLinkUser_StatusUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}

	data := &users.LinkData{}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}
	err := user.LinkUser(sess, data)

	expectedError := errors.Join(errors.New(users.REST_USER_UNAUTHORIZED_ERROR), fmt.Errorf("%s", "401 Unauthorized"))
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error\n%v foo, \ngot\n%v foo", expectedError, err)
	}
}

func TestLinkUser_StatusNotFound(t *testing.T) {
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

	data := &users.LinkData{}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}
	err := user.LinkUser(sess, data)

	expectedError := errors.Join(errors.New(users.REST_USER_NO_USERS_ERROR), fmt.Errorf("%s", "404 Not Found"))
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error\n%v, \ngot\n%v", expectedError, err)
	}
}

func TestLinkUser_StatusInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}

	data := &users.LinkData{}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}
	err := user.LinkUser(sess, data)

	expectedError := errors.Join(errors.New(users.REST_USER_GENERAL_ERROR), fmt.Errorf("%s", "500 Internal Server Error"))
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error\n%v, \ngot\n%v", expectedError, err)
	}
}

func TestLinkUser_InvalidResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `invalid json`)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}

	data := &users.LinkData{}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}
	err := user.LinkUser(sess, data)

	expectedError := errors.Join(errors.New(users.REST_USER_PARSE_ERROR), fmt.Errorf("%s", "invalid character 'i' looking for beginning of value"))
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error\n%v, \ngot\n%v", expectedError, err)
	}
}

func TestLinkUser_StatusBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error": "https://docs.enode.io/problems/validation-error", "messafe": "Multiple validation errors, see issues list."}`)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}

	data := &users.LinkData{}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}

	err := user.LinkUser(sess, data)

	expectedError := errors.Join(errors.New(users.REST_USER_VALLIDATION_ERROR), fmt.Errorf("%s", "400 Bad Request"))
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error\n%v, \ngot\n%v", expectedError, err)
	}
}

func TestLinkUser_StatusForbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{
  "type": "https://docs.enode.io/problems/forbidden",
  "title": "Connections limit reached.",
  "detail": "Unable to create more connections for ClientID: a7bedf14-c3eb-4c2b-a08f-b34a1f70808d"
}`)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}

	data := &users.LinkData{}
	user := &users.User{Id: "user-1", CreatedAt: time.Now()}

	err := user.LinkUser(sess, data)

	expectedError := errors.Join(errors.New(users.REST_USER_CONNECTION_LIMIT_REACHED), fmt.Errorf("%s", "403 Forbidden"))
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error\n%v, \ngot\n%v", expectedError, err)
	}
}
