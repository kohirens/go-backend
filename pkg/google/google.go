package google

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kohirens/go-backend"
	"github.com/kohirens/sso"
	"github.com/kohirens/sso/pkg/google"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/session"
	"github.com/kohirens/www/validation"
)

// Legend:
// * f - field
// * max - maximum
const (
	// Used as a hint when the user attempts to login with the provider.
	fEmail = "email"
	fCode  = "code"
	fState = "state"
	name   = "google"
)

var (
	// LoginRedirect A location the client will be sent after a successful callback.
	LoginRedirect = "/"
	// Log Set a logger, must be compatible with Kohirens stdlib/logger.
	Log = &logger.Standard{}
	// SignOutRedirect A location to send the client after they sign out.
	SignOutRedirect = "/"
)

// AuthLink Build link to authenticate with Google.
func AuthLink(w http.ResponseWriter, r *http.Request, app backend.App) {
	email, emailOK := validation.Email(r.URL.Query().Get(fEmail))
	if !emailOK {
		email = "" // It's not required, so it is O.K. to leave it out.
	}

	p, e1 := app.AuthManager().Get(backend.KeyGoogleProvider)
	if e1 != nil {
		backend.HandleError(e1, w)
		return
	}
	gp := p.(sso.OIDCProvider)

	authURI, e2 := gp.AuthLink(email)
	if e2 != nil {
		backend.HandleError(e2, w)
		return
	}

	s := fmt.Sprintf(`{"status": %q, "link": %q}`, "ok", authURI)

	_, e3 := w.Write([]byte(s))
	if e3 != nil {
		backend.HandleError(fmt.Errorf(stderr.EncodeJSON, e3.Error()), w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// SignIn Begin the authentication process for a client.
func SignIn(w http.ResponseWriter, r *http.Request, app backend.App) {
	if e := r.ParseForm(); e != nil {
		backend.HandleError(fmt.Errorf(stderr.ParseSignInData, e.Error()), w)
		return
	}

	email := r.PostForm.Get(fEmail)
	_, emailOK := validation.Email(email)
	if email != "" && !emailOK {
		backend.HandleError(backend.NewReferralError(
			"",
			stderr.ValidEmail,
			"/?m=invalid-email",
			http.StatusTemporaryRedirect,
			true,
		), w)
		return
	}

	p, e1 := app.AuthManager().Get(backend.KeyGoogleProvider)
	if e1 != nil {
		backend.HandleError(e1, w)
		return
	}
	gp := p.(sso.OIDCProvider)

	authURI, e2 := gp.AuthLink(email)
	if e2 != nil {
		backend.HandleError(e2, w)
		return
	}

	// set a redirect for the browser.
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Location", authURI)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// SignOut Invalidate a authentication token.
func SignOut(w http.ResponseWriter, _ *http.Request, app backend.App) {
	endpoint := "/?signed-out=1"

	p, e2 := app.AuthManager().Get(backend.KeyGoogleProvider)
	if e2 != nil {
		backend.HandleError(e2, w)
		return
	}
	gp := p.(sso.OIDCProvider)

	if e := gp.SignOut(); e != nil {
		Log.Errf(stderr.SignOut, e)
	}

	body := []byte(fmt.Sprintf(backend.MetaRefresh, endpoint))
	_, e3 := w.Write(body)
	if e3 != nil {
		backend.HandleError(fmt.Errorf(stderr.WriteResponseBody, e3.Error()), w)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Location", SignOutRedirect)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Callback Handles callback request initiated from a Google
// authentication server when the client chose to sign in with Google.
func Callback(w http.ResponseWriter, r *http.Request, app backend.App) {
	Log.Dbugf("%v", stdout.GoogleCallback)

	queryParams := r.URL.Query()
	code := queryParams.Get(fCode)
	state := queryParams.Get(fState)

	gpX, e1 := app.AuthManager().Get(backend.KeyGoogleProvider)
	if e1 != nil {
		backend.HandleError(e1, w)
		return
	}
	gp := gpX.(*google.Provider)

	// Exchange the 1 time code for an ID and refresh tokens.
	if e2 := gp.ExchangeCodeForToken(state, code); e2 != nil {
		backend.HandleError(e2, w)
		return
	}

	// Get client account info
	amX, e3 := app.ServiceManager().Get(backend.KeyAccountManager)
	if e3 != nil {
		backend.HandleError(e3, w)
		return
	}
	am := amX.(backend.AccountManager)

	// Retrieve the session manager.
	smX, e4 := app.Service(backend.KeySessionManager)
	if e4 != nil {
		backend.HandleError(e4, w)
		return
	}
	sm := smX.(*session.Manager)

	// Get user agent data.
	userAgent := r.Header.Get("User-Agent")
	Log.Infof(stdout.UserAgent, userAgent)
	sessionID := sm.ID()
	Log.Infof(stdout.SessionID, sessionID)

	Log.Infof("%v", stdout.EncryptedCookie)
	ec, e12 := GetEncryptedCookie(r, app)
	if e12 != nil {
		Log.Warnf("%v", e12.Error())
	}

	var account *backend.Account
	var loginInfo *sso.LoginInfo
	var makeNewAccount bool
	// Do you have a cookie or not?
	if ec != nil {
		// There MUST be an account ID and a device ID tied to the login,
		// so assume they are validate them before use.
		loginInfo, account = YesCookie(ec, am, gp, sessionID, userAgent)
	} else {
		var e error
		loginInfo, account, e = NoCookie(am, gp, sessionID, userAgent)
		var err *google.ErrNoLoginInfo
		if errors.As(e, &err) {
			makeNewAccount = true
		}
	}

	// If you have no login info, then you should never have an account,
	// the account is only made during login, and it serves as a way to tie
	// multiple providers to a single account.
	// When you're logged in on a different device, but then later use another
	// device but choose a different provider, then this will cause a new
	// account to be made for you. The solution is to log in to eiter account
	// and invite that other account to be merged.
	if ec == nil && loginInfo == nil || makeNewAccount {
		Log.Infof("%v", stdout.MakeAccount)
		var e error
		account, e = registerNewAccount(am, gp)
		if e != nil {
			// TODO: Send them to a page that states: "Something went wrong, please try again later"
			// TODO: This should be custom to the app calling it, so allow the developer to set where
			// TODO: the client will be sent.
			// TODO: Set a temporary redirect.
			panic("something has gone wrong, please try again later")
		}
		if loginInfo == nil {
			Log.Infof(stdout.MakeLoginInfo, gp.Name())

			li, ex := gp.RegisterLoginInfo(account.ID, sessionID, userAgent)
			if ex != nil {
				panic("something has gone wrong, please try again later")
			}
			loginInfo = li
		}
	}

	deviceID := gp.DeviceID()
	if ec == nil && loginInfo != nil {
		Log.Infof("%v", stdout.AddDevice)
		// Since consent has just be granted, add this device,
		// overwriting if it exists. This is OK
		newDevice := sso.NewDevice(userAgent, sessionID, gp.Name())
		loginInfo.Devices[newDevice.ID] = newDevice
		deviceID = newDevice.ID
	}

	Log.Dbugf(stdout.DeviceID, deviceID)
	Log.Dbugf(stdout.AccountID, account.ID)

	Log.Infof("%v", stdout.UpdateLoginInfo)
	if e := gp.UpdateLoginInfo(deviceID, sessionID, userAgent); e != nil {
		backend.HandleError(e, w)
		return
	}

	Log.Infof("%v", stdout.EncryptedCookieValue)
	if e := SetEncryptedCookie(account.ID, deviceID, userAgent, w, app); e != nil {
		backend.HandleError(e, w)
		return
	}

	// send user to a predetermined link or the dashboard.
	w.Header().Set("Location", LoginRedirect)
	w.WriteHeader(http.StatusSeeOther)
}

// RegisterNewAccount Make a new account only when a client has a successful
// login.
func registerNewAccount(
	am backend.AccountManager,
	gp *google.Provider,
) (*backend.Account, error) {
	Log.Dbugf("%v", stdout.RegisterAccount)

	account, e1 := am.AddWithProvider(gp.ClientID(), gp.Name())
	if e1 != nil {
		return nil, e1
	}
	Log.Dbugf(stdout.NewAccount, account.ID)

	// TODO: Change this to gp.Profile() which will have client ID, email address, first, and last name.
	account.GoogleId = gp.ClientID()
	account.Email = gp.ClientEmail()

	return account, nil
}

func NoCookie(
	am backend.AccountManager,
	gp *google.Provider,
	sessionID,
	userAgent string,
) (*sso.LoginInfo, *backend.Account, error) {
	Log.Infof("%v", stdout.LookupLoginInfo)

	li, e1 := gp.LoadLoginInfo("", sessionID, userAgent)
	if e1 != nil {
		return nil, nil, e1
	}

	Log.Infof(stdout.LookupAccount, li.AccountID)

	// if login found, use it to find the linked account.
	account, e2 := am.Lookup(li.AccountID)
	if e2 != nil {
		// Something is really wrong if you have login information, but cannot
		// find the account.
		// You must eject the user, maybe even delete the login info.
		panic(stderr.GetAccount)
	}

	return li, account, nil
}

func YesCookie(
	ec *EncryptedCookie,
	am backend.AccountManager,
	gp *google.Provider,
	sessionID,
	userAgent string,
) (*sso.LoginInfo, *backend.Account) {
	loginInfo, e1 := gp.LoadLoginInfo(ec.DID, sessionID, userAgent)
	if loginInfo == nil || e1 != nil {
		// something is really strange if you have a cookie but cannot retrieve
		// the loginInfo.
		// NOTE: Its OK if the device cannot be found, it could have been
		// manually deleted.
		Log.Errf("%v", e1.Error())
		// TODO Figure out how to proceed.
		panic(stderr.LookupLoginInfo)
	}

	// Get the account.
	account, e2 := am.Lookup(ec.AID)
	if e2 != nil {
		// something is really strange if you have an account number but cannot
		// find it.
		Log.Errf("%v", e2.Error())
		// TODO Figure out how to proceed.
		panic(stderr.GetAccount)
	}

	return loginInfo, account
}
