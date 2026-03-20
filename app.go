package backend

import (
	"net/http"

	"github.com/kohirens/www/awslambda"
	"github.com/kohirens/www/session"
	"github.com/kohirens/www/storage"
)

type App interface {
	AddRoute(endpoint string, handler Route)
	AddService(key string, service interface{})
	ProviderManager() ProviderManager
	Provider(string) any
	Decrypt(message []byte) ([]byte, error)
	Encrypt(message []byte) ([]byte, error)
	LoadGPG()
	Name() string
	RouteNotFound(handler Route)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	ServeLambda(event *awslambda.Input) (*awslambda.Output, error)
	Service(key string) (interface{}, error)
	ServiceManager() ServiceManager
	TmplManager() TemplateManager
	Session() (*session.Manager, error)
}

// New A nNew initialized application instance.
func New(
	name string,
	router RouteManager,
	serviceManager ServiceManager,
	tmpl TemplateManager,
	authManager ProviderManager,
	store storage.Storage,
) App {
	return &Api{
		name:            name,
		serviceManager:  serviceManager,
		router:          router,
		tmplManager:     tmpl,
		providerManager: authManager,
		storage:         store,
	}
}

// NewWithDefaults initialize a new backend application. The name MUST match
// the filename of GPG key stored in JSON format and located in /secrets of the
// storage.
func NewWithDefaults(name string, store storage.Storage) App {
	return New(
		name,
		NewRouteManager(),
		NewServiceManager(),
		NewTemplateManager(store, TmplDir, TmplSuffix),
		NewProviderManager(),
		store,
	)
}
