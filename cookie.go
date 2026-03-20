package backend

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const EncryptedCookieName = "_ec_"

type EncryptedCookie struct {
	LastActivity time.Time
	Name         string
	Value        []byte
}

// DecryptCookie using a GPG key to allow working with the cookie data.
//
//	The GPG key has to have been loaded.
func DecryptCookie(name string, r *http.Request, a App) (*EncryptedCookie, error) {
	cookie, e1 := r.Cookie(name)
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

// EncryptCookie encrypts data before storing it in a cookie.
func EncryptCookie(
	name,
	value,
	path string,
	w http.ResponseWriter,
	a App,
) error {
	ec := &EncryptedCookie{
		Name:         name,
		Value:        []byte(value),
		LastActivity: time.Now(),
	}

	ecBytes, e1 := json.Marshal(ec)
	if e1 != nil {
		return fmt.Errorf(stderr.EncodeJSON, e1.Error())
	}

	encodeMessage, e2 := a.Encrypt(ecBytes)
	if e2 != nil {
		return e2
	}

	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  base64.StdEncoding.EncodeToString(encodeMessage),
		Path:   path,
		Secure: true,
	})

	return nil
}
