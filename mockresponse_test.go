package backend

import (
	"fmt"
	"net/http"
)

type MockResponse struct {
	ExpectedBody       []byte
	ExpectedHeaders    http.Header
	Headers            http.Header
	ExpectedStatusCode int
}

func (m MockResponse) Header() http.Header {
	if m.Headers == nil {
		m.Headers = make(http.Header)
	}
	return m.Headers
}

func (m MockResponse) Write(bytes []byte) (int, error) {
	idx := 0
	if m.ExpectedBody != nil {
		var val byte
		for idx, val = range bytes {
			if m.ExpectedBody[idx] != val {
				panic("unexpected bytes found in the body")
			}
		}
	}
	return idx, nil
}

func (m MockResponse) WriteHeader(statusCode int) {
	if statusCode != m.ExpectedStatusCode {
		panic(fmt.Sprintf("unexpected status code %v", statusCode))
	}
}
