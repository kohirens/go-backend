package google

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kohirens/go-backend"
	"github.com/kohirens/go-login"
	"github.com/kohirens/sso"
	"github.com/kohirens/sso/pkg/google"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/validation"
)

// Legend:
// * f - field
// * max - maximum
const (
	// Used as a hint when the user attempts to login with the provider.
	fEmail            = "email"
	fCode             = "code"
	fState            = "state"
	KeyGoogleProvider = "gp"
	name              = "google"
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

	p, e1 := app.ProviderManager().Get(KeyGoogleProvider)
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

	p, e1 := app.ProviderManager().Get(KeyGoogleProvider)
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

	p, e2 := app.ProviderManager().Get(KeyGoogleProvider)
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

	gpX, e1 := app.ProviderManager().Get(KeyGoogleProvider)
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

	// Load the client's profile.

	// Get and decrypt the clientApp data.
	Log.Infof("%v", stdout.EncryptedCookie)
	ec, e3 := backend.DecryptCookie(backend.EncryptedCookieName, r, app)
	if e3 != nil {
		Log.Warnf("%v", e3.Error())
	}

	var clientApp *login.ClientApp
	if ec != nil {
		var err error
		clientApp, err = login.LoadClientApp(ec.Value)
		if err != nil {
			Log.Warnf("%v", err.Error())
		}
	} else {
		clientApp = login.NewClientApp()
	}

	Log.Infof(stdout.ClientApp, clientApp.Id)
	Log.Infof("%v", clientApp.LastDate.UTC().Format(time.RFC3339))

	// Get the storage manager so we can pull the account and profile.
	//storeX, e4 := app.ServiceManager().Get("store")
	//if e4 != nil {
	//	backend.HandleError(e4, w)
	//	return
	//}
	//store := storeX.(storage.Storage)

	//var account *login.Account
	//
	//// Lookup the account in the cookie.
	//if clientApp.AccountId != "" {
	//	var err error
	//	account, err = login.LoadAccount(clientApp.AccountId, store)
	//	if err != nil {
	//		backend.HandleError(err, w)
	//		return
	//	}
	//}
	//
	//// Make a new account if one does not exist.
	//if account == nil {
	//	login.NewAccount()
	//	// TODO: Pull the profile.
	//	// TODO: Update the cookie.
	//	// TODO: Lookup the login.
	//	profileId, e5 := login.LoadProfileMap(gp.ClientID(), store)
	//	if e5 != nil {
	//		backend.HandleError(e5, w)
	//		return
	//	}
	//	// TODO: Lookup the profile.
	//	profile, e6 := login.LoadProfile(profileId, store)
	//	if e6 != nil {
	//		backend.HandleError(e6, w)
	//		return
	//	}
	//	Log.Dbugf("profile name: %v", profile.Name)
	//}
	//
	//// Retrieve the session manager.
	//smX, e4 := app.Service(backend.KeySessionManager)
	//if e4 != nil {
	//	backend.HandleError(e4, w)
	//	return
	//}
	//sm := smX.(*session.Manager)
	//
	//// Get user agent data.
	//userAgent := r.Header.Get("User-Agent")
	//Log.Infof(stdout.UserAgent, userAgent)
	//sessionID := sm.ID()
	//Log.Infof(stdout.SessionID, sessionID)
	//
	//var account *backend.Account
	//// If you have no login info, then you should never have an account,
	//// the account is only made during login, and it serves as a way to tie
	//// multiple providers to a single account.
	//// When you're logged in on a different device, but then later use another
	//// device but choose a different provider, then this will cause a new
	//// account to be made for you. The solution is to log in to eiter account
	//// and invite that other account to be merged.
	//if ec == nil && loginInfo == nil || makeNewAccount {
	//	Log.Infof("%v", stdout.MakeAccount)
	//	var e error
	//	account, e = registerNewAccount(am, gp)
	//	if e != nil {
	//		// TODO: Send them to a page that states: "Something went wrong, please try again later"
	//		// TODO: This should be custom to the app calling it, so allow the developer to set where
	//		// TODO: the client will be sent.
	//		// TODO: Set a temporary redirect.
	//		panic("something has gone wrong, please try again later")
	//	}
	//	if loginInfo == nil {
	//		Log.Infof(stdout.MakeLoginInfo, gp.Name())
	//
	//		li, ex := gp.RegisterLoginInfo(account.ID, sessionID, userAgent)
	//		if ex != nil {
	//			panic("something has gone wrong, please try again later")
	//		}
	//		loginInfo = li
	//	}
	//}
	//
	//Log.Dbugf(stdout.AccountID, account.ID)
	//
	//Log.Infof("%v", stdout.EncryptedCookieValue)
	//if e := backend.EncryptCookie(backend.EncryptedCookieName, "", userAgent, w, app); e != nil {
	//	backend.HandleError(e, w)
	//	return
	//}

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
