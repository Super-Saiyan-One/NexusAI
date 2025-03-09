// auth/auth.go
package auth

import "nexus-ai/api/hailuo/config"

type Authenticator struct{}

func NewAuthenticator() *Authenticator {
	return &Authenticator{}
}

func (a *Authenticator) AuthorizationHeader() string {
	return "Bearer " + config.API_KEY
}
