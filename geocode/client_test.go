package geocode_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/connectfit-team/naverapi/geocode"
	"github.com/connectfit-team/naverapi/internal/testhelper"
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
	normalResp = &geocode.Response{
		Status: "OK",
		Meta:   geocode.Meta{TotalCount: 1, Page: 1, Count: 1},
		Addresses: []geocode.Address{
			{
				RoadAddress:    "도로명 주소",
				JibunAddress:   "지번 주소",
				EnglishAddress: "English Address",
				X:              "127",
				Y:              "37",
				AddressElements: []geocode.AddressElement{
					{
						Types:    []string{"POSTAL_CODE"},
						LongName: "1111",
					},
				},
			},
		},
	}
	emptyResp = &geocode.Response{
		Status: "OK",
	}
	engResp = &geocode.Response{
		Status: "OK",
		Meta:   geocode.Meta{TotalCount: 1, Page: 1, Count: 1},
		Addresses: []geocode.Address{
			{
				RoadAddress:    "English Address",
				JibunAddress:   "English Address",
				EnglishAddress: "English Address",
				X:              "127",
				Y:              "37",
				AddressElements: []geocode.AddressElement{
					{
						Types:    []string{"POSTAL_CODE"},
						LongName: "1111",
					},
				},
			},
		},
	}
	coordinateResp = &geocode.Response{
		Status: "OK",
		Meta:   geocode.Meta{TotalCount: 1, Page: 1, Count: 1},
		Addresses: []geocode.Address{
			{
				RoadAddress:    "도로명 주소",
				JibunAddress:   "지번 주소",
				EnglishAddress: "English Address",
				X:              "127",
				Y:              "37",
				Distance:       1,
				AddressElements: []geocode.AddressElement{
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

	mux.HandleFunc(geocode.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodGet)

		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY-ID", wantID)
		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY", wantSecret)

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

	mux.HandleFunc(geocode.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodGet)

		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY-ID", wantID)
		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY", wantSecret)

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

	mux.HandleFunc(geocode.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodGet)

		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY-ID", wantID)
		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY", wantSecret)

		checkQuery(t, r, validAddr)
		checkCoordinate(t, r, fmt.Sprintf("%f,%f", lon, lat))

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(coordinateResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr, geocode.WithCoordinate(lon, lat))
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

	mux.HandleFunc(geocode.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodGet)

		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY-ID", wantID)
		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY", wantSecret)

		checkQuery(t, r, validAddr)
		checkLanguage(t, r, string(geocode.LanguageEng))

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(engResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr, geocode.WithLanguage(geocode.LanguageEng))
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

	mux.HandleFunc(geocode.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodGet)

		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY-ID", wantID)
		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY", wantSecret)

		checkQuery(t, r, validAddr)
		checkCode(t, r, fmt.Sprintf("%s@%s", "HCODE", strings.Join(hCodes, ";")))

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(normalResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr, geocode.WithHCode(hCodes...))
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

	mux.HandleFunc(geocode.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodGet)

		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY-ID", wantID)
		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY", wantSecret)

		checkQuery(t, r, validAddr)
		checkPage(t, r, "2")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(emptyResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr, geocode.WithPage(pageVal))
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

	mux.HandleFunc(geocode.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestRequestMethod(t, r, http.MethodGet)

		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY-ID", wantID)
		testhelper.TestRequestHeader(t, r, "X-NCP-APIGW-API-KEY", wantSecret)

		checkQuery(t, r, validAddr)
		checkCount(t, r, "0")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(emptyResp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	res, err := client.Query(context.Background(), validAddr, geocode.WithCount(countVal))
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if diff := cmp.Diff(res, emptyResp); diff != "" {
		t.Errorf("Not expected response: %s", diff)
	}
}
