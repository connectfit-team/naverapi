package mailer_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/connectfit-team/naverapi/mailer"
	"github.com/google/go-cmp/cmp"
)

const validCreateMailResponse = `
	{
		"requestId": "test-request-id",
		"count": 1
	}
`

const validCreateFilesResponse = `
	{
		"tempRequestId": "test-temp-request-id",
		"files": [
			{
				"fileName": "test-filename-1",
				"fileSize": 42,
				"fileId": "test-file-id-1"
			},
			{
				"fileName": "test-filename-2",
				"fileSize": 1337,
				"fileId": "test-file-id-2"
			}
		]
	}
`

var testCreateMailRequest = mailer.CreateMailRequest{
	SenderAddress: "test-sender-address",
	SenderName:    "test-sender-name",
	Title:         "test-title",
	Body:          "test-body",
	Recipients: []*mailer.Recipient{
		{
			Address:    "test-address",
			Name:       "test-name",
			Type:       "test-type",
			Parameters: []string{"test-parameter-1", "test-parameter-2"},
		},
	},
	AttachFileIDs: []string{"test-file-id-1", "test-file-id-2"},
}

var testFiles = []mailer.File{
	{
		Name:    "test-name-1",
		Content: []byte("test-content-1"),
	},
	{
		Name:    "test-name-2",
		Content: []byte("test-content-2"),
	},
}

func TestCloudOutboundMailerClient_CreateMail(t *testing.T) {
	client, mux, teardown := setupTestCloudOutboundMailerClient()
	defer teardown()

	mux.HandleFunc(mailer.EndpointMails, func(w http.ResponseWriter, r *http.Request) {
		checkCreateMailRequest(t, r)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, validCreateMailResponse)
	})

	got, err := client.CreateMail(context.Background(), testCreateMailRequest)
	if err != nil {
		t.Fatalf("createMail request was given a valid request but failed: %v", err)
	}

	want := mailer.CreateMailResponse{
		RequestID: "test-request-id",
		Count:     1,
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Response differ from the expected one: %s", diff)
	}
}

func TestCloudOutboundMailerClient_CreateMail_ShouldFailIfWrongResponseStatusCode(t *testing.T) {
	client, mux, teardown := setupTestCloudOutboundMailerClient()
	defer teardown()

	mux.HandleFunc(mailer.EndpointMails, func(w http.ResponseWriter, r *http.Request) {
		checkCreateMailRequest(t, r)

		w.WriteHeader(http.StatusBadRequest)
	})

	_, err := client.CreateMail(context.Background(), testCreateMailRequest)
	if err == nil {
		t.Fatalf("createMail request should fail when the server send status code %d", http.StatusBadRequest)
	}
}

func TestCloudOutboundMailerClient_CreateMail_ShouldFailIfMalformedResponseBody(t *testing.T) {
	client, mux, teardown := setupTestCloudOutboundMailerClient()
	defer teardown()

	mux.HandleFunc(mailer.EndpointMails, func(w http.ResponseWriter, r *http.Request) {
		checkCreateMailRequest(t, r)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("malformed response body :)"))
	})

	_, err := client.CreateMail(context.Background(), testCreateMailRequest)
	if err == nil {
		t.Fatalf("createMail request should fail when the server send malformed response body")
	}
}

func TestCloudOutboundMailerClient_CreateFiles(t *testing.T) {
	client, mux, teardown := setupTestCloudOutboundMailerClient()
	defer teardown()

	mux.HandleFunc(mailer.EndpointFiles, func(w http.ResponseWriter, r *http.Request) {
		checkCreateFilesRequest(t, r)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, validCreateFilesResponse)
	})

	got, err := client.CreateFiles(context.Background(), testFiles)
	if err != nil {
		t.Fatalf("createFile request was given a valid request but failed: %v", err)
	}

	want := mailer.CreateFileResponse{
		TempRequestID: "test-temp-request-id",
		Files: []*mailer.ResponseFileInfo{
			{
				FileName: "test-filename-1",
				FileSize: 42,
				FileID:   "test-file-id-1",
			},
			{
				FileName: "test-filename-2",
				FileSize: 1337,
				FileID:   "test-file-id-2",
			},
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Response differ from the expected one: %s", diff)
	}
}

func TestCloudOutboundMailerClient_CreateFile_ShouldFailIfWrongResponseStatusCode(t *testing.T) {
	client, mux, teardown := setupTestCloudOutboundMailerClient()
	defer teardown()

	mux.HandleFunc(mailer.EndpointFiles, func(w http.ResponseWriter, r *http.Request) {
		checkCreateFilesRequest(t, r)

		w.WriteHeader(http.StatusBadRequest)
	})

	_, err := client.CreateFiles(context.Background(), testFiles)
	if err == nil {
		t.Fatalf("createFile request should fail when the server send malformed response body")
	}
}

func TestCloudOutboundMailerClient_CreateFiles_ShouldFailIfMalformedResponseBody(t *testing.T) {
	client, mux, teardown := setupTestCloudOutboundMailerClient()
	defer teardown()

	mux.HandleFunc(mailer.EndpointFiles, func(w http.ResponseWriter, r *http.Request) {
		checkCreateFilesRequest(t, r)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("malformed response body :)"))
	})

	_, err := client.CreateFiles(context.Background(), testFiles)
	if err == nil {
		t.Fatalf("createFile request should fail when the server send status code %d", http.StatusBadRequest)
	}
}
