package backend

import (
	"net/http"

	"github.com/kohirens/storage"
	"github.com/kohirens/www/awslambda"
	"github.com/kohirens/www/session"
)

type App interface {
	AddRoute(endpoint string, handler Route)
	AddService(key string, service interface{})
	ProviderManager() *ProviderManager
	Decrypt(message []byte) ([]byte, error)
	Encrypt(message []byte) ([]byte, error)
	LoadGPG()
	Name() string
	RouteNotFound(handler Route)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	ServeLambda(event *awslambda.Input) (*awslambda.Output, error)
	ServiceManager() ServiceManager
	TmplManager() TemplateManager
	SessionManager() (*session.Manager, error)
}

// New backend initialized instance.
func New(
	name string,
	router RouteManager,
	serviceManager ServiceManager,
	tmpl TemplateManager,
	pm *ProviderManager,
	store storage.Storage,
) App {
	return &Api{
		accountManager:  NewAccountManager(store),
		name:            name,
		serviceManager:  serviceManager,
		router:          router,
		tmplManager:     tmpl,
		providerManager: pm,
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
