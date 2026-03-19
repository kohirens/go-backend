package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kohirens/sso"
	"github.com/kohirens/www/awslambda"
	"github.com/kohirens/www/gpg"
	"github.com/kohirens/www/session"
	"github.com/kohirens/www/storage"
)

// Api serves as the backend server for managing routes (a.k.a endpoints),
// services, authentication providers, and a template engine.
// These components are available to your routes (which you define). Also,
// these components are replaceable as long as the meet the interface
// requirements.
type Api struct {
	providerManager ProviderManager
	capsule         *gpg.Capsule
	gpgKey          *appKey
	name            string
	router          RouteManager
	serviceManager  ServiceManager
	storage         storage.Storage
	tmplManager     TemplateManager
}

// AddProvider Wrapper method that adds an auth provider to the ProviderManager
// for retrieval during request handling.
func (a *Api) AddProvider(key string, provider sso.OIDCProvider) {
	a.providerManager.Add(key, provider)
}

// AddRoute Maps a function to a http.HandlerFunc so that it will respond when
// the route (a.k.a endpoint) is requested.
func (a *Api) AddRoute(endpoint string, handler Route) {
	a.router.Add(endpoint, handler)
}

func (a *Api) AddService(key string, service interface{}) {
	a.serviceManager.Add(key, service)
}

// Decrypt Decode a message using the apps key.
func (a *Api) Decrypt(message []byte) ([]byte, error) {
	return a.capsule.Decrypt(message)
}

// Encrypt cipher a message using the apps key.
func (a *Api) Encrypt(subject []byte) ([]byte, error) {
	return a.capsule.Encrypt(subject)
}

// LoadGPG Pull the GPG key from <storage>/secret/<app-name>
func (a *Api) LoadGPG() {
	Log.Dbugf("%v", stdout.LoadGPG)

	gpgData, e1 := a.storage.Load(PrefixSecrets + "/" + a.Name() + ".json")
	if e1 != nil {
		panic(e1.Error())
	}

	gpgKey := &appKey{}
	if e := json.Unmarshal(gpgData, &gpgKey); e != nil {
		panic(fmt.Sprintf(stderr.DecodeJSON, e.Error()))
	}

	// Encrypt the data and store in a secure cookie.
	capsule, e9 := gpg.NewCapsuleString(gpgKey.PublicKey, gpgKey.PrivateKey, gpgKey.PassPhrase)
	if e9 != nil {
		panic(e9)
	}

	a.gpgKey = gpgKey
	a.capsule = capsule
}

// Name A name/ID given to the application.
func (a *Api) Name() string {
	return a.name
}

// Provider get an OIDC provider from the manager.
func (a *Api) Provider(authProvider string) any {
	p, e1 := a.providerManager.Get(authProvider)
	if e1 != nil {
		Log.Errf(stderr.AuthProviderLookup, e1.Error())
	}

	return p
}

// ProviderManager Return the authentication manager.
func (a *Api) ProviderManager() ProviderManager {
	return a.providerManager
}

// RouteNotFound Add a http.HandlerFunc to return a response when a route is
// not found.
func (a *Api) RouteNotFound(handler Route) {
	a.router.NotFound(handler)
}

// ServeHTTP Will be called for every request to this server. There is no need
// to register individual handlers for each pattern or use confusing middleware
// logic.
func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	processRequest(w, r, a)
}

