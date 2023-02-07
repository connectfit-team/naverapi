package mailer_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/connectfit-team/naverapi/internal/testhelper"
	"github.com/connectfit-team/naverapi/mailer"
	"github.com/google/go-cmp/cmp"
)

func TestMailServic_CreateMail(t *testing.T) {
	client, mux, teardown := setupTestCloudOutboundMailerClient()
	defer teardown()

	mux.HandleFunc(mailer.EndpointMails, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodPost)

		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", "856915200000")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", "test-access-key")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", "F1YxxwEjDRZmNLxqqDFz53OpbvLrMCqEsv9tLxoBcWE=")
		testhelper.TestRequestHeader(t, r, "Content-Type", "application/json")

		testhelper.TestRequestBody(t, r, `{"senderAddress":"test-sender-address","senderName":"test-sender-name","title":"test-title","body":"test-body","recipients":[{"address":"test-address","name":"test-name","type":"test-type","parameters":["test-parameter-1","test-parameter-2"]}],"attachFileIds":["test-file-id-1","test-file-id-2"]}`)

		w.WriteHeader(http.StatusCreated)
		respBody := mailer.CreateMailResponse{
			RequestID: "test-request-id",
			Count:     1,
		}
		err := json.NewEncoder(w).Encode(respBody)
		if err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})

	req := mailer.CreateMailRequest{
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
	got, err := client.CreateMail(context.Background(), req)
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
		testhelper.TestRequestMethod(t, r, http.MethodPost)

		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", "856915200000")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", "test-access-key")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", "F1YxxwEjDRZmNLxqqDFz53OpbvLrMCqEsv9tLxoBcWE=")
		testhelper.TestRequestHeader(t, r, "Content-Type", "application/json")

		testhelper.TestRequestBody(t, r, `{"senderAddress":"test-sender-address","senderName":"test-sender-name","title":"test-title","body":"test-body","recipients":[{"address":"test-address","name":"test-name","type":"test-type","parameters":["test-parameter-1","test-parameter-2"]}],"attachFileIds":["test-file-id-1","test-file-id-2"]}`)

		w.WriteHeader(http.StatusBadRequest)
	})

	req := mailer.CreateMailRequest{
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
	_, err := client.CreateMail(context.Background(), req)
	if err == nil {
		t.Fatalf("createMail request should fail when the server send status code %d", http.StatusBadRequest)
	}
}

func TestCloudOutboundMailerClient_CreateMail_ShouldFailIfMalformedResponseBody(t *testing.T) {
	client, mux, teardown := setupTestCloudOutboundMailerClient()
	defer teardown()

	mux.HandleFunc(mailer.EndpointMails, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodPost)

		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", "856915200000")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", "test-access-key")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", "F1YxxwEjDRZmNLxqqDFz53OpbvLrMCqEsv9tLxoBcWE=")
		testhelper.TestRequestHeader(t, r, "Content-Type", "application/json")

		testhelper.TestRequestBody(t, r, `{"senderAddress":"test-sender-address","senderName":"test-sender-name","title":"test-title","body":"test-body","recipients":[{"address":"test-address","name":"test-name","type":"test-type","parameters":["test-parameter-1","test-parameter-2"]}],"attachFileIds":["test-file-id-1","test-file-id-2"]}`)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("malformed response body :)"))
	})

	req := mailer.CreateMailRequest{
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
	_, err := client.CreateMail(context.Background(), req)
	if err == nil {
		t.Fatalf("createMail request should fail when the server send malformed response body")
	}
}

func TestCloudOutboundMailerClient_CreateFiles(t *testing.T) {
	client, mux, teardown := setupTestCloudOutboundMailerClient()
	defer teardown()

	mux.HandleFunc(mailer.EndpointFiles, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodPost)

		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", "856915200000")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", "test-access-key")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", "q1JhbYivx0lU//wBoOyh+yn/y7+Lg9Ez/Xj6FDzxap4=")

		testhelper.TestRequestFormFiles(t, r, "fileList", []string{"test-content-1", "test-content-2"})

		w.WriteHeader(http.StatusCreated)
		respBody := mailer.CreateFileResponse{
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
		err := json.NewEncoder(w).Encode(respBody)
		if err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})

	files := []mailer.File{
		{
			Name:    "test-name-1",
			Content: []byte("test-content-1"),
		},
		{
			Name:    "test-name-2",
			Content: []byte("test-content-2"),
		},
	}
	got, err := client.CreateFiles(context.Background(), files)
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
		testhelper.TestRequestMethod(t, r, http.MethodPost)

		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", "856915200000")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", "test-access-key")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", "q1JhbYivx0lU//wBoOyh+yn/y7+Lg9Ez/Xj6FDzxap4=")

		testhelper.TestRequestFormFiles(t, r, "fileList", []string{"test-content-1", "test-content-2"})

		w.WriteHeader(http.StatusBadRequest)
	})

	files := []mailer.File{
		{
			Name:    "test-name-1",
			Content: []byte("test-content-1"),
		},
		{
			Name:    "test-name-2",
			Content: []byte("test-content-2"),
		},
	}
	_, err := client.CreateFiles(context.Background(), files)
	if err == nil {
		t.Fatalf("createFile request should fail when the server send malformed response body")
	}
}

func TestCloudOutboundMailerClient_CreateFiles_ShouldFailIfMalformedResponseBody(t *testing.T) {
	client, mux, teardown := setupTestCloudOutboundMailerClient()
	defer teardown()

	mux.HandleFunc(mailer.EndpointFiles, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodPost)

		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", "856915200000")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", "test-access-key")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", "q1JhbYivx0lU//wBoOyh+yn/y7+Lg9Ez/Xj6FDzxap4=")

		testhelper.TestRequestFormFiles(t, r, "fileList", []string{"test-content-1", "test-content-2"})

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("malformed response body :)"))
	})

	files := []mailer.File{
		{
			Name:    "test-name-1",
			Content: []byte("test-content-1"),
		},
		{
			Name:    "test-name-2",
			Content: []byte("test-content-2"),
		},
	}
	_, err := client.CreateFiles(context.Background(), files)
	if err == nil {
		t.Fatalf("createFile request should fail when the server send status code %d", http.StatusBadRequest)
	}
}
