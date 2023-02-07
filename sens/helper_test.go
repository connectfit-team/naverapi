package sens_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/connectfit-team/naverapi/sens"
)

type fixedTimeClock struct {
	fixedTime time.Time
}

func (ftc fixedTimeClock) Now() time.Time { return ftc.fixedTime }

func setupTestSENSClient() (client *sens.Client, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()

	srv := httptest.NewServer(mux)

	srvURL, _ := url.Parse(srv.URL)
	client, _ = sens.NewClient(
		"test-access-key",
		"test-secret-key",
		srv.Client(),
	)
	client.Clock = &fixedTimeClock{
		fixedTime: time.Date(1997, 02, 26, 0, 0, 0, 0, time.UTC),
	}
	client.BaseURL = srvURL

	return client, mux, srv.Close
}
