package backend

import "net/http"

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
