package vehicles_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/addihorn/enode-gosdk/pkg/auth"
	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/users"
	"github.com/addihorn/enode-gosdk/pkg/vehicles"
)

func TestListUserVehicles_Success(t *testing.T) {
	// Arrange
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  "https://api.enode.com",
			Access_token: "valid_token",
		},
	}

	expectedVehicles := map[string]*vehicles.Vehicle{
		"vehicle1": {
			Id:          "vehicle1",
			Vendor:      "TESLA",
			UserId:      "fobbar",
			IsReachable: true,
			LastSeen:    time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC),
		},
	}
	expectedVehicles["vehicle1"].Information.VIN = "ABC123DEF4567890"
	expectedVehicles["vehicle1"].ChargeState.Range = 255
	expectedVehicles["vehicle1"].Odometer.Distance = 1234567890

	vehicleList := make([]*vehicles.Vehicle, len(expectedVehicles))
	i := 0
	for _, vehicle := range expectedVehicles {
		vehicleList[i] = vehicle
		i++
	}

	jsonData, err := json.Marshal(vehicles.Data{Data: vehicleList})
	if err != nil {
		t.Fatalf("Error marshaling expected vehicles: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}))
	defer server.Close()

	sess.Authentication.Environment = server.URL

	// Act
	actualVehicles, err := vehicles.ListUserVehicles(sess, "user1")

	// Assert
	if err != nil {
		t.Fatalf("Error calling ListVehicles: %v", err)
	}

	if len(actualVehicles) != len(expectedVehicles) {
		t.Fatalf("Expected %d vehicles, got %d", len(expectedVehicles), len(actualVehicles))
	}

	for id, expectedVehicle := range expectedVehicles {
		actualVehicle, ok := actualVehicles[id]
		if !ok {
			t.Fatalf("Expected vehicle with ID %s not found", id)
		}

		if !reflect.DeepEqual(actualVehicle, expectedVehicle) {
			t.Fatalf("Expected vehicle \n%+v\n, got \n%+v\n", expectedVehicle, actualVehicle)
		}
	}
}

func TestListUserVehicles_InvalidToken(t *testing.T) {
	// Arrange
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  "https://api.enode.com",
			Access_token: "valid_token",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	sess.Authentication.Environment = server.URL

	// Act
	_, err := vehicles.ListUserVehicles(sess, "user1")

	// Assert
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}

	expectedError := errors.Join(errors.New(vehicles.REST_VEHICLE_UNAUTHORIZED_ERROR), fmt.Errorf("%d %+v", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)))
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected error %v, got %v", expectedError, err)
	}
}

func TestListUserVehicles_BadGateway(t *testing.T) {
	// Arrange
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  "https://api.enode.com",
			Access_token: "valid_token",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	sess.Authentication.Environment = server.URL

	// Act
	userId := "superUser123"
	_, err := vehicles.ListUserVehicles(sess, userId)

	// Assert
	if err == nil {
		t.Fatal("Expected error for bad gateway, got nil")
	}

	expectedError := errors.Join(errors.New(vehicles.REST_VEHICLE_TRANSFER_ERROR), fmt.Errorf("Get %s: Bad Gateway", server.URL+"/users/"+userId+"/vehicles"))
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected error %v, got %v", expectedError, err)
	}
}

func TestListUserVehicles_InternalServerError(t *testing.T) {
	// Arrange
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  "https://api.enode.com",
			Access_token: "valid_token",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	sess.Authentication.Environment = server.URL

	// Act
	_, err := vehicles.ListUserVehicles(sess, "user1")

	// Assert
	if err == nil {
		t.Fatal("Expected error for internal server error, got nil")
	}

	expectedError := errors.Join(errors.New(vehicles.REST_VEHICLE_GENERAL_ERROR), fmt.Errorf("%d %+v", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected error %v, got %v", expectedError, err)
	}
}

func TestListUserVehicles_InvalidJSON(t *testing.T) {
	// Arrange
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  "https://api.enode.com",
			Access_token: `valid_token`,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid_json"))
	}))
	defer server.Close()

	sess.Authentication.Environment = server.URL

	// Act
	_, err := vehicles.ListUserVehicles(sess, "user1")

	// Assert
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}

	expectedError := errors.Join(errors.New(vehicles.REST_VEHICLE_PARSE_ERROR), errors.New("invalid character 'i' looking for beginning of value"))
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected error %v, got %v", expectedError, err)
	}
}

func TestListUserVehicles_UserNotFound(t *testing.T) {
	// Arrange
	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  "https://api.enode.com",
			Access_token: `valid_token`,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("invalid_json"))
	}))
	defer server.Close()

	sess.Authentication.Environment = server.URL

	// Act
	_, err := vehicles.ListUserVehicles(sess, "user1")

	// Assert
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}

	expectedError := errors.Join(errors.New(users.REST_USER_NO_USERS_ERROR), fmt.Errorf("404 Not Found"))
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected error %v, got %v", expectedError, err)
	}
}
