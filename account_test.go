package backend

import (
	"testing"

	"github.com/kohirens/www/storage"
)

func TestAccountExec_Lookup(t *testing.T) {
	fixedStore, _ := storage.NewLocalStorage(fixtureDir)

	tests := []struct {
		name    string
		store   storage.Storage
		id      string
		wantErr bool
	}{
		{
			"pull_account",
			fixedStore,
			"1234",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &AccountExec{
				store: tt.store,
			}
			got, err := am.Lookup(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.ID != tt.id {
				t.Errorf("Lookup() got = %v, want %v", got, tt.id)
				return
			}
		})
	}
}

func TestAccountExec_Add(t *testing.T) {
	fixedStore, _ := storage.NewLocalStorage(tmpDir)

	tests := []struct {
		name,
		providerID,
		providerName string
		store   storage.Storage
		wantErr bool
	}{
		{
			"add_account",
			"1234",
			"google",
			fixedStore,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &AccountExec{
				store: tt.store,
			}
			got, err := am.AddWithProvider(tt.providerID, tt.providerName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.GoogleId != tt.providerID {
				t.Errorf("Add() got = %v, want %v", got.GoogleId, tt.providerID)
				return
			}

			got2, err2 := am.Lookup(got.ID)
			if (err2 != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err2, tt.wantErr)
				return
			}
			if got2.ID != got.ID {
				t.Errorf("Add() error, newly added account could not be loaded")
				return
			}
		})
	}
}
