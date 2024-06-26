package session

import "github.com/addihorn/enode-gosdk/pkg/auth"

type Session struct {
	authentication *auth.Authentication
}

func NewSession(authSession *auth.Authentication) *Session {
	return &Session{authentication: authSession}
}
