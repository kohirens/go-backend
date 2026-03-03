package backend

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/kohirens/www/storage"
)

type Account struct {
	AppleID  string `json:"apple_id"`
	Email    string `json:"email"`
	FistName string `json:"fist_name"`
	GoogleId string `json:"google_id"`
	ID       string `json:"id"`
	LastName string `json:"last_name"`
}

type AccountManager interface {
	AddWithProvider(providerID, providerName string) (*Account, error)
	Lookup(id string) (*Account, error)
	// Location The generated filename/key/path where the account
	// should be located in storage.
	Location(id string) string
}

// AccountExec Short for account executive, is an implementation of
// AccountManager.
type AccountExec struct {
	store storage.Storage
}

// AddWithProvider Make a new account using an OIDC provider.
func (am *AccountExec) AddWithProvider(clientID, providerName string) (*Account, error) {
	// Generate an account ID.
	id, e1 := uuid.NewV7()
	if e1 != nil {
		return nil, fmt.Errorf(stderr.UUID, e1.Error())
	}

	account := &Account{
		ID: id.String(), //TODO: generate a guid
	}

	switch providerName {
	case "apple":
		account.AppleID = clientID
	case "google":
		account.GoogleId = clientID
	}

	accountBytes, e2 := json.Marshal(account)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.DecodeJSON, e2.Error())
	}

	if e := am.store.Save(am.Location(account.ID), accountBytes); e != nil {
		return nil, e
	}

	return account, nil
}

func (am *AccountExec) Location(id string) string {
	return fmt.Sprintf(PrefixAccounts+"/%v.json", id)
}

// Lookup Search for an account in storage.
func (am *AccountExec) Lookup(id string) (*Account, error) {
	aData, e1 := am.store.Load(am.Location(id))
	if e1 != nil {
		return nil, &AccountNotFoundError{id}
	}

	account := &Account{}
	if e2 := json.Unmarshal(aData, &account); e2 != nil {
		return nil, fmt.Errorf(stderr.DecodeJSON, e2.Error())
	}

	return account, nil
}
