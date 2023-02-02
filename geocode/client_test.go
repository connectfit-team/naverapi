package geocode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	validAddr   = "valid-addr"
	invalidAddr = "invalid-addr"
)

var (
	wantID     = "test-client-id"
	wantSecret = "test-client-secret"
)

var (
	normalResp = &Response{
		Status: "OK",
		Meta:   Meta{TotalCount: 1, Page: 1, Count: 1},
		Addresses: []Address{
			{
				RoadAddress:    "도로명 주소",
				JibunAddress:   "지번 주소",
				EnglishAddress: "English Address",
				X:              "127",
				Y:              "37",
				AddressElements: []AddressElement{
					{
						Types:    []string{"POSTAL_CODE"},
						LongName: "1111",
					},
				},
			},
		},
	}
	emptyResp = &Response{
		Status: "OK",
	}
	engResp = &Response{
		Status: "OK",
		Meta:   Meta{TotalCount: 1, Page: 1, Count: 1},
		Addresses: []Address{
			{
				RoadAddress:    "English Address",
				JibunAddress:   "English Address",
				EnglishAddress: "English Address",
				X:              "127",
				Y:              "37",
				AddressElements: []AddressElement{
					{
						Types:    []string{"POSTAL_CODE"},
						LongName: "1111",
					},
				},
			},
		},
	}
	coordinateResp = &Response{
		Status: "OK",
		Meta:   Meta{TotalCount: 1, Page: 1, Count: 1},
		Addresses: []Address{
			{
				RoadAddress:    "도로명 주소",
				JibunAddress:   "지번 주소",
				EnglishAddress: "English Address",
				X:              "127",
				Y:              "37",
				Distance:       1,
				AddressElements: []AddressElement{
					{
						Types:    []string{"POSTAL_CODE"},
						LongName: "1111",
					},
				},
			},
		},
	}
)

func TestClient_Query_Valid(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r)
		checkHeaderID(t, r)
		checkHeaderSecret(t, r)

		checkQuery(t, r, validAddr)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(normalResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if diff := cmp.Diff(res, normalResp); diff != "" {
		t.Errorf("Not expected response: %s", diff)
	}
}

func TestClient_Query_Invalid(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r)
		checkHeaderID(t, r)
		checkHeaderSecret(t, r)

		checkQuery(t, r, invalidAddr)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(emptyResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), invalidAddr)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if diff := cmp.Diff(res, emptyResp); diff != "" {
		t.Errorf("Not expected response: %s", diff)
	}
}

func TestClient_Coordinate(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	lon, lat := 37.0, 127.0

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r)
		checkHeaderID(t, r)
		checkHeaderSecret(t, r)

		checkQuery(t, r, validAddr)
		checkCoordinate(t, r, fmt.Sprintf("%f,%f", lon, lat))

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(coordinateResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr, WithCoordinate(lon, lat))
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if diff := cmp.Diff(res, coordinateResp); diff != "" {
		t.Errorf("Not expected response: %s", diff)
	}
}

func TestClient_Language(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r)
		checkHeaderID(t, r)
		checkHeaderSecret(t, r)

		checkQuery(t, r, validAddr)
		checkLanguage(t, r, string(LanguageEng))

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(engResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr, WithLanguage(LanguageEng))
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if diff := cmp.Diff(res, engResp); diff != "" {
		t.Errorf("Not expected response: %s", diff)
	}
}

func TestClient_HCode(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	hCodes := []string{"value1", "value2"}

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r)
		checkHeaderID(t, r)
		checkHeaderSecret(t, r)

		checkQuery(t, r, validAddr)
		checkCode(t, r, fmt.Sprintf("%s@%s", string(hCode), strings.Join(hCodes, ";")))

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(normalResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr, WithHCode(hCodes...))
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if diff := cmp.Diff(res, normalResp); diff != "" {
		t.Errorf("Not expected response: %s", diff)
	}
}

func TestClient_Page(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	pageVal := 2

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r)
		checkHeaderID(t, r)
		checkHeaderSecret(t, r)

		checkQuery(t, r, validAddr)
		checkPage(t, r, pageVal)

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(emptyResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr, WithPage(pageVal))
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if diff := cmp.Diff(res, emptyResp); diff != "" {
		t.Errorf("Not expected response: %s", diff)
	}
}

func TestClient_Count(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	countVal := 0

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r)
		checkHeaderID(t, r)
		checkHeaderSecret(t, r)

		checkQuery(t, r, validAddr)
		checkCount(t, r, countVal)

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(emptyResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr, WithCount(countVal))
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if diff := cmp.Diff(res, emptyResp); diff != "" {
		t.Errorf("Not expected response: %s", diff)
	}
}

func checkCount(t *testing.T, r *http.Request, v int) (equal bool) {
	expect := strconv.Itoa(v)
	actual := r.URL.Query().Get(countConst)
	equal = expect == actual
	if !equal {
		t.Errorf("Expected count %s but got %s", expect, actual)
	}
	return equal
}

func checkPage(t *testing.T, r *http.Request, v int) (equal bool) {
	expect := strconv.Itoa(v)
	actual := r.URL.Query().Get(pageConst)
	equal = expect == actual
	if !equal {
		t.Errorf("Expected page %s but got %s", expect, actual)
	}
	return equal
}

func checkCode(t *testing.T, r *http.Request, v string) (equal bool) {
	actual := r.URL.Query().Get(filterConst)
	equal = v == actual
	if !equal {
		t.Fatalf("Expected code %s but got %s", v, actual)
	}
	return equal
}

func checkLanguage(t *testing.T, r *http.Request, v string) (equal bool) {
	actual := r.URL.Query().Get(languageConst)
	equal = v == actual
	if !equal {
		t.Errorf("Expected language %s but got %s", v, actual)
	}
	return equal
}

func checkCoordinate(t *testing.T, r *http.Request, v string) (equal bool) {
	actual := r.URL.Query().Get(coordinateConst)
	equal = v == actual
	if !equal {
		t.Errorf("Expected coordinate %s but got %s", v, actual)
	}
	return equal
}

func checkQuery(t *testing.T, r *http.Request, v string) (equal bool) {
	actual := r.URL.Query().Get(queryConst)
	equal = v == actual
	if !equal {
		t.Errorf("Expected query %s but got %s", v, actual)
	}
	return equal
}

func checkHeaderID(t *testing.T, r *http.Request) {
	id := r.Header.Get(clientIDHeaderKey)
	if wantID != id {
		t.Errorf("Expected Header %s=%s but got %s=%s",
			clientIDHeaderKey, wantID, clientIDHeaderKey, id)
	}
}

func checkHeaderSecret(t *testing.T, r *http.Request) {
	secret := r.Header.Get(clientSecretHeaderKey)
	if wantSecret != secret {
		t.Errorf("Expected Header %s=%s but got %s=%s",
			clientSecretHeaderKey, wantSecret, clientSecretHeaderKey, secret)
	}
}

func checkMethod(t *testing.T, r *http.Request) {
	if http.MethodGet != r.Method {
		t.Errorf("Expected Method %s but got %s", http.MethodGet, r.Method)
	}
}

func setupTestClient() (client *Client, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()

	srv := httptest.NewServer(mux)

	srvURL, _ := url.Parse(srv.URL)
	client, _ = NewClient(
		"test-client-id",
		"test-client-secret",
		srv.Client(),
	)
	client.BaseURL = srvURL

	return client, mux, srv.Close
}
