package main

import (
	"fmt"
	"os"

	"github.com/addihorn/enode-gosdk/pkg/auth"
	"github.com/addihorn/enode-gosdk/pkg/enums/environments"
	"github.com/addihorn/enode-gosdk/pkg/enums/languages"
	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/users"
	"github.com/addihorn/enode-gosdk/pkg/vendors"
	"github.com/joho/godotenv"
)

func main() {
	//change directory to direct to your project directory
	err := godotenv.Load("/workspaces/enode-gosdk/.env")

	if err != nil {
		fmt.Errorf(err.Error())
	}

	client_id := os.Getenv("ENODE_CLIENT_ID")
	client_secret := os.Getenv("ENODE_CLIENT_SECRET")

	authentication, err := auth.NewAuthentication(client_id, client_secret, environments.SANDBOX, true)
	if err != nil {
		fmt.Println("main: was not able to create a new session \n", err)
	}

	// get all users
	sess := session.NewSession(authentication)
	userList, _ := users.ListUsers(sess)
	fmt.Printf("%+v\n", userList)

	// get specific user
	user, err := users.GetUser(sess, "1ab23cd4")
	if err == nil {
		fmt.Printf("User Data: %+v\n", user)
	} else {
		fmt.Println(err.Error())
	}

	//link user to new devices
	user = &users.User{Id: "foobar"}
	linkData := users.LinkData{
		Type:        vendors.BATTERY,
		Language:    languages.ENGLISH_UK,
		Scopes:      []string{"battery:read:data"},
		RedirectUri: "http://localhost:3000",
	}
	fmt.Printf("%+v\n", user.Link(sess, &linkData)) // print error
	fmt.Printf("%+v\n", linkData.LinkAccessData)    // print link data

	// when using net/http redirect to url linkData.LinkAccessData.LinkUrl

}
