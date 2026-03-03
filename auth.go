package backend

import (
	"fmt"

	"github.com/kohirens/sso"
)

// AuthManager handles storing and retrieval of OIDC providers when an endpoint
// handler function is called. Granting the ability to authenticate the request.
type AuthManager interface {
	Add(name string, provider sso.OIDCProvider)
	Get(name string) (sso.OIDCProvider, error)
}

// Authorizer A default authorization manager.
type Authorizer struct {
	providers map[string]sso.OIDCProvider
}

// NewAuthManager Return an initialized default authorization manager.
func NewAuthManager() AuthManager {
	return &Authorizer{
		providers: make(map[string]sso.OIDCProvider),
	}
}

// Add Store an OIDC provider to retrieve for a later time.
func (ap *Authorizer) Add(name string, provider sso.OIDCProvider) {
	ap.providers[name] = provider
}

// Get Return an OIDC provider or throw an error.
func (ap *Authorizer) Get(name string) (sso.OIDCProvider, error) {
	// get from session which one the user chose.
	p, ok := ap.providers[name]
	if !ok {
		return nil, fmt.Errorf(stderr.ProviderNotFound, name)
	}
	return p, nil
}
