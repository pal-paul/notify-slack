package slack

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SlackResponse handles parsing out errors from the web api.
type SlackResponse struct {
	Ok               bool                  `json:"ok"`
	Error            string                `json:"error"`
	Channel          string                `json:"channel"`
	Ts               string                `json:"ts"`
	ResponseMetadata SlackResponseMetadata `json:"response_metadata"`
}

type SlackResponseMetadata struct {
	Cursor   string   `json:"next_cursor"`
	Messages []string `json:"messages"`
	Warnings []string `json:"warnings"`
}

// RateLimitedError represents the rate limit response from slack
type RateLimitedError struct {
	RetryAfter time.Duration
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("slack rate limit exceeded, retry after %s", e.RetryAfter)
}

func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode == http.StatusTooManyRequests {
		retry, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return err
		}
		return &RateLimitedError{time.Duration(retry) * time.Second}
	}
	// Slack seems to send an HTML body along with 5xx error codes. Don't parse it.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (m *slack) postRequest(
	endpoint string,
	headers map[string]string,
	reqBody []byte,
) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(
		m.ctx,
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return resp, fmt.Errorf("error post to slack: %v", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.token))
	resp, err = m.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	if err := checkStatusCode(resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (m *slack) postForm(
	endpoint string,
	headers map[string]string,
	values url.Values,
) (resp *http.Response, err error) {
	reqBody := strings.NewReader(values.Encode())
	req, err := http.NewRequestWithContext(m.ctx, http.MethodPost, endpoint, reqBody)
	if err != nil {
		return resp, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.token))
	resp, err = m.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	if err := checkStatusCode(resp); err != nil {
		return resp, err
	}
	return resp, nil
}
