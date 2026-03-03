# Go Backend

Keeps it simple with developing web applications with just the tools you need.

```go
package main

import (
	"net/http"

	"github.com/kohirens/go-backend"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/storage"
)

const (
	assetsFilesDir   = "../../frontend/assets"
	certFile         = "/root/pki/certs/server.crt"
	certKey          = "/root/pki/private/server.key"
	templateFilesDir = "../../templates"
)

var (
	log = &logger.Standard{}
)

func main() {
	var mainErr error

	defer func() {
		if mainErr != nil {
			log.Errf("main error: %v", mainErr)
		}
	}()

	logger.VerbosityLevel = 6

	// Initialize the backend API storage.
	// Initialize a storage handler for the backend.
	store, e2 := storage.NewLocalStorage("./")
	if e2 != nil {
		mainErr = e2
		return
	}

	// Initialize the backend API.
	app := backend.NewWithDefaults("webapp", store)
	
	// Add all the routes you want.
	loadRoutes(app, &Responder{})

	mainErr = http.ListenAndServeTLS(":443", certFile, certKey, app)
}


type Responder struct {}
// Health check response.
func (s *Responder) Health(w http.ResponseWriter, r *http.Request) {
	if _, e := w.Write([]byte("OK")); e != nil {
		log.Errf("internal error %v", e.Error())
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func loadRoutes(app backend.App, responder *Responder) {
	app.AddRoute("/health", responder.Health)
}
```