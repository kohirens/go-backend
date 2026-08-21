package backend

import (
	"testing"

	"github.com/kohirens/sso/oidc"
	"github.com/kohirens/www/storage"
)

func TestAccountManager_RetrieveAnAccount(t *testing.T) {
	fixedStore, _ := storage.NewLocalStorage(fixtureDir)

	tests := []struct {
		name          string
		store         storage.Storage
		accountLinkID string
		id            string
		wantErr       bool
	}{
		{
			"pull_account",
			fixedStore,
			"0001",
			"1234",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := NewAccountManager(tt.store)

			got, err := am.RetrieveAnAccount(tt.accountLinkID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAnAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.ID() != tt.id {
				t.Errorf("RetrieveAnAccount() got = %v, want %v", got, tt.id)
				return
			}
		})
	}
}

func TestAccountManager_addWithProvider(t *testing.T) {
	fixedStore, _ := storage.NewLocalStorage(tmpDir)

	tests := []struct {
		name     string
		provider oidc.Provider
		store    storage.Storage
		wantErr  bool
	}{
		{
			"success",
			&MockProvider{},
			fixedStore,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &AccountManager{
				store: tt.store,
			}
			got, err := am.addWithProvider(tt.provider)

			if (err != nil) != tt.wantErr {
				t.Errorf("addWithProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.ID() != tt.provider.ClientID() {
				t.Errorf("addWithProvider() got = %v, want %v", got.ID(), tt.provider.ClientID())
				return
			}

			got2, err2 := am.RetrieveAnAccount(got.ID())
			if (err2 != nil) != tt.wantErr {
				t.Errorf("addWithProvider() error = %v, wantErr %v", err2, tt.wantErr)
				return
			}
			if got2.ID() != got.ID() {
				t.Errorf("addWithProvider() error, newly added account could not be loaded")
				return
			}
		})
	}
}
