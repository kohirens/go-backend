package google

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/kohirens/go-backend"
	"github.com/kohirens/stdlib/test"
)

func TestAuthLink(t *testing.T) {
	goodAuth := backend.NewProviderManager()
	goodAuth.Add(backend.KeyGoogleProvider, &MockProvider{
		ExpectedAuthLink: "good-link",
	})

	cases := []struct {
		name    string
		w       http.ResponseWriter
		r       *http.Request
		a       backend.App
		wantErr bool
	}{
		{
			"provider_not_found",
			&test.MockResponseWriter{
				ExpectedBody:       nil,
				ExpectedHeaders:    nil,
				Headers:            nil,
				ExpectedStatusCode: 500,
			},
			&http.Request{
				URL: &url.URL{
					Scheme:   "https",
					Host:     "google.com",
					Path:     "/auth/google/callback",
					RawQuery: "email=test@example.com",
				},
			},
			&MockApp{
				Authorizer: backend.NewProviderManager(),
			},
			true,
		},
		{
			"return_link",
			&test.MockResponseWriter{
				Headers:            nil,
				ExpectedStatusCode: 200,
			},
			&http.Request{
				URL: &url.URL{
					RawQuery: "email=test@example.com",
				},
				Header: http.Header{},
			},
			&MockApp{
				Authorizer: goodAuth,
			},
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			AuthLink(c.w, c.r, c.a)
		})
	}
}

func TestSignIn(t *testing.T) {
	goodAuth := backend.NewProviderManager()
	goodAuth.Add(backend.KeyGoogleProvider, &MockProvider{
		ExpectedAuthLink: "good-link",
	})

	cases := []struct {
		name    string
		w       http.ResponseWriter
		r       *http.Request
		a       backend.App
		wantErr bool
	}{
		{
			"provider_not_found",
			&test.MockResponseWriter{
				ExpectedBody:       nil,
				ExpectedHeaders:    nil,
				Headers:            nil,
				ExpectedStatusCode: 500,
			},
			&http.Request{
				URL: &url.URL{
					Scheme:   "https",
					Host:     "google.com",
					Path:     "/auth/google/callback",
					RawQuery: "email=test@example.com",
				},
			},
			&MockApp{
				Authorizer: backend.NewProviderManager(),
			},
			true,
		},
		{
			"redirect_to_google_auth_server",
			&test.MockResponseWriter{
				Headers:            nil,
				ExpectedStatusCode: 307,
			},
			&http.Request{
				URL: &url.URL{
					RawQuery: "email=test@example.com",
				},
				Header: http.Header{},
			},
			&MockApp{
				Authorizer: goodAuth,
			},
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			SignIn(c.w, c.r, c.a)
		})
	}
}

func TestSignOut(t *testing.T) {
	goodAuth := backend.NewProviderManager()
	goodAuth.Add(backend.KeyGoogleProvider, &MockProvider{
		ExpectedAuthLink: "good-link",
	})

	cases := []struct {
		name    string
		w       http.ResponseWriter
		r       *http.Request
		a       backend.App
		wantErr bool
	}{
		{
			"provider_not_found",
			&test.MockResponseWriter{
				ExpectedStatusCode: 500,
			},
			&http.Request{
				URL: &url.URL{
					Scheme:   "https",
					Host:     "google.com",
					Path:     "/auth/google/callback",
					RawQuery: "email=test@example.com",
				},
			},
			&MockApp{
				Authorizer: backend.NewProviderManager(),
			},
			true,
		},
		{
			"redirect_to_google_auth_server",
			&test.MockResponseWriter{
				ExpectedStatusCode: 307,
			},
			nil,
			&MockApp{
				Authorizer: goodAuth,
			},
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			SignOut(c.w, c.r, c.a)
		})
	}
}
