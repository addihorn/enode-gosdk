package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/addihorn/enode-gosdk/pkg/auth"
	"github.com/addihorn/enode-gosdk/pkg/environments"
)

func BenchmarkMain(b *testing.B) {

	client_id := "<YOUR-CLIENT_ID>"
	client_secret := "<YOU-CLIENT-SECRET>"

	session, err := auth.NewAuthentication(client_id, client_secret, environments.SANDBOX, false)
	if err != nil {
		fmt.Println("main: was not able to create a new session", err)
	}

	fmt.Printf("%+v\n", session)

	endTimer := time.NewTimer(60 * time.Second)

	<-endTimer.C
}
