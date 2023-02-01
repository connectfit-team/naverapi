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
)

var (
	invalidAddr       = ""
	validAddrExist    = "valid-addr-exist"
	validAddrNotExist = "valid-addr-not-exist"
	errorAddr         = "error-addr"
)

func TestClient_Query(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		wantID := "test-client-id"
		wantSecret := "test-client-secret"

		// {"error":{"errorCode":"300","message":"Not Found Exception","details":"URL not found."}}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method but got %s", r.Method)
		}
		// {"error":{"errorCode":"200","message":"Authentication Failed","details":"Invalid authentication information."}}
		if got := r.Header.Get(ClientIDHeaderKey); got != wantID {
			t.Error("invalid auth")
		}
		if got := r.Header.Get(ClientSecretHeaderKey); got != wantSecret {
			t.Error("invalid auth")
		}

		queryParam := r.URL.Query().Get(query)
		// {"status":"INVALID_REQUEST","errorMessage":"query is INVALID"}
		if queryParam == "" {
			t.Error("must include 'query' as parameter")
		}

		var resp *Response
		code := http.StatusOK

		switch queryParam {
		case validAddrExist:
			resp = &Response{
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
		case errorAddr:
			code = http.StatusInternalServerError
		default: // OK but address not exist
			// {"status":"OK","meta":{"totalCount":0,"count":0},"addresses":[],"errorMessage":""}
			resp = &Response{
				Status: "OK",
			}
		}

		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(*resp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})

	ctx := context.Background()
	// empty query
	_, err := client.Query(ctx, invalidAddr)
	if err == nil {
		t.Error("Expected invalid error but nil")
	}
	// not invalid but address not exist
	res, err := client.Query(ctx, validAddrNotExist)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if res.Status != OK {
		t.Errorf("Expected status 'OK' but got : %v", res.Status)
	}
	if int(res.Meta.Count) != len(res.Addresses) {
		t.Errorf("Expected meta.count == address.length but got : count=%d and length=%d",
			res.Meta.Count, len(res.Addresses))
	}
	// valid and exist address
	res, err = client.Query(ctx, validAddrExist)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if res.Status != OK {
		t.Errorf("Expected status 'OK' but got : %v", res.Status)
	}
	if int(res.Meta.Count) != len(res.Addresses) {
		t.Errorf("Expected meta.count == address.length but got : count=%d and length=%d",
			res.Meta.Count, len(res.Addresses))
	}
}

