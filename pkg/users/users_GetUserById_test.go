package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/addihorn/enode-gosdk/pkg/auth"
	"github.com/addihorn/enode-gosdk/pkg/session"
)

func TestGetUserById_ParseError(t *testing.T) {
	// Create a test server that will return a JSON response with invalid format
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{invalid_json}`)
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
	_, err := GetUserById(sess, "test_user_id")

	// Check if the error is not nil and contains the expected error message
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedError := errors.Join(errors.New(REST_USER_PARSE_ERROR), errors.New("invalid character 'i' looking for beginning of object key string"))
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestGetUserById_EmptyResponse(t *testing.T) {
	// Create a test server that will return an empty response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
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
	_, err := GetUserById(sess, "test_user_id")

	// Check if the error is not nil and contains the expected error message
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedError := errors.Join(errors.New(REST_USER_READ_ERROR), io.EOF)
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestGetUserById_NetworkError(t *testing.T) {
	userId := "test_user_id"
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

	// Call GetUserById function with a valid user ID
	_, err := GetUserById(sess, userId)

	// Check if the error is not nil and contains the expected error message
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedError := errors.Join(errors.New(REST_USER_TRANSFER_ERROR), fmt.Errorf("Get %s/users/%s: Bad Gateway", ts.URL, userId))
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestGetUserById_Unauthorized(t *testing.T) {
	// Create a test server that will return an unauthorized response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
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
	_, err := GetUserById(sess, "test_user_id")

	// Check if the error is not nil and contains the expected error message
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedError := errors.Join(errors.New(REST_USER_UNAUTHORIZED_ERROR), fmt.Errorf("401 Unauthorized"))
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestGetUserById_GeneralError(t *testing.T) {
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
	_, err := GetUserById(sess, "test_user_id")

	// Check if the error is not nil and contains the expected error message
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedError := errors.Join(errors.New(REST_USER_GENERAL_ERROR), fmt.Errorf("500 Internal Server Error"))
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestGetUserById_Success(t *testing.T) {
	// Arrange
	userId := "test-user-id"
	expectedUser := &User{Id: userId, CreatedAt: time.Now()}
	jsonResponse, _ := json.Marshal(expectedUser)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}))
	defer ts.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  ts.URL,
			Access_token: "test-access-token",
		},
	}

	// Act
	actualUser, err := GetUserById(sess, userId)

	// Assert
	if err != nil {
		t.Errorf("GetUserById returned an error: %v", err)
	}
	if actualUser.Id != expectedUser.Id {
		t.Errorf("GetUserById returned incorrect user ID. Expected %s, got %s", expectedUser.Id, actualUser.Id)
	}
}