// ServeLambda Provide an HTTP response for an AWS Lambda function. A little
// extra because the AWS even it not compatible with http.Request, same for its
// response, which also has special considerations.
func (a *Api) ServeLambda(event *awslambda.Input) (*awslambda.Output, error) {
	Log.Infof("%v", stdout.Started)

	if errRes := awslambda.PreliminaryChecks(event); errRes != nil {
		return errRes, nil
	}

	w := awslambda.NewResponse()

	r, e1 := awslambda.NewRequest(event)
	if e1 != nil {
		Log.Errf("%v", e1.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return w, nil
	}

	processRequest(w, r, a)

	return w, nil
}

// Service returns a service by name.
func (a *Api) Service(key string) (interface{}, error) {
	return a.serviceManager.Get(key)
}

// ServiceManager Get the handler for retrieving services.
func (a *Api) ServiceManager() ServiceManager {
	return a.serviceManager
}

// Session Get the session manager.
func (a *Api) Session() (*session.Manager, error) {
	x, e1 := a.serviceManager.Get(KeySessionManager)
	if e1 != nil {
		return nil, e1
	}
	return x.(*session.Manager), nil
}

// TmplManager Template engine that renders templates.
func (a *Api) TmplManager() TemplateManager {
	return a.tmplManager
}

// processRequest responsibilities:
//  1. Initialize/Load an HTTP session for client requests.
//  2. Load logic to process a request and write a response.
//  3. Save the session before sending an HTTP response.
func processRequest(w http.ResponseWriter, r *http.Request, a *Api) {
	rawPath := r.URL.Path
	Log.Infof("request %v %v", r.Method, rawPath)

	idCookie, _ := r.Cookie(session.IDKey)
	if e := a.RestoreSessionData(w, idCookie); e != nil {
		Log.Errf("%v", e.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	sess, e1 := a.Session()
	if e1 != nil {
		Log.Errf("%v", e1.Error())
	}
	sessionTime := sess.Expiration().UTC().Sub(time.Now().UTC())
	Log.Dbugf(stdout.SessionTime, sessionTime)

	// Add common variables to the template manager.
	a.tmplManager.AppendVars(Variables{
		"HTTP_Method": r.Method,
		"URL_Path":    rawPath,
		"URL_Query":   r.URL.RawQuery,
		"SessionTime": int(sessionTime.Seconds()),
	})

	// Find the route to respond to the request.
	fn := a.router.Find(rawPath)

	fn(w, r, a)

	Log.Infof("%v", stdout.PageDone)

	if e := a.SaveSessionData(w, r); e != nil {
		Log.Errf("%v", e.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *Api) RestoreSessionData(w http.ResponseWriter, idCookie *http.Cookie) error {
	sm, e1 := a.Session()
	if e1 != nil {
		return e1
	}

	sm.Load(w, idCookie)

	// TODO pull from the cookie which provider the client chose.
	gp, e2 := a.providerManager.Get(KeyGoogleProvider)
	if e2 != nil {
		return e2
	}

	gpData := sm.Get(sso.SessionTokenGoogle)
	if gpData != nil { // restore from the saved session.
		// TODO: Test if you can overwrite members of an initialized struct from a json.Unmarshal.
		//var savedGp *sso.GoogleProvider
		if e := json.Unmarshal(gpData, &gp); e != nil {
			var je *json.UnmarshalTypeError
			if errors.As(e, &je) {
				return fmt.Errorf(stderr.UnmarshalJSON, je.Field, je.Value, je.Offset)
			}
			return fmt.Errorf(stderr.DecodeJSON, e.Error())
		}
	}

	return nil
}

func (a *Api) SaveSessionData(w http.ResponseWriter, r *http.Request) error {
	sm, e1 := a.Session()
	if e1 != nil {
		return e1
	}

	authProvider, e2 := a.providerManager.Get(KeyGoogleProvider)
	if authProvider == nil {
		return e2
	}

	gpData, e3 := json.Marshal(authProvider)
	if e3 != nil {
		return e3
	}

	// When you restore the Google provider from the session the previous token
	// should also be restored.
	sm.Set(sso.SessionTokenGoogle, gpData)

	if e := sm.Save(); e != nil {
		return e
	}

	return nil
}

// Storage Retrieve the storage service from the service manager.
func (a *Api) Storage() storage.Storage {
	return a.storage
}

type appKey struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	PassPhrase string `json:"pass_phrase"`
}
