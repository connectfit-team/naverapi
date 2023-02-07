package mailer_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/connectfit-team/naverapi/mailer"
)

type fixedTimeClock struct {
	fixedTime time.Time
}

func (ftc fixedTimeClock) Now() time.Time { return ftc.fixedTime }

func setupTestCloudOutboundMailerClient() (client *mailer.CloudOutboundMailerClient, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()

	srv := httptest.NewServer(mux)

	srvURL, _ := url.Parse(srv.URL)
	client, _ = mailer.NewCloudOutboundMailerClient(
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
