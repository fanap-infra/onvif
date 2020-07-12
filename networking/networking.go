package networking

import (
	"bytes"
	"net/http"
	"time"

	dac "github.com/xinsnake/go-http-digest-auth-client"
)

// SendSoap send soap message
func SendSoap(username, password, endpoint, message string) (*http.Response, error) {
	dr := dac.NewRequest(username, password, "POST", endpoint, message)
	resp, err := dr.Execute()
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// SendSoapWithTimeout send soap message with timeOut
func SendSoapWithTimeout(endpoint string, message []byte, timeout time.Duration) (*http.Response, error) {
	httpClient := &http.Client{
		Timeout: timeout,
	}

	return httpClient.Post(endpoint, "application/soap+xml; charset=utf-8", bytes.NewReader(message))
}
