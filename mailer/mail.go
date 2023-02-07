package mailer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/connectfit-team/naverapi/internal/httputil"
)

const (
	CloudOutboundMailerDefaultBaseURL = "https://mail.apigw.ntruss.com" // [Base URL of the API]: https://api.ncloud-docs.com/docs/en/ai-application-service-cloudoutboundmailer

	EndpointFiles = "/api/v1/files" // [Files endpoint]: https://api.ncloud-docs.com/docs/en/ai-application-service-cloudoutboundmailer-createfile
	EndpointMails = "/api/v1/mails" // [Mails endpoint]: https://api.ncloud-docs.com/docs/en/ai-application-service-cloudoutboundmailer-createmailrequest
)

// CloudOutboundMailerClient is a REST client providing methods to interact with the
// Naver Cloud Outbound Mailer API.
type CloudOutboundMailerClient struct {
	// HTTPClient is the HTTP client used internally to perform the requests.
	HTTPClient *http.Client
	// BaseURL is the base URL which prefix every request's URL path.
	// e.g. https://sens.apigw.ntruss.com
	BaseURL *url.URL
	// AccessKey is the access key (from portal or sub account)
	AccessKey string
	// SecretKey is the secret key (from portal or sub account)
	SecretKey string
	// Clock provides the current time used to fill the `x-ncp-apigw-timestamp`
	// header of each request to the API.
	// It has been made public mainly for testing purpose to avoid polluting the
	// public API with extra parameters or options.
	Clock Clock
}

// NewCloudOutboundMailerClient returns a new Naver Cloud Outbound mail API client given a
// base URL, an access key and a secret key to authenticate to the API.
// It uses the http.DefaultClient unless you provide your own.
func NewCloudOutboundMailerClient(accessKey, secretKey string, httpClient *http.Client) (*CloudOutboundMailerClient, error) {
	url, err := url.Parse(CloudOutboundMailerDefaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("malformed base URL %q: %w", CloudOutboundMailerDefaultBaseURL, err)
	}

	svc := &CloudOutboundMailerClient{
		HTTPClient: http.DefaultClient,
		BaseURL:    url,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Clock:      &realClock{},
	}

	if httpClient != nil {
		svc.HTTPClient = httpClient
	}

	return svc, nil
}

// Error represents an API request error.
type Error struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

// CreateMailRequest represents a createMail request.
//
// See https://api.ncloud-docs.com/docs/ai-application-service-cloudoutboundmailer-createmailrequest
type CreateMailRequest struct {
	SenderAddress string       `json:"senderAddress"`
	SenderName    string       `json:"senderName"`
	Title         string       `json:"title"`
	Body          string       `json:"body"`
	Recipients    []*Recipient `json:"recipients"`
	AttachFileIDs []string     `json:"attachFileIds,omitempty"`
}

type Recipient struct {
	Address    string   `json:"address"`
	Name       string   `json:"name"`
	Type       string   `json:"type"` //_R:수신자, C:참조자, B:숨은 참조자
	Parameters []string `json:"parameters,omitempty"`
}

// CreateMailResponse represents the response sent by the Naver Cloud Outbound
// API after a createMail request.
//
// See https://api.ncloud-docs.com/docs/ai-application-service-cloudoutboundmailer-createmailrequest
type CreateMailResponse struct {
	RequestID string `json:"requestId"`
	Count     int    `json:"count"`
	Error     Error  `json:"error,omitempty"`
}

// CreateMail sends a createMail request to the API using the given request
// parameters.
//
// See https://api.ncloud-docs.com/docs/ai-application-service-cloudoutboundmailer-createmailrequest
func (comc *CloudOutboundMailerClient) CreateMail(ctx context.Context, req CreateMailRequest) (CreateMailResponse, error) {
	endpoint := comc.BaseURL.JoinPath(EndpointMails)
	httpReq, err := httputil.NewJSONBodyRequest(ctx, http.MethodPost, endpoint.String(), req)
	if err != nil {
		return CreateMailResponse{}, fmt.Errorf("could not build the createMail request: %w", err)
	}

	timestamp := strconv.FormatInt(comc.Clock.Now().UnixMilli(), 10)
	err = httputil.SetNCloudRequestHeaders(httpReq, EndpointMails, timestamp, comc.AccessKey, comc.SecretKey)
	if err != nil {
		return CreateMailResponse{}, fmt.Errorf("failed to set the HTTP header of the request: %w", err)
	}

	resp, err := comc.HTTPClient.Do(httpReq)
	if err != nil {
		return CreateMailResponse{}, fmt.Errorf("could not perform the HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return CreateMailResponse{}, fmt.Errorf("request failed with code %d: %s", resp.StatusCode, resp.Status)
	}

	var responseBody CreateMailResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return CreateMailResponse{}, fmt.Errorf("could not response decode body: %w", err)
	}

	return responseBody, nil
}

// File represents a file to upload to Naver Cloud.
type File struct {
	Name    string
	Content []byte
}

// CreateFileResponse represents the response sent by the Naver Cloud Outbound
// API after a createFile request.
//
// See https://api.ncloud-docs.com/docs/ai-application-service-cloudoutboundmailer-createmailrequest
type CreateFileResponse struct {
	TempRequestID string              `json:"tempRequestId"`
	Files         []*ResponseFileInfo `json:"files"`
	Error         Error               `json:"error,omitempty"`
}

type ResponseFileInfo struct {
	FileName string `json:"fileName"`
	FileSize int    `json:"fileSize"`
	FileID   string `json:"fileId"`
}

// CreateFiles sends a createFile request to the API using the given request
// parameters.
//
// See https://api.ncloud-docs.com/docs/ai-application-service-cloudoutboundmailer-createmailrequest
func (comc *CloudOutboundMailerClient) CreateFiles(ctx context.Context, files []File) (CreateFileResponse, error) {
	endpoint := comc.BaseURL.JoinPath(EndpointFiles).String()
	req, err := newMultipartFormFileBodyRequest(ctx, http.MethodPost, endpoint, files)
	if err != nil {
		return CreateFileResponse{}, fmt.Errorf("could not build the createFile request: %w", err)
	}

	timestamp := strconv.FormatInt(comc.Clock.Now().UnixMilli(), 10)
	err = httputil.SetNCloudRequestHeaders(req, EndpointFiles, timestamp, comc.AccessKey, comc.SecretKey)
	if err != nil {
		return CreateFileResponse{}, fmt.Errorf("failed to set the HTTP header of the request: %w", err)
	}

	resp, err := comc.HTTPClient.Do(req)
	if err != nil {
		return CreateFileResponse{}, fmt.Errorf("could not perform the HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return CreateFileResponse{}, fmt.Errorf("request failed with code %d: %s", resp.StatusCode, resp.Status)
	}

	var responseBody CreateFileResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return CreateFileResponse{}, fmt.Errorf("could not decode response body: %w", err)
	}

	return responseBody, nil
}
