package http_client

import (
	"context"
	"net/http"
)

type httpClient struct {
	client *http.Client
	ctx    context.Context
}

type HttpClientInterface interface {
	// Get a http request to url with headers
	//
	// Parameters:
	//   - url: string
	//   - headers: map[string]string
	//
	// Returns:
	//   - []byte: response body
	//   - int: response status code
	//   - error: error
	Get(url string, headers map[string]string) ([]byte, int, error)

	// Post a http request to url with headers
	//
	// Parameters:
	//   - url: string
	//   - postBody: []byte
	//   - headers: map[string]string
	//
	// Returns:
	//   - []byte: response body
	//   - int: response status code
	//   - error: error
	Post(url string, postBody []byte, headers map[string]string) ([]byte, int, error)

	// Put a http request to url with headers
	//
	// Parameters:
	//   - url: string
	//   - postBody: []byte
	//   - headers: map[string]string
	//
	// Returns:
	//   - []byte: response body
	//   - int: response status code
	//   - error: error
	Put(url string, postBody []byte, headers map[string]string) ([]byte, int, error)

	// Delete a http request to url with headers
	//
	// Parameters:
	//   - url: string
	//   - postBody: []byte
	//   - headers: map[string]string
	//
	// Returns:
	//   - []byte: response body
	//   - int: response status code
	//   - error: error
	Delete(url string, postBody []byte, headers map[string]string) ([]byte, int, error)
}
