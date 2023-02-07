package mailer_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/connectfit-team/naverapi/internal/testhelper"
	"github.com/connectfit-team/naverapi/mailer"
)

const (
	testRequestTimestamp = "856915200000"
	testRequestAccessKey = "test-access-key"
)

const (
	createMailRequestAPIGatewaySignature = "F1YxxwEjDRZmNLxqqDFz53OpbvLrMCqEsv9tLxoBcWE="

	createMailRequestBody = `{"senderAddress":"test-sender-address","senderName":"test-sender-name","title":"test-title","body":"test-body","recipients":[{"address":"test-address","name":"test-name","type":"test-type","parameters":["test-parameter-1","test-parameter-2"]}],"attachFileIds":["test-file-id-1","test-file-id-2"]}`
)

const (
	createFilesRequestAPIGatewaySignature = "q1JhbYivx0lU//wBoOyh+yn/y7+Lg9Ez/Xj6FDzxap4="
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

func checkCreateMailRequest(t *testing.T, r *http.Request) {
	testhelper.TestRequestMethod(t, r, http.MethodPost)

	testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", testRequestTimestamp)
	testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", testRequestAccessKey)
	testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", createMailRequestAPIGatewaySignature)
	testhelper.TestRequestHeader(t, r, "Content-Type", "application/json")

	testhelper.TestRequestBody(t, r, createMailRequestBody)
}

func checkCreateFilesRequest(t *testing.T, r *http.Request) {
	testhelper.TestRequestMethod(t, r, http.MethodPost)

	testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", testRequestTimestamp)
	testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", testRequestAccessKey)
	testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", createFilesRequestAPIGatewaySignature)

	testhelper.TestRequestFormFiles(t, r, "fileList", []string{"test-content-1", "test-content-2"})
}
