package net

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/usagifm/dating-app/lib/logger"
)

const (
	MethodPost                = "POST"
	MethodGet                 = "GET"
	MethodPut                 = "PUT"
	MethodDelete              = "DELETE"
	ContentType               = "Content-Type"
	ApplicationJson           = "application/json"
	ApplicationFormUrlEncoded = "application/x-www-form-urlencoded"
	CacheControl              = "Cache-Control"
	NoCache                   = "no-cache"
	Authorization             = "Authorization"
	ServerError               = 500
)

func Post(ctx context.Context, endpoint string, headers map[string]string, payload []byte) ([]byte, int, error) {
	req, err := http.NewRequest(MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when build request, err: ", err)
		return nil, ServerError, err
	}

	// Set the content type to JSON (adjust according to your needs)
	req.Header.Set(ContentType, ApplicationJson)

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Perform the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when do request, err: ", err)
		return nil, ServerError, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when read response, err: ", err)
		return nil, ServerError, err
	}

	return bodyBytes, resp.StatusCode, nil
}

// PostForm performs an HTTP POST request with a form-urlencoded body
func PostForm(ctx context.Context, endpoint string, formValues url.Values, headers map[string]string) ([]byte, int, error) {
	// Encode the form values
	formBody := formValues.Encode()

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(formBody))
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when build request, err: ", err)
		return nil, ServerError, err
	}

	// Set the Content-Type header to application/x-www-form-urlencoded
	req.Header.Set(ContentType, ApplicationFormUrlEncoded)

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Perform the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when do request, err: ", err)
		return nil, ServerError, err
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when read response, err: ", err)
		return nil, ServerError, err
	}

	return bodyBytes, resp.StatusCode, nil
}

func Put(ctx context.Context, endpoint string, headers map[string]string, payload []byte) ([]byte, int, error) {
	req, err := http.NewRequest(MethodPut, endpoint, bytes.NewReader(payload))
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when build request, err: ", err)
		return nil, ServerError, err
	}

	// Set the content type to JSON (adjust according to your needs)
	req.Header.Set(ContentType, ApplicationJson)

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Perform the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when do request, err: ", err)
		return nil, ServerError, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when read response, err: ", err)
		return nil, ServerError, err
	}

	return bodyBytes, resp.StatusCode, nil
}

func Get(ctx context.Context, endpoint string, headers map[string]string, queryParams map[string]string) ([]byte, int, error) {
	// Create a new URL with the provided endpoint and query parameters
	u, err := url.Parse(endpoint)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error parsing endpoint, err: ", err)
		return nil, ServerError, err
	}

	// Add query parameters to the URL
	q := u.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(MethodGet, u.String(), nil)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when build request, err: ", err)
		return nil, ServerError, err
	}

	// Set the content type to JSON (adjust according to your needs)
	req.Header.Set(ContentType, ApplicationJson)

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Perform the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when do request, err: ", err)
		return nil, ServerError, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when read response, err: ", err)
		return nil, ServerError, err
	}

	return bodyBytes, resp.StatusCode, nil
}

func Delete(ctx context.Context, endpoint string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest(MethodDelete, endpoint, nil)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when build request, err: ", err)
		return nil, ServerError, err
	}

	// Set the content type to JSON (adjust according to your needs)
	req.Header.Set(ContentType, ApplicationJson)

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Perform the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when do request, err: ", err)
		return nil, ServerError, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error when read response, err: ", err)
		return nil, ServerError, err
	}

	return bodyBytes, resp.StatusCode, nil
}