func TestClient_Language(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		wantID := "test-client-id"
		wantSecret := "test-client-secret"

		// {"error":{"errorCode":"300","message":"Not Found Exception","details":"URL not found."}}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method but got %s", r.Method)
		}
		// {"error":{"errorCode":"200","message":"Authentication Failed","details":"Invalid authentication information."}}
		if got := r.Header.Get(ClientIDHeaderKey); got != wantID {
			t.Error("invalid auth")
		}
		if got := r.Header.Get(ClientSecretHeaderKey); got != wantSecret {
			t.Error("invalid auth")
		}

		queryParam := r.URL.Query().Get(query)
		// {"status":"INVALID_REQUEST","errorMessage":"query is INVALID"}
		if queryParam == "" {
			t.Error("must include 'query' as parameter")
		}

		langParam := r.URL.Query().Get(language)
		toEnglish := langParam == string(Eng)

		var resp *Response
		code := http.StatusOK

		switch queryParam {
		case validAddrExist:
			addr := Address{
				RoadAddress:    "도로명 주소",
				JibunAddress:   "지번 주소",
				EnglishAddress: "영어 주소",
				X:              "127",
				Y:              "37",
				AddressElements: []AddressElement{
					{
						Types:    []string{"POSTAL_CODE"},
						LongName: "1111",
					},
				},
			}
			if toEnglish {
				addr.RoadAddress = "English Address"
				addr.JibunAddress = "English Address"
				addr.EnglishAddress = "English Address"
			}
			resp = &Response{
				Status:    "OK",
				Meta:      Meta{TotalCount: 1, Page: 1, Count: 1},
				Addresses: []Address{addr},
			}
		case errorAddr:
			code = http.StatusInternalServerError
		default: // OK but address not exist
			// {"status":"OK","meta":{"totalCount":0,"count":0},"addresses":[],"errorMessage":""}
			resp = &Response{
				Status: "OK",
			}
		}

		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(*resp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})

	ctx := context.Background()
	// set english language option
	res, err := client.Language(Eng).Query(ctx, validAddrExist)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if res.Status != OK {
		t.Errorf("Expected status 'OK' but got : %v", res.Status)
	}
	if int(res.Meta.Count) != len(res.Addresses) {
		t.Errorf("Expected meta.count == address.length but got : count=%d and length=%d",
			res.Meta.Count, len(res.Addresses))
	}
	if res.Addresses[0].RoadAddress != "English Address" {
		t.Errorf("Expected English Address but got : %s", res.Addresses[0].RoadAddress)
	}

	// find default
	res, err = client.Query(ctx, validAddrExist)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if res.Status != OK {
		t.Errorf("Expected status 'OK' but got : %v", res.Status)
	}
	if int(res.Meta.Count) != len(res.Addresses) {
		t.Errorf("Expected meta.count == address.length but got : count=%d and length=%d",
			res.Meta.Count, len(res.Addresses))
	}
	if res.Addresses[0].RoadAddress == "English Address" {
		t.Errorf("Expected Korean Address but got : %s", res.Addresses[0].RoadAddress)
	}
}

func TestClient_Coordinate(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		wantID := "test-client-id"
		wantSecret := "test-client-secret"

		// {"error":{"errorCode":"300","message":"Not Found Exception","details":"URL not found."}}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method but got %s", r.Method)
		}
		// {"error":{"errorCode":"200","message":"Authentication Failed","details":"Invalid authentication information."}}
		if got := r.Header.Get(ClientIDHeaderKey); got != wantID {
			t.Error("invalid auth")
		}
		if got := r.Header.Get(ClientSecretHeaderKey); got != wantSecret {
			t.Error("invalid auth")
		}

		queryParam := r.URL.Query().Get(query)
		// {"status":"INVALID_REQUEST","errorMessage":"query is INVALID"}
		if queryParam == "" {
			t.Error("must include 'query' as parameter")
		}

		calculateDistance := false
		coordinateParam := r.URL.Query().Get(coordinate)
		if coordinateParam != "" {
			points := strings.Split(coordinateParam, ",")
			if len(points) != 2 {
				t.Errorf("Expected longitue, latitude but got : %s", coordinateParam)
			}
			_, err := strconv.ParseFloat(points[0], 64)
			if err != nil {
				t.Errorf("Expected longitude as type float but got : %v", points[0])
			}
			_, err = strconv.ParseFloat(points[1], 64)
			if err != nil {
				t.Errorf("Expected latitude as type float but got : %v", points[1])
			}
			calculateDistance = true
		}

		var resp *Response
		code := http.StatusOK

		switch queryParam {
		case validAddrExist:
			addr := Address{
				RoadAddress:    "도로명 주소",
				JibunAddress:   "지번 주소",
				EnglishAddress: "영어 주소",
				X:              "127",
				Y:              "37",
				AddressElements: []AddressElement{
					{
						Types:    []string{"POSTAL_CODE"},
						LongName: "1111",
					},
				},
			}
			if calculateDistance {
				addr.Distance = 100
			}
			resp = &Response{
				Status:    "OK",
				Meta:      Meta{TotalCount: 1, Page: 1, Count: 1},
				Addresses: []Address{addr},
			}
		case errorAddr:
			code = http.StatusInternalServerError
		default: // OK but address not exist
			// {"status":"OK","meta":{"totalCount":0,"count":0},"addresses":[],"errorMessage":""}
			resp = &Response{
				Status: "OK",
			}
		}

		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(*resp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})

	ctx := context.Background()
	// calculate distance
	res, err := client.Coordinate(127, 37).Query(ctx, validAddrExist)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if res.Status != OK {
		t.Errorf("Expected status 'OK' but got : %v", res.Status)
	}
	if int(res.Meta.Count) != len(res.Addresses) {
		t.Errorf("Expected meta.count == address.length but got : count=%d and length=%d",
			res.Meta.Count, len(res.Addresses))
	}
	if res.Addresses[0].Distance != 100 {
		t.Errorf("Expected distance but got : %f", res.Addresses[0].Distance)
	}

	// find default
	res, err = client.Query(ctx, validAddrExist)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if res.Status != OK {
		t.Errorf("Expected status 'OK' but got : %v", res.Status)
	}
	if int(res.Meta.Count) != len(res.Addresses) {
		t.Errorf("Expected meta.count == address.length but got : count=%d and length=%d",
			res.Meta.Count, len(res.Addresses))
	}
	if res.Addresses[0].Distance != 0 {
		t.Errorf("Expected zero distance but got : %f", res.Addresses[0].Distance)
	}
}

