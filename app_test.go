package backend

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/kohirens/go-login"
	"github.com/kohirens/sso/oidc"
	"github.com/kohirens/storage"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		provider oidc.Provider
		want     reflect.Type
	}{
		{
			"add_provider",
			&MockProvider{},
			reflect.TypeOf(&MockProvider{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := NewWithDefaults("test", nil)

			fixture.ProviderManager().Add("gp", tt.provider)
			got, _ := fixture.ProviderManager().Get("gp")
			gotType := reflect.TypeOf(got)

			if gotType != tt.want {
				t.Errorf("New() = %v, want %T", gotType.Name(), tt.want.Name())
			}
		})
	}
}

type MockProvider struct {
	user                  string
	m                     map[string]oidc.Provider
	ExpectedAuthLink      string
	ExpectedApp           string
	ExpectedClientID      string
	ExpectedEmail         string
	ExpectedName          string
	ExpectedAuthLinkError error
}

func (mp *MockProvider) Callback(params url.Values) error {
	//TODO implement me
	panic("implement me")
}

func (mp *MockProvider) UserInfo() oidc.UserInfo {
	switch mp.user {
	default:
		return login.NewUserInfo(
			"mctest_t@example.com",
			"tester",
			"McTest",
			"555-555-5555",
			"en-US",
		)
	}
}

func (mp *MockProvider) String() string {
	//TODO implement me
	panic("implement me")
}

func (mp *MockProvider) Add(name string, provider oidc.Provider) {
	mp.m[name] = provider
}

func (mp *MockProvider) Get(name string) (oidc.Provider, error) {
	return mp.m[name], nil
}

func xTestNewWithDefaults(t *testing.T) {
	type args struct {
		name  string
		store storage.Storage
	}
	tests := []struct {
		name string
		args args
		want App
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWithDefaults(tt.args.name, tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithDefaults() = %v, want %v", got, tt.want)
			}
		})
	}
}
