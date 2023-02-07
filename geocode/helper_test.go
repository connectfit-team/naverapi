package geocode_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/connectfit-team/naverapi/geocode"
)

func checkURLQueryValue(t *testing.T, r *http.Request, key string, expected string) {
	actual := r.URL.Query().Get(key)
	if actual != expected {
		t.Errorf("Expected %s for key %s but got %s", expected, key, actual)
	}
}

func checkCount(t *testing.T, r *http.Request, expected string) {
	checkURLQueryValue(t, r, "count", expected)
}

func checkPage(t *testing.T, r *http.Request, expected string) {
	checkURLQueryValue(t, r, "page", expected)
}

func checkCode(t *testing.T, r *http.Request, expected string) {
	checkURLQueryValue(t, r, "filter", expected)
}

func checkLanguage(t *testing.T, r *http.Request, expected string) {
	checkURLQueryValue(t, r, "language", expected)
}

func checkCoordinate(t *testing.T, r *http.Request, expected string) {
	checkURLQueryValue(t, r, "coordinate", expected)
}

func checkQuery(t *testing.T, r *http.Request, expected string) {
	checkURLQueryValue(t, r, "query", expected)
}

func setupTestClient() (client *geocode.Client, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()

	srv := httptest.NewServer(mux)

	srvURL, _ := url.Parse(srv.URL)
	client, _ = geocode.NewClient(
		"test-client-id",
		"test-client-secret",
		srv.Client(),
	)
	client.BaseURL = srvURL

	return client, mux, srv.Close
}
