package backend

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRouter_Route(t *testing.T) {
	fix := func(w http.ResponseWriter, r *http.Request, a App) error { return nil }
	tests := []struct {
		name     string
		handler  Route
		endpoint string
		want     Route
	}{
		{
			name:     "found",
			handler:  fix,
			endpoint: "*.html",
			want:     fix,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := &Router{
				routes:          make(map[string]Route),
				notFoundHandler: func(w http.ResponseWriter, r *http.Request, a App) error { return nil },
			}

			router.Add(tt.endpoint, tt.handler)
			got := router.Find(tt.endpoint)

			gotRef := fmt.Sprintf("%v", got)
			wantRef := fmt.Sprintf("%v", tt.want)
			if gotRef != wantRef {
				t.Errorf("Find() = %v, want %v", gotRef, wantRef)
				return
			}
		})
	}
}
