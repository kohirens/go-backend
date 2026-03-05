package google

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kohirens/go-backend"
)

type EncryptedCookie struct {
	AID          string
	DID          string
	UA           string
	LastActivity time.Time
}

const ecCookieName = "_ec_"

func GetEncryptedCookie(r *http.Request, a backend.App) (*EncryptedCookie, error) {
	cookie, e1 := r.Cookie(ecCookieName)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.ECCookie, e1.Error())
	}

	ecBytes, e2 := base64.StdEncoding.DecodeString(cookie.Value)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.DecodeBase64, e2.Error())
	}

	message, e3 := a.Decrypt(ecBytes)
	if e3 != nil {
		return nil, e3
	}

	ec := &EncryptedCookie{}
	if e := json.Unmarshal(message, ec); e != nil {
		return nil, fmt.Errorf(stderr.DecodeJSON, e.Error())
	}

	return ec, nil
}

func SetEncryptedCookie(
	accountID,
	deviceID,
	userAgent string,
	w http.ResponseWriter,
	a backend.App,
) error {
	ec := &EncryptedCookie{
		AID:          accountID,
		DID:          deviceID,
		UA:           userAgent,
		LastActivity: time.Now(),
	}

	ecBytes, e15 := json.Marshal(ec)
	if e15 != nil {
		return fmt.Errorf(stderr.EncodeJSON, e15.Error())
	}

	encodeMessage, e10 := a.Encrypt(ecBytes)
	if e10 != nil {
		return e10
	}

	http.SetCookie(w, &http.Cookie{
		Name:   ecCookieName,
		Value:  base64.StdEncoding.EncodeToString(encodeMessage),
		Path:   "/",
		Secure: true,
	})
	return nil
}
