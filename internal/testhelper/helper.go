package testhelper

import (
	"io"
	"net/http"
	"testing"
)

func TestRequestMethod(t *testing.T, r *http.Request, method string) {
	t.Helper()

	if r.Method != method {
		t.Errorf("Expected POST request method but got %s", r.Method)
	}
}

func TestRequestHeader(t *testing.T, r *http.Request, headerKey, want string) {
	t.Helper()

	if got := r.Header.Get(headerKey); got != want {
		t.Errorf("Expected header %q to have value %q but got %q instead", headerKey, want, got)
	}
}

func TestRequestBody(t *testing.T, r *http.Request, want string) {
	t.Helper()

	got, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Failed to read the request body: %v", err)
	}

	if string(got) != want {
		t.Errorf("Expected body was:\n%s\nbut got:\n%s", want, got)
	}
}

func TestRequestFormFiles(t *testing.T, r *http.Request, key string, want []string) {
	t.Helper()

	err := r.ParseMultipartForm(0)
	if err != nil {
		t.Errorf("failed to parse the request body as multipart/form-data: %v", err)
	}

	fileHeaders := r.MultipartForm.File[key]
	for i, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			t.Errorf("could not open the file %s associated to the multipart request: %v", fileHeader.Filename, err)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			t.Errorf("could not read the content of the file %s associated to the multipart request: %v", fileHeader.Filename, err)
		}

		contentStr := string(content)
		if contentStr != want[i] {
			t.Errorf("expected content for file %s was:\n%s\nBut got:\n%s", fileHeader.Filename, want[i], contentStr)
		}
	}
}
