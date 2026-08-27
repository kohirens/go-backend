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

## Login Flow

1. Load a GPG key into the keychain
    1. Use GOpenGPG library.
    2. Pull the GPG key from storage.
    3. Load the GPG keys from file.
2. Store data encrypted with the webapp GPG key into a cookie.
    1. When the user has authenticated, but before you return a response:
        1. Encrypt some info about their login status with the GPG key.
        2. Save this encrypted value in a secure cookie.
3. Decrypt data with the webapp GPG key.
    1. Check if a user is logged in:
        1. Look for a secure cookie.
        2. Take the value and try to decrypt it with the webapp GPG key.
        3. If there is a valid info, then:
            1. pull the account based on user validated info.
            2. Get the device ID.
            3. Use the device ID and search for it in the clients account,
               if there is a match, then:
                1. Pull the provider the logged in with on the device.
                2. Pull the provider login info from storage, if found, then
                    1. Check to see if authentication token has expired:
                        1. If not, then restore it.
                        2. If yes, then re-authenticate or get a fresh token.
