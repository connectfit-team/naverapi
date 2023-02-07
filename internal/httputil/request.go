package httputil

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewJSONBodyRequest(ctx context.Context, method, endpoint string, body any) (*http.Request, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("could not marshal the request body struct to JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("could not create the %s HTTP request given the URL %q: %w", method, endpoint, err)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func SetNCloudRequestHeaders(req *http.Request, url, timestamp, accessKey, secretKey string) error {
	apigwSignature, err := formatAPIGatewaySignature(http.MethodPost, url, timestamp, accessKey, secretKey)
	if err != nil {
		return fmt.Errorf("could not format the API gateway signature: %w", err)
	}

	req.Header.Add("x-ncp-apigw-timestamp", timestamp)
	req.Header.Add("x-ncp-iam-access-key", accessKey)
	req.Header.Add("x-ncp-apigw-signature-v2", apigwSignature)

	return nil
}

// See https://api.ncloud-docs.com/docs/ai-application-service-cloudoutboundmailer
func formatAPIGatewaySignature(method, url, timestamp, accessKey, secretKey string) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(method)
	buf.WriteString(" ")
	buf.WriteString(url)
	buf.WriteString("\n")
	buf.WriteString(timestamp)
	buf.WriteString("\n")
	buf.WriteString(accessKey)
	message := buf.String()

	hm := hmac.New(sha256.New, []byte(secretKey))
	_, err := hm.Write([]byte(message))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hm.Sum(nil)), err
}
