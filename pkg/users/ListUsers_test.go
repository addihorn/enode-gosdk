package users_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/addihorn/enode-gosdk/pkg/auth"
	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/users"
)

func TestListUsers_JSONUnmarshallingError(t *testing.T) {
	// Create a test server that will return a response with invalid JSON
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{invalid: "json"}`)
	}))
	defer ts.Close()

	// Create a session with the test server's URL
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  ts.URL,
			Access_token: "test_token",
		},
	}

	// Call the GetUsers function and capture the error
	_, err := users.ListUsers(sess)

	// Check if the error is not nil and contains the expected error message
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedError := errors.Join(errors.New(users.REST_USER_PARSE_ERROR), errors.New("invalid character 'i' looking for beginning of object key string"))
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestListUsers_NetworkError(t *testing.T) {
	// Create a test server that will return a network error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Network Error", http.StatusBadGateway)
	}))
	defer ts.Close()
	// Create a session with the test server's URL
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  ts.URL,
			Access_token: "test_token",
		},
	}

	// Call the GetUsers function and capture the error
	_, err := users.ListUsers(sess)

	// Check if the error is not nil and contains the expected error message
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedError := errors.Join(errors.New(users.REST_USER_TRANSFER_ERROR), errors.New("Get "+ts.URL+"/users: Bad Gateway"))
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestListUsers_EmptyResponse(t *testing.T) {
	// Create a test server that will return an empty response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a session with the test server's URL
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  ts.URL,
			Access_token: "test_token",
		},
	}

	// Call the GetUsers function and capture the error
	_, err := users.ListUsers(sess)

	// Check if the error is not nil and contains the expected error message
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedError := errors.Join(errors.New(users.REST_USER_READ_ERROR), io.EOF)
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestListUsers_Unauthorized(t *testing.T) {
	// Create a test server that will return an unauthorized response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	// Create a session with the test server's URL
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  ts.URL,
			Access_token: "test_token",
		},
	}

	// Call the GetUsers function and capture the error
	_, err := users.ListUsers(sess)

	// Check if the error is not nil and contains the expected error message
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedError := errors.Join(errors.New(users.REST_USER_UNAUTHORIZED_ERROR), fmt.Errorf("401 Unauthorized"))
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestListUsers_GeneralError(t *testing.T) {
	// Create a test server that will return a general error response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	// Create a session with the test server URL
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  ts.URL,
			Access_token: "test_token",
		},
	}

	// Call GetUserById function with a valid user ID
	_, err := users.ListUsers(sess)

	// Check if the error is not nil and contains the expected error message
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedError := errors.Join(errors.New(users.REST_USER_GENERAL_ERROR), fmt.Errorf("500 Internal Server Error"))
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestListUsers_Success(t *testing.T) {
	// Create a test server that will return a valid response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"data": [{"id": "1", "createdAt": "2022-01-01T00:00:00Z"}]}`)
	}))
	defer ts.Close()

	// Create a session with the test server's URL
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  ts.URL,
			Access_token: "test_token",
		},
	}

	// Call the GetUsers function and capture the result
	usr, err := users.ListUsers(sess)

	// Check if the error is nil and the result contains the expected user
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedUser := &users.User{
		Id:        "1",
		CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if _, ok := usr["1"]; !ok {
		t.Error("Expected user not found in result")
	}
	if usr["1"].Id != expectedUser.Id || usr["1"].CreatedAt != expectedUser.CreatedAt {
		t.Errorf("Expected user: %+v, but got: %+v", expectedUser, usr["1"])
	}
}
