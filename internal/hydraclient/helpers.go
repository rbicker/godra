package hydraclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// error response for all kind of flows.
type errorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorDebug       string `json:"error_debug"`
	StatusCode       int    `json:"status_code"`
}

// determineError extracts the error message from the
// http response if status code is unexpected.
func determineError(res *http.Response) error {
	if res.StatusCode < 200 || res.StatusCode > 302 {
		var resBody errorResponse
		err := json.NewDecoder(res.Body).Decode(&resBody)
		if err != nil {
			return err
		}
		return fmt.Errorf("%v; %s; %s; %s", resBody.StatusCode, resBody.Error, resBody.ErrorDescription, resBody.ErrorDebug)
	}
	return nil
}

// Get queries the given challenge from hydra.
// The flow can be "login", "consent" or "logout".
func (c *Client) Get(flow, challenge string) (*http.Response, error) {
	if flow != "login" && flow != "consent" && flow != "logout" {
		return nil, fmt.Errorf("invalid flow: %s", flow)
	}
	if challenge == "" {
		return nil, fmt.Errorf("empty challenge given for flow %s", flow)
	}
	params := url.Values{}
	params.Add(fmt.Sprintf("%s_challenge", flow), challenge)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/oauth2/auth/requests/%s?%s", c.hydraPrivateURL, flow, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("X-Forwarded-Proto", "https")
	client := http.Client{
		Timeout: time.Second * 5,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	if err = determineError(res); err != nil {
		return nil, fmt.Errorf("unexpected response: %w", err)
	}
	return res, nil
}

// Put sends an "accept" or "deny" request with the given body
// to hydra, depending on the given action.
// The flow can be "login", "consent" or "logout".
func (c *Client) Put (flow, action, challenge string, body []byte) (*http.Response, error) {
	if flow != "login" && flow != "consent" && flow != "logout" {
		return nil, fmt.Errorf("invalid flow: %s", flow)
	}
	if challenge == "" {
		return nil, fmt.Errorf("empty challenge given for flow %s", flow)
	}
	if action != "accept" && action != "deny" {
		return nil, fmt.Errorf("invalid action: %s", action)
	}
	params := url.Values{}
	params.Add(fmt.Sprintf("%s_challenge", flow), challenge)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/oauth2/auth/requests/%s/%s?%s", c.hydraPrivateURL, flow, action, params.Encode()),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("X-Forwarded-Proto", "https")
	client := http.Client{
		Timeout: time.Second * 5,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	if err = determineError(res); err != nil {
		return nil, fmt.Errorf("unexpected response: %w", err)
	}
	return res, nil
}