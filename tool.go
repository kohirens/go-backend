package backend

import (
	"net/http"

	"github.com/kohirens/go-login"
	"github.com/kohirens/www/storage"
)

func getClientApp(store storage.Storage, r *http.Request, app App) (*login.ClientApp, error) {
	ec, e1 := DecryptCookie(EncryptedCookieName, r, app)
	if e1 != nil {
		return nil, e1
	}

	if ec != nil {
		// TODO: Convert encrypted cookie to the proper object.
		clientApp, err := login.LoadClientApp(string(ec.Value), store)
		if err != nil {
			Log.Errf("%v", err.Error())
		}
		if clientApp != nil {
			return clientApp, nil
		}
	}

	return nil, nil
}

func HandleError(err error, w http.ResponseWriter) {
	switch e := err.(type) {
	case *ReferralError:
		if e.Location != "" {
			w.Header().Set("Location", e.Location)
		}
		if e.Body != nil {
			_, eX := w.Write(e.Body)
			if eX != nil {
				Log.Errf(stderr.WriteResponse, eX.Error())
			}
		}
		w.Header().Set("Content-Type", e.ContentType)
		w.WriteHeader(e.Code)

		if e.Log {
			Log.Errf("%v", e.Error())
		}
	default:
		Log.Errf("%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
