package session

import "github.com/addihorn/enode-gosdk/pkg/auth"

type Session struct {
	Authentication *auth.Authentication
}

func NewSession(authSession *auth.Authentication) *Session {
	return &Session{Authentication: authSession}
}
