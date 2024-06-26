package main

import (
	"fmt"

	"github.com/addihorn/enode-gosdk/pkg/auth"
	"github.com/addihorn/enode-gosdk/pkg/environments"
	"github.com/addihorn/enode-gosdk/pkg/session"
)

func main() {

	client_id := "<YOUR-CLIENT_ID>"
	client_secret := "<YOU-CLIENT-SECRET>"

	authentication, err := auth.NewAuthentication(client_id, client_secret, environments.SANDBOX, true)
	if err != nil {
		fmt.Println("main: was not able to create a new session \n", err)
	}

	_ = session.NewSession(authentication)
}
