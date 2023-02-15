package sens_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/connectfit-team/naverapi/internal/testhelper"
	"github.com/connectfit-team/naverapi/sens"
	"github.com/google/go-cmp/cmp"
)

func TestSENSClient_SendSMS(t *testing.T) {
	client, mux, teardown := setupTestSENSClient()
	defer teardown()

	mux.HandleFunc(sens.EndpointMessages, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodPost)

		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", "856915200000")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", "test-access-key")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", "fNdhE7DMQ3QT2Ev3GxZxSpzg3vJmHNBMjAPG78IqBX8=")
		testhelper.TestRequestHeader(t, r, "Content-Type", "application/json")

		testhelper.TestRequestBody(t, r, `{"type":"LMS","contentType":"AD","countryCode":"82","from":"test-from","subject":"test-subject","content":"test-content","messages":[{"to":"test-to-1","subject":"test-subject-1","content":"test-content-1"},{"to":"test-to-2","subject":"test-subject-2","content":"test-content-2"}],"reserveTime":"test-reserve-time","reserveTimeZone":"test-reserve-time-zone","scheduleCode":"test-schedule-code"}`)

		w.WriteHeader(http.StatusAccepted)
		respBody := sens.SendSMSResponse{
			RequestID:   "test-request-id",
			RequestTime: "test-request-time",
			StatusCode:  "test-status-code",
			StatusName:  "success",
		}
		err := json.NewEncoder(w).Encode(respBody)
		if err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})

	req := sens.SendSMSRequest{
		Type:        sens.SMSTypeLMS,
		ContentType: sens.ContentTypeAD,
		CountryCode: sens.CountryCodeKorea,
		From:        "test-from",
		Subject:     "test-subject",
		Content:     "test-content",
		Messages: []sens.Message{
			{
				To:      "test-to-1",
				Subject: "test-subject-1",
				Content: "test-content-1",
			},
			{
				To:      "test-to-2",
				Subject: "test-subject-2",
				Content: "test-content-2",
			},
		},
		ReserveTime:     "test-reserve-time",
		ReserveTimeZone: "test-reserve-time-zone",
		ScheduleCode:    "test-schedule-code",
	}
	got, err := client.SendSMS(context.Background(), req)
	if err != nil {
		t.Fatalf("Send SMS request was given a valid request but failed: %v", err)
	}

	want := sens.SendSMSResponse{
		RequestID:   "test-request-id",
		RequestTime: "test-request-time",
		StatusCode:  "test-status-code",
		StatusName:  "success",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Response differ from the expected one: %s", diff)
	}
}

func TestSENSClient_SendSMS_ShouldFailIfWrongResponseStatusCode(t *testing.T) {
	client, mux, teardown := setupTestSENSClient()
	defer teardown()

	mux.HandleFunc(sens.EndpointMessages, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodPost)

		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", "856915200000")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", "test-access-key")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", "fNdhE7DMQ3QT2Ev3GxZxSpzg3vJmHNBMjAPG78IqBX8=")
		testhelper.TestRequestHeader(t, r, "Content-Type", "application/json")

		testhelper.TestRequestBody(t, r, `{"type":"LMS","contentType":"AD","countryCode":"82","from":"test-from","subject":"test-subject","content":"test-content","messages":[{"to":"test-to-1","subject":"test-subject-1","content":"test-content-1"},{"to":"test-to-2","subject":"test-subject-2","content":"test-content-2"}],"reserveTime":"test-reserve-time","reserveTimeZone":"test-reserve-time-zone","scheduleCode":"test-schedule-code"}`)

		w.WriteHeader(http.StatusBadRequest)
	})

	req := sens.SendSMSRequest{
		Type:        sens.SMSTypeLMS,
		ContentType: sens.ContentTypeAD,
		CountryCode: sens.CountryCodeKorea,
		From:        "test-from",
		Subject:     "test-subject",
		Content:     "test-content",
		Messages: []sens.Message{
			{
				To:      "test-to-1",
				Subject: "test-subject-1",
				Content: "test-content-1",
			},
			{
				To:      "test-to-2",
				Subject: "test-subject-2",
				Content: "test-content-2",
			},
		},
		ReserveTime:     "test-reserve-time",
		ReserveTimeZone: "test-reserve-time-zone",
		ScheduleCode:    "test-schedule-code",
	}
	_, err := client.SendSMS(context.Background(), req)
	if err == nil {
		t.Fatalf("Send SMS request should fail when the server send status code %d", http.StatusBadRequest)
	}
}

func TestSENSClient_SendSMS_ShouldFailIfMalformedResponseBody(t *testing.T) {
	client, mux, teardown := setupTestSENSClient()
	defer teardown()

	mux.HandleFunc(sens.EndpointMessages, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodPost)

		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Timestamp", "856915200000")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Iam-Access-Key", "test-access-key")
		testhelper.TestRequestHeader(t, r, "X-Ncp-Apigw-Signature-V2", "fNdhE7DMQ3QT2Ev3GxZxSpzg3vJmHNBMjAPG78IqBX8=")
		testhelper.TestRequestHeader(t, r, "Content-Type", "application/json")

		testhelper.TestRequestBody(t, r, `{"type":"LMS","contentType":"AD","countryCode":"82","from":"test-from","subject":"test-subject","content":"test-content","messages":[{"to":"test-to-1","subject":"test-subject-1","content":"test-content-1"},{"to":"test-to-2","subject":"test-subject-2","content":"test-content-2"}],"reserveTime":"test-reserve-time","reserveTimeZone":"test-reserve-time-zone","scheduleCode":"test-schedule-code"}`)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("malformed response body :)"))
	})

	req := sens.SendSMSRequest{
		Type:        sens.SMSTypeLMS,
		ContentType: sens.ContentTypeAD,
		CountryCode: sens.CountryCodeKorea,
		From:        "test-from",
		Subject:     "test-subject",
		Content:     "test-content",
		Messages: []sens.Message{
			{
				To:      "test-to-1",
				Subject: "test-subject-1",
				Content: "test-content-1",
			},
			{
				To:      "test-to-2",
				Subject: "test-subject-2",
				Content: "test-content-2",
			},
		},
		ReserveTime:     "test-reserve-time",
		ReserveTimeZone: "test-reserve-time-zone",
		ScheduleCode:    "test-schedule-code",
	}
	_, err := client.SendSMS(context.Background(), req)
	if err == nil {
		t.Fatalf("Send SMS request should fail when the server send malformed response body")
	}
}
