package sens

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/connectfit-team/naverapi/internal/httputil"
)

const (
	SENSDefaultBaseURL = "https://sens.apigw.ntruss.com"

	EndpointSMSAPI   = "/sms/v2"
	EndpointMessages = EndpointSMSAPI + "services/ncp:sms:kr:273248310490:sms_service/messages"
)

var (
	// ErrSendSMSFailed is returned when the `statusName` in the response after
	// requesting to send a SMS is not "success".
	ErrSendSMSFailed = errors.New(`the send SMS request's response did not return status "success"`)
)

type SMSType string

const (
	SMSTypeSMS SMSType = "SMS"
	SMSTypeLMS SMSType = "LMS"
	SMSTypeMMS SMSType = "MMS"
)

type SMSContentType string

const (
	ContentTypeSMS SMSContentType = "COMM"
	ContentTypeAD  SMSContentType = "AD"
)

type SMSCountryCode string

const (
	CountryCodeKorea SMSCountryCode = "82"
)

// Client is a REST client providing methods to interact with
// Naver Cloud Platform SENS(Simple & Easy Notification Service).
type Client struct {
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

// NewClient returns a new Naver Cloud Platform SMS API client given a
// base URL, an access key and a secret key to authenticate to the API.
// It uses the http.DefaultClient unless you provide your own.
func NewClient(accessKey, secretKey string, httpClient *http.Client) (*Client, error) {
	baseURL, err := url.Parse(SENSDefaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("malformed base URL %q: %w", baseURL, err)
	}

	svc := &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    baseURL,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Clock:      &realClock{},
	}

	if httpClient != nil {
		svc.HTTPClient = httpClient
	}

	return svc, nil
}

// SendSMSRequest represents a the REST request to send a SMS.
//
// See https://api.ncloud-docs.com/docs/en/ai-application-service-sens-smsv2
type SendSMSRequest struct {
	Type            SMSType        `json:"type"`                      // SMS 타입 (SMS | LMS | MMS) - 필수
	ContentType     SMSContentType `json:"contentType"`               // 메세지 타입 (COMM(일반) | AD(광고)) - 필수
	CountryCode     SMSCountryCode `json:"countryCode,omitempty"`     // 국가 코드(default 82) - 선택
	From            string         `json:"from"`                      // 문자 발송 번호 - 필수
	Subject         string         `json:"subject,omitempty"`         // 문자 제목(LMS, MMS 만 사용) - 선택
	Content         string         `json:"content"`                   // 기본 문자 내용(EUC-KR 인코딩, 지원 외 이모지 발송 실패) - 필수
	Messages        []Message      `json:"messages"`                  // 문자 정보 - 필수
	ReserveTime     string         `json:"reserveTime,omitempty"`     // 예약 시간("yyyy-MM-dd HH:mm") - 선택
	ReserveTimeZone string         `json:"reserveTimeZone,omitempty"` // 예약 시간 타임존("Asia/Seoul") - 선택
	ScheduleCode    string         `json:"scheduleCode,omitempty"`
}

type Message struct {
	To      string `json:"to"`                // 문자 수신 번호 - 필수
	Subject string `json:"subject,omitempty"` // 개별 문자 내용 - 선택
	Content string `json:"content,omitempty"` // 개별 문자 제목(LMS, MMS 만 사용) - 선택
}

// SendSMSResponse represents the response sent by the SENS SMS API
// API after a request to send a SMS.
//
// See https://api.ncloud-docs.com/docs/en/ai-application-service-sens-smsv2
type SendSMSResponse struct {
	RequestID   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	StatusCode  string `json:"statusCode"`
	StatusName  string `json:"statusName"`
}

// SendSMS sends a request to send a SMS to the SMS API using the given request
// parameters.
//
// See https://api.ncloud-docs.com/docs/en/ai-application-service-sens-smsv2
func (ss *Client) SendSMS(ctx context.Context, req SendSMSRequest) (SendSMSResponse, error) {
	endpoint := ss.BaseURL.JoinPath(EndpointMessages).String()
	httpReq, err := httputil.NewJSONBodyRequest(ctx, http.MethodPost, endpoint, req)
	if err != nil {
		return SendSMSResponse{}, fmt.Errorf("could not build the send message request: %w", err)
	}

	timestamp := strconv.FormatInt(ss.Clock.Now().UnixMilli(), 10)
	err = httputil.SetNCloudRequestHeaders(httpReq, EndpointMessages, timestamp, ss.AccessKey, ss.SecretKey)
	if err != nil {
		return SendSMSResponse{}, fmt.Errorf("failed to set the HTTP header of the request: %w", err)
	}

	resp, err := ss.HTTPClient.Do(httpReq)
	if err != nil {
		return SendSMSResponse{}, fmt.Errorf("could not perform the HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return SendSMSResponse{}, fmt.Errorf("request failed with code %d: %s", resp.StatusCode, resp.Status)
	}

	var responseBody SendSMSResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return SendSMSResponse{}, fmt.Errorf("could not response decode body: %w", err)
	}

	if responseBody.StatusName != "success" {
		return SendSMSResponse{}, ErrSendSMSFailed
	}

	return responseBody, nil
}
