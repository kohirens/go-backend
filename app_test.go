package backend

import (
	"reflect"
	"testing"

	"github.com/kohirens/sso"
	"github.com/kohirens/www/storage"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		provider sso.OIDCProvider
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
			got := fixture.Provider("gp")
			gotType := reflect.TypeOf(got)

			if gotType != tt.want {
				t.Errorf("New() = %v, want %T", gotType.Name(), tt.want.Name())
			}
		})
	}
}

type MockProvider struct {
	m map[string]sso.OIDCProvider
}

func (mp *MockProvider) Add(name string, provider sso.OIDCProvider) {
	mp.m[name] = provider
}

func (m *MockProvider) Get(name string) (sso.OIDCProvider, error) {
	return m.m[name], nil
}
func (p *MockProvider) AuthLink(loginHint string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p *MockProvider) Name() string {
	//TODO implement me
	panic("implement me")
}

func (p *MockProvider) Application() string {
	//TODO implement me
	panic("implement me")
}

func (p *MockProvider) ClientEmail() string {
	//TODO implement me
	panic("implement me")
}

func (p *MockProvider) ClientID() string {
	//TODO implement me
	panic("implement me")
}

func (p *MockProvider) SignOut() error {
	//TODO implement me
	panic("implement me")
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
