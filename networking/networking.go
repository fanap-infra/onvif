package networking

import (
	"bytes"
	"net/http"
	"strings"
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

func SendSoapWithDigestAndHttpReq(username, password, endpoint, message string, timeout time.Duration) (*http.Response, error) {
	t := dac.NewTransport(username, password)
	t.HTTPClient = &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(message))

	if err != nil {
		return nil, err
	}

	resp, err := t.RoundTrip(req)
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