func TestClient_HCodeFilter(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		wantID := "test-client-id"
		wantSecret := "test-client-secret"

		// {"error":{"errorCode":"300","message":"Not Found Exception","details":"URL not found."}}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method but got %s", r.Method)
		}
		// {"error":{"errorCode":"200","message":"Authentication Failed","details":"Invalid authentication information."}}
		if got := r.Header.Get(ClientIDHeaderKey); got != wantID {
			t.Error("invalid auth")
		}
		if got := r.Header.Get(ClientSecretHeaderKey); got != wantSecret {
			t.Error("invalid auth")
		}

		queryParam := r.URL.Query().Get(query)
		// {"status":"INVALID_REQUEST","errorMessage":"query is INVALID"}
		if queryParam == "" {
			t.Error("must include 'query' as parameter")
		}

		fcode := ""
		fcodeVal := ""
		codeFilterParam := r.URL.Query().Get(filter)
		if codeFilterParam != "" {
			sep := strings.Split(codeFilterParam, "@")
			if len(sep) != 2 {
				t.Errorf("Expected code correct filter formant 'HCODE@[code1];[code2]' but got : %s", codeFilterParam)
			}
			if sep[0] != string(hCode) && sep[0] != string(bCode) {
				t.Errorf("Expected code filter HCODE or BCODE but got : %v", sep[0])
			}
			fcode = sep[0]
			sepCode := strings.Split(sep[1], ";")
			_, err := strconv.ParseInt(sepCode[1], 10, 64)
			if err != nil {
				t.Errorf("Expected code filter as type int but got : %v", sepCode[1])
			}
			fcodeVal = sep[1]
		}

		var resp *Response
		code := http.StatusOK

		switch queryParam {
		case validAddrExist:
			addr := Address{
				RoadAddress:     "도로명 주소",
				JibunAddress:    "지번 주소",
				EnglishAddress:  "영어 주소",
				X:               "127",
				Y:               "37",
				AddressElements: []AddressElement{{}},
			}
			if fcode != "" {
				addr.AddressElements[0].Types = []string{fcode}
				addr.AddressElements[0].Code = fcodeVal
			}
			resp = &Response{
				Status:    "OK",
				Meta:      Meta{TotalCount: 1, Page: 1, Count: 1},
				Addresses: []Address{addr},
			}
		case errorAddr:
			code = http.StatusInternalServerError
		default: // OK but address not exist
			// {"status":"OK","meta":{"totalCount":0,"count":0},"addresses":[],"errorMessage":""}
			resp = &Response{
				Status: "OK",
			}
		}

		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(*resp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})
	first := "4113554500"
	second := "4113555000"
	ctx := context.Background()
	// multiple hcode filter
	res, err := client.
		HCode(first).
		HCode(second).
		Query(ctx, validAddrExist)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if res.Status != OK {
		t.Errorf("Expected status 'OK' but got : %v", res.Status)
	}
	if int(res.Meta.Count) != len(res.Addresses) {
		t.Errorf("Expected meta.count == address.length but got : count=%d and length=%d",
			res.Meta.Count, len(res.Addresses))
	}
	if res.Addresses[0].AddressElements[0].Types[0] != string(hCode) {
		t.Errorf("Expected code filter as 'HCODE' but got : %v",
			res.Addresses[0].AddressElements[0].Types[0])
	}
	if res.Addresses[0].AddressElements[0].Code != fmt.Sprintf("%s;%s", first, second) {
		t.Errorf("Expected code as '%s' but got : %v",
			fmt.Sprintf("%s;%s", first, second), res.Addresses[0].AddressElements[0].Code)
	}
}

