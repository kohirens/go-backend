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
	loadRoutes(app)

	mainErr = http.ListenAndServeTLS(":443", certFile, certKey, app)
}

// Health check response.
func Health(w http.ResponseWriter, _ *http.Request,_ backend.App) {
	if _, e := w.Write([]byte("OK")); e != nil {
		log.Errf("internal error %v", e.Error())
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

// HomePage the index page.
func HomePage(w http.ResponseWriter, _ *http.Request, app backend.App) {
    _, e1 := app.TmplManager().RenderFiles(w, map[string]any{}, "layout.html", "index.html")
	if e1 != nil {
		backend.HandleError(e1, w)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func loadRoutes(app backend.App) {
	// This seems like u·ro·bo·ros.
	app.AddRoute("/health", Health)
	app.AddRoute("/", HomePage)
}
```