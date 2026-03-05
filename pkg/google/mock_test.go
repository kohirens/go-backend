package google

import (
	"net/http"

	"github.com/kohirens/go-backend"
	"github.com/kohirens/www/awslambda"
)

type MockApp struct {
	Authorizer backend.AuthManager
	name       string
}

func (m *MockApp) LoadGPG() {
	//TODO implement me
	panic("implement me")
}

func (m *MockApp) Decrypt(message []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockApp) Encrypt(message []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockApp) AddRoute(endpoint string, handler backend.Route) {
	//TODO implement me
	panic("implement me")
}

func (m *MockApp) AddService(key string, service interface{}) {
	//TODO implement me
	panic("implement me")
}

func (m *MockApp) AuthManager() backend.AuthManager {
	return m.Authorizer
}

func (m *MockApp) Name() string {
	return m.name
}

func (m *MockApp) ServiceManager() backend.ServiceManager {
	//TODO implement me
	panic("implement me")
}

func (m *MockApp) TmplManager() backend.TemplateManager {
	//TODO implement me
	panic("implement me")
}

func (m *MockApp) RouteNotFound(handler backend.Route) {
	//TODO implement me
	panic("implement me")
}

func (m *MockApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (m *MockApp) ServeLambda(event *awslambda.Input) (*awslambda.Output, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockApp) Service(key string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

type MockProvider struct {
	ExpectedAuthLink      string
	ExpectedApp           string
	ExpectedClientID      string
	ExpectedEmail         string
	ExpectedName          string
	ExpectedAuthLinkError error
}

func (m *MockProvider) AuthLink(loginHint string) (string, error) {
	return m.ExpectedAuthLink, m.ExpectedAuthLinkError
}

func (m *MockProvider) Name() string {
	return m.ExpectedName
}

func (m *MockProvider) Application() string {
	return m.ExpectedApp
}

func (m *MockProvider) ClientEmail() string {
	return m.ExpectedEmail

}

func (m *MockProvider) ClientID() string {
	return m.ExpectedClientID
}

func (m *MockProvider) SignOut() error {
	//TODO implement me
	return nil
}
