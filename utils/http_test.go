package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeRequestWithUserAgentHeader(t *testing.T) {
	var userAgent string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent = r.Header.Get("User-Agent")
	}))
	defer ts.Close()

	expectedUserAgent := fmt.Sprintf("gonews:v%s (by /u/drgarcia1986)", Version)
	MakeRequest("GET", ts.URL, nil)
	if userAgent != expectedUserAgent {
		t.Errorf("Expected %s, got %s", expectedUserAgent, userAgent)
	}

}