func TestClient_Page(t *testing.T) {
	client, mux, tearDown := setupTestClient()
	defer tearDown()

	mux.HandleFunc(Endpoint, func(w http.ResponseWriter, r *http.Request) {
		wantID := "test-client-id"
		wantSecret := "test-client-secret"

		// {"error":{"errorCode":"300","message":"Not Found Exception","details":"URL not found."}}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method but got %s", r.Method)
		}
		// {"error":{"errorCode":"200","message":"Authentication Failed","details":"Invalid authentication information."}}
		if got := r.Header.Get(ClientIDHeaderKey); got != wantID {
			t.Error("invalid auth")
		}
		if got := r.Header.Get(ClientSecretHeaderKey); got != wantSecret {
			t.Error("invalid auth")
		}

		queryParam := r.URL.Query().Get(query)
		// {"status":"INVALID_REQUEST","errorMessage":"query is INVALID"}
		if queryParam == "" {
			t.Error("must include 'query' as parameter")
		}

		pageParam := r.URL.Query().Get(page)
		if pageParam != "" {
			_, err := strconv.ParseInt(pageParam, 10, 64)
			if err != nil {
				t.Errorf("Expected page as type int but got : %v", pageParam)
			}
		}

		var resp *Response
		code := http.StatusOK

		switch queryParam {
		case validAddrExist:
			if pageParam == "1" {
				resp = &Response{
					Status: "OK",
					Meta:   Meta{TotalCount: 1, Page: 1, Count: 1},
					Addresses: []Address{{
						RoadAddress:    "도로명 주소",
						JibunAddress:   "지번 주소",
						EnglishAddress: "영어 주소",
						X:              "127",
						Y:              "37",
					}},
				}
			} else {
				resp = &Response{
					Status: "OK",
					Meta:   Meta{TotalCount: 0, Page: 0, Count: 0},
				}
			}
		case errorAddr:
			code = http.StatusInternalServerError
		default: // OK but address not exist
			// {"status":"OK","meta":{"totalCount":0,"count":0},"addresses":[],"errorMessage":""}
			resp = &Response{
				Status: "OK",
			}
		}

		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(*resp); err != nil {
			t.Errorf("encoding the response body shouldn't fail but got: %v", err)
		}
	})

	ctx := context.Background()
	// first page
	res, err := client.
		Page(1).
		Query(ctx, validAddrExist)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if res.Status != OK {
		t.Errorf("Expected status 'OK' but got : %v", res.Status)
	}
	if int(res.Meta.Count) != len(res.Addresses) {
		t.Errorf("Expected meta.count == address.length but got : count=%d and length=%d",
			res.Meta.Count, len(res.Addresses))
	}
	// second page
	res, err = client.
		Page(2).
		Query(ctx, validAddrExist)
	if err != nil {
		t.Errorf("Expected nil but got : %v", err)
	}
	if res.Status != OK {
		t.Errorf("Expected status 'OK' but got : %v", res.Status)
	}
	if int(res.Meta.Count) != len(res.Addresses) {
		t.Errorf("Expected meta.count == address.length but got : count=%d and length=%d",
			res.Meta.Count, len(res.Addresses))
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
