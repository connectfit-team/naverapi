package geocode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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
	normalResp = Response{
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
	emptyResp = Response{
		Status: "OK",
	}
	engResp = Response{
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
	coordinateResp = Response{
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
	if !reflect.DeepEqual(*res, normalResp) {
		t.Error("Not expected response")
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
	if !reflect.DeepEqual(*res, emptyResp) {
		t.Error("Not expected response")
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

		resp := coordinateResp
		if !checkCoordinate(t, r, fmt.Sprintf("%f,%f", lon, lat)) {
			resp = normalResp
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Coordinate(lon, lat).
		Query(context.Background(), validAddr)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if !reflect.DeepEqual(*res, coordinateResp) {
		t.Error("Not expected response")
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

		resp := engResp
		if !checkLanguage(t, r, string(LanguageEng)) {
			resp = normalResp
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Language(LanguageEng).
		Query(context.Background(), validAddr)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if !reflect.DeepEqual(*res, engResp) {
		t.Error("Not expected response")
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
		resp := normalResp
		if !checkCode(t, r, fmt.Sprintf("%s@%s", string(hCode), strings.Join(hCodes, ";"))) {
			resp = emptyResp
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.HCode(hCodes...).
		Query(context.Background(), validAddr)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if !reflect.DeepEqual(*res, normalResp) {
		t.Error("Not expected response")
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
		resp := normalResp
		if checkPage(t, r, pageVal) {
			resp = emptyResp
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Page(pageVal).
		Query(context.Background(), validAddr)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if !reflect.DeepEqual(*res, emptyResp) {
		t.Error("Not expected response")
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
		resp := normalResp
		if checkCount(t, r, countVal) {
			resp = emptyResp
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Count(countVal).
		Query(context.Background(), validAddr)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if !reflect.DeepEqual(*res, emptyResp) {
		t.Error("Not expected response")
	}
}

func checkCount(t *testing.T, r *http.Request, v int) (equal bool) {
	vs := fmt.Sprintf("%d", v)
	if !cmp.Equal(vs, r.URL.Query().Get(count)) {
		t.Errorf("Expected count %s but got %s", vs, r.URL.Query().Get(count))
	}
	return cmp.Equal(vs, r.URL.Query().Get(count))
}

func checkPage(t *testing.T, r *http.Request, v int) (equal bool) {
	vs := fmt.Sprintf("%d", v)
	if !cmp.Equal(vs, r.URL.Query().Get(page)) {
		t.Errorf("Expected page %s but got %s", vs, r.URL.Query().Get(page))
	}
	return cmp.Equal(vs, r.URL.Query().Get(page))
}

func checkCode(t *testing.T, r *http.Request, v string) (equal bool) {
	if !cmp.Equal(v, r.URL.Query().Get(filter)) {
		t.Fatalf("Expected code %s but got %s", v, r.URL.Query().Get(filter))
	}
	return cmp.Equal(v, r.URL.Query().Get(filter))
}

func checkLanguage(t *testing.T, r *http.Request, v string) (equal bool) {
	if !cmp.Equal(v, r.URL.Query().Get(language)) {
		t.Errorf("Expected language %s but got %s", v, r.URL.Query().Get(language))
	}
	return cmp.Equal(v, r.URL.Query().Get(language))
}

func checkCoordinate(t *testing.T, r *http.Request, v string) (equal bool) {
	if !cmp.Equal(v, r.URL.Query().Get(coordinate)) {
		t.Errorf("Expected coordinate %s but got %s", v, r.URL.Query().Get(coordinate))
	}
	return cmp.Equal(v, r.URL.Query().Get(coordinate))
}

func checkQuery(t *testing.T, r *http.Request, v string) (equal bool) {
	if !cmp.Equal(v, r.URL.Query().Get(query)) {
		t.Errorf("Expected query %s but got %s", v, r.URL.Query().Get(query))
	}
	return cmp.Equal(v, r.URL.Query().Get(query))
}

func checkHeaderID(t *testing.T, r *http.Request) {
	if !cmp.Equal(wantID, r.Header.Get(clientIDHeaderKey)) {
		t.Errorf("Expected Header %s=%s but got %s=%s",
			clientIDHeaderKey, wantID, clientIDHeaderKey, r.Header.Get(clientIDHeaderKey))
	}
}

func checkHeaderSecret(t *testing.T, r *http.Request) {
	if !cmp.Equal(wantSecret, r.Header.Get(clientSecretHeaderKey)) {
		t.Errorf("Expected Header %s=%s but got %s=%s",
			clientSecretHeaderKey, wantSecret, clientSecretHeaderKey, r.Header.Get(clientSecretHeaderKey))
	}
}

func checkMethod(t *testing.T, r *http.Request) {
	if !cmp.Equal(http.MethodGet, r.Method) {
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
