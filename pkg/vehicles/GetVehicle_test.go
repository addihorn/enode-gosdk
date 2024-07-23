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
	"github.com/addihorn/enode-gosdk/pkg/vehicles"
)

func TestGetVehicle_StatusOK(t *testing.T) {
	// Arrange
	expectedVehicle := &vehicles.Vehicle{
		Id:          "vehicle1",
		Vendor:      "TESLA",
		UserId:      "fobbar",
		IsReachable: true,
		LastSeen:    time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC)}
	jsonResponse, _ := json.Marshal(expectedVehicle)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}

	// Act
	actualVehicle, err := vehicles.GetVehicle(sess, "vehicle1")

	// Assert
	if err != nil {
		t.Fatalf("Error calling GetVehicle: %v", err)
	}

	if !reflect.DeepEqual(actualVehicle, expectedVehicle) {
		t.Fatalf("Expected vehicle \n%+v\n, got \n%+v\n", expectedVehicle, actualVehicle)
	}
}

func TestGetVehicle_StatusUnauthorized(t *testing.T) {
	// Arrange
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

	// Act
	_, err := vehicles.GetVehicle(sess, "123")

	// Assert
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}

	expectedError := errors.Join(errors.New(vehicles.REST_VEHICLE_UNAUTHORIZED_ERROR), fmt.Errorf("%d %+v", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)))
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected error %v, got %v", expectedError, err)
	}
}

func TestGetVehicle_StatusBadGateway(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	sess := &session.Session{
		Authentication: &auth.Authentication{
			Environment:  server.URL,
			Access_token: "test_token",
		},
	}

	vehicleId := "123"

	// Act

	_, err := vehicles.GetVehicle(sess, vehicleId)

	// Assert
	if err == nil {
		t.Fatal("Expected error for bad gateway, got nil")
	}

	expectedError := errors.Join(errors.New(vehicles.REST_VEHICLE_TRANSFER_ERROR), fmt.Errorf("Get %s: Bad Gateway", server.URL+"/vehicles/"+vehicleId))
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected error %v, got %v", expectedError, err)
	}
}

func TestGetVehicle_StatusInternalServerError(t *testing.T) {
	// Arrange
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

	// Act
	_, err := vehicles.GetVehicle(sess, "123")

	// Assert
	if err == nil {
		t.Fatal("Expected error for internal server error, got nil")
	}

	expectedError := errors.Join(errors.New(vehicles.REST_VEHICLE_GENERAL_ERROR), fmt.Errorf("%d %+v", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected error %v, got %v", expectedError, err)
	}
}

func TestGetVehicle_VehicleNotFound(t *testing.T) {
	// Arrange
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

	// Act
	_, err := vehicles.GetVehicle(sess, "123")

	// Assert
	if err == nil {
		t.Fatal("Expected error for internal server error, got nil")
	}

	expectedError := errors.Join(errors.New(vehicles.REST_VEHICLE_NO_VEHICLE_ERROR), fmt.Errorf("404 Not Found"))
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected error %v, got %v", expectedError, err)
	}
}
