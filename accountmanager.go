package backend

import (
	"github.com/kohirens/go-login"
	"github.com/kohirens/sso/oidc"
	"github.com/kohirens/www/storage"
)

type AccountManager struct {
	store storage.Storage
}

// AddWithProvider Make a new account using an OIDC provider.
func (am *AccountManager) addWithProvider(provider oidc.Provider) (*login.Account, error) {
	return login.NewAccountByProvider(provider, am.store)
}

// RetrieveAnAccount finds an existing account or returns a new account.
func (am *AccountManager) RetrieveAnAccount(accountLinkId string) (*login.Account, error) {
	return login.FindAccount(accountLinkId, am.store)
}

func NewAccountManager(store storage.Storage) *AccountManager {
	return &AccountManager{store: store}
}
