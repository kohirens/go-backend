package backend

import (
	"fmt"
	"net/url"

	"github.com/kohirens/sso/oidc"
)

// ProviderManager handles storing and retrieval of OIDC providers when an endpoint
// handler function is called. Granting the ability to authenticate the request.
type ProviderManager struct {
	providers map[string]oidc.Provider
}

func (pm *ProviderManager) Add(name string, provider oidc.Provider) {
	pm.providers[name] = provider
}

// AuthLinks generates an authentication link for every OIDC provider listed.
func (pm *ProviderManager) AuthLinks() map[string]string {
	links := make(map[string]string)

	for _, provider := range pm.providers {
		authLink, e1 := provider.AuthLink("")
		if e1 != nil {

		}
		links[provider.Name()] = authLink
	}

	return links
}

func (pm *ProviderManager) Callback(url url.Values) error {
	// TODO: Determine the provider from the request.
	providerName := "google"

	provider, ok := pm.providers[providerName]
	if !ok {
		return fmt.Errorf(stderr.ProviderNotFound, providerName)
	}

	return provider.Callback(url)
}

func (pm *ProviderManager) Get(name string) (oidc.Provider, error) {
	provider, ok := pm.providers[name]
	if !ok {
		return nil, fmt.Errorf(stderr.ProviderNotFound, name)

	}

	return provider, nil
}

// NewProviderManager Return an initialized default authorization manager.
func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		providers: make(map[string]oidc.Provider),
	}
}
