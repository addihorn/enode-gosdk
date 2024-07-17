package main

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/addihorn/enode-gosdk/pkg/auth"
	"github.com/addihorn/enode-gosdk/pkg/enums/environments"
	"github.com/addihorn/enode-gosdk/pkg/enums/languages"
	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/users"
	"github.com/addihorn/enode-gosdk/pkg/vehicles"
	"github.com/addihorn/enode-gosdk/pkg/vendors"
	"github.com/joho/godotenv"
)

func BenchmarkMain(b *testing.B) {

	client_id := "<YOUR-CLIENT_ID>"
	client_secret := "<YOUR-CLIENT-SECRET>"

	session, err := auth.NewAuthentication(client_id, client_secret, environments.SANDBOX, false)
	if err != nil {
		fmt.Println("main: was not able to create a new session", err)
	}

	fmt.Printf("%+v\n", session)

	endTimer := time.NewTimer(60 * time.Second)

	<-endTimer.C
}

func Test_UserCRUD(t *testing.T) {
	//change directory to direct to your project directory
	err := godotenv.Load("/workspaces/enode-gosdk/.env")

	if err != nil {
		t.Errorf("integration: was not able to read .env file")
	}

	client_id := os.Getenv("ENODE_CLIENT_ID")
	client_secret := os.Getenv("ENODE_CLIENT_SECRET")

	authentication, err := auth.NewAuthentication(client_id, client_secret, environments.SANDBOX, true)
	if err != nil {
		t.Errorf("integration: was not able to create a new session:\n%+v\n", err)
	}

	// get all users
	sess := session.NewSession(authentication)
	userList, _ := users.ListUsers(sess)
	fmt.Printf("%+v\n", userList)

	//link user to new devices
	user := &users.User{Id: "foobar"}
	linkData := users.LinkData{
		Type:        vendors.BATTERY,
		Language:    languages.ENGLISH_UK,
		Scopes:      []string{"battery:read:data"},
		RedirectUri: "http://localhost:3000",
	}
	fmt.Printf("%+v\n", user.Link(sess, &linkData)) // print error
	fmt.Printf("%+v\n", linkData.LinkAccessData)    // print link data

	// when using net/http redirect to url linkData.LinkAccessData.LinkUrl

	//read user foobar
	// get specific user
	user, err = users.GetUser(sess, user.Id)
	if err == nil {
		fmt.Printf("User Data: %+v\n", user)
	} else {
		t.Errorf("integration: unable to read users data:\n%+v\n", err)
		fmt.Println(err.Error())
	}
	//unlink user foobar

	if err := user.Unlink(sess); err != nil {
		t.Errorf("integration: error while unlinking user:\n%+v\n", err)
	}

	user, err = users.GetUser(sess, user.Id)

	expectedError := errors.Join(errors.New(users.REST_USER_NO_USERS_ERROR), fmt.Errorf("404 Not Found"))
	if expectedError.Error() != err.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}

	if user != nil {
		t.Errorf("User should have been deleted\n%+v\n", user)
	}
}

func Test_VehicleCRUD(t *testing.T) {
	//change directory to direct to your project directory
	err := godotenv.Load("/workspaces/enode-gosdk/.env")

	if err != nil {
		t.Errorf("integration: was not able to read .env file")
	}

	client_id := os.Getenv("ENODE_CLIENT_ID")
	client_secret := os.Getenv("ENODE_CLIENT_SECRET")

	authentication, err := auth.NewAuthentication(client_id, client_secret, environments.SANDBOX, true)
	if err != nil {
		t.Errorf("integration: was not able to create a new session:\n%+v\n", err)
	}

	// get all vehicles
	sess := session.NewSession(authentication)
	vehicleList, _ := vehicles.ListVehicles(sess)
	fmt.Printf("%+v\n", vehicleList)

}
