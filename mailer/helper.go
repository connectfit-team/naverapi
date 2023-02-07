package mailer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func newMultipartFormFileBodyRequest(ctx context.Context, method, endpoint string, files []File) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	err := addFilesToRequest(writer, files...)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close the multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("could not create the %s HTTP request given the URL %q: %w", method, endpoint, err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

func addFilesToRequest(writer *multipart.Writer, files ...File) error {
	for _, file := range files {
		err := addFileToRequest(writer, file)
		if err != nil {
			return fmt.Errorf("could not add file %s to the request body: %w", file.Name, err)
		}
	}
	return nil
}

func addFileToRequest(writer *multipart.Writer, file File) error {
	part, err := writer.CreateFormFile("fileList", file.Name)
	if err != nil {
		return fmt.Errorf("could not create the form file for %s: %w", file.Name, err)
	}

	_, err = io.Copy(part, bytes.NewBuffer(file.Content))
	if err != nil {
		return fmt.Errorf("failed to copy the content of %s into the form: %w", file.Name, err)
	}

	return nil
}
