package backend

import (
	"fmt"

	"github.com/kohirens/sso"
)

// ProviderManager handles storing and retrieval of OIDC providers when an endpoint
// handler function is called. Granting the ability to authenticate the request.
type ProviderManager interface {
	Add(name string, provider sso.OIDCProvider)
	Get(name string) (sso.OIDCProvider, error)
}

// OIDCProvider A default authorization manager.
type OIDCProvider struct {
	providers map[string]sso.OIDCProvider
}

// NewProviderManager Return an initialized default authorization manager.
func NewProviderManager() ProviderManager {
	return &OIDCProvider{
		providers: make(map[string]sso.OIDCProvider),
	}
}

// Add Store an OIDC provider to retrieve for a later time.
func (ap *OIDCProvider) Add(name string, provider sso.OIDCProvider) {
	ap.providers[name] = provider
}

// Get Return an OIDC provider or throw an error.
func (ap *OIDCProvider) Get(name string) (sso.OIDCProvider, error) {
	// get from session which one the user chose.
	p, ok := ap.providers[name]
	if !ok {
		return nil, fmt.Errorf(stderr.ProviderNotFound, name)
	}
	return p, nil
}
