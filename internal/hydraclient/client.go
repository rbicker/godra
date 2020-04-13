package hydraclient

import (
	"encoding/json"
	"fmt"
)

// Client represents a hyra client.
type Client struct {
	hydraPrivateURL string
}

// ensure Client implements the HydraClient interface.
var _ HydraClient = Client{}

// SetHydraPrivateURL sets the hydra server's private URL.
func (c *Client) SetHydraPrivateURL(url string) {
	c.hydraPrivateURL = url
}

type getLoginRequestResponse struct {
	Skip    bool   `json:"skip"`
	Subject string `json:"subject"`
}

func (r getLoginRequestResponse) GetSkip() bool {
	return r.Skip
}
func (r getLoginRequestResponse) GetSubject() string {
	return r.Subject
}

// GetLoginRequest queries the login request from hydra.
func (c Client) GetLoginRequest(challenge string) (GetLoginRequestResponse, error) {
	res, err := c.Get("login", challenge)
	if err != nil {
		return nil, fmt.Errorf("receiving login request failed: %w", err)
	}
	defer res.Body.Close()
	var resBody getLoginRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
	}
	return resBody, nil
}

type acceptLoginRequestResponse struct {
	RedirectTo string `json:"redirect_to"`
}

func (r acceptLoginRequestResponse) GetRedirectTo() string {
	return r.RedirectTo
}

// AcceptLoginRequest accepts the login request by
// responding to the hydra server.
func (c Client) AcceptLoginRequest(challenge string, remember bool, rememberFor int, subject string) (AcceptLoginRequestResponse, error) {
	reqBody, err := json.Marshal(struct {
		Remember    bool   `json:"remember"`
		RememberFor int    `json:"remember_for"`
		Subject     string `json:"subject"`
	}{
		Remember:    remember,
		RememberFor: rememberFor,
		Subject:     subject,
	})
	if err != nil {
		return nil, fmt.Errorf("cloud not create request body: %w", err)
	}
	res, err := c.Put("login", "accept", challenge, reqBody)
	if err != nil {
		return nil, fmt.Errorf("accepting login request failed: %w", err)
	}
	var resBody acceptLoginRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w",err)
	}
	return resBody, nil
}

type rejectLoginRequestResponse struct {
	RedirectTo string `json:"redirect_to"`
}

func (r rejectLoginRequestResponse) GetRedirectTo() string {
	return r.RedirectTo
}

// RejectLoginRequest rejects the login request
// by responding to the hydra server.
func (c Client) RejectLoginRequest(challenge string, errorID string, errorDescription string) (RejectLoginRequestResponse, error) {
	reqBody, err := json.Marshal(map[string]string{
		"error":             errorID,
		"error_description": errorDescription,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create request body: %w", err)
	}
	res, err := c.Put("login", "reject", challenge, reqBody)
	if err != nil {
		return nil, fmt.Errorf("rejecting login request failed: %w", err)
	}
	defer res.Body.Close()
	var resBody rejectLoginRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w", err)
	}
	return resBody, nil
}

type getConsentRequestResponse struct {
	Subject                      string   `json:"subject"`
	RequestedScope               []string `json:"requested_scope"`
	RequestedAccessTokenAudience []string `json:"requested_access_token_audience"`
}

func (r getConsentRequestResponse) GetSubject() string {
	return r.Subject
}

func (r getConsentRequestResponse) GetRequestedScope() []string {
	return r.RequestedScope
}

func (r getConsentRequestResponse) GetRequestedAccessTokenAudience() []string {
	return r.RequestedAccessTokenAudience
}

// GetConsentRequest queries additional
// information from the hydra server.
func (c Client) GetConsentRequest(challenge string) (GetConsentRequestResponse, error) {
	// query consent request information from hydra
	// using given login challenge
	res, err := c.Get("consent", challenge)
	if err != nil {
		return nil, fmt.Errorf("receiving consent request failed: %w", err)
	}
	defer res.Body.Close()
	var resBody getConsentRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w", err)
	}
	return resBody, nil
}

type acceptConsentRequestResponse struct {
	RedirectTo string `json:"redirect_to"`
}

func (r acceptConsentRequestResponse) GetRedirectTo() string {
	return r.RedirectTo
}

// AcceptConsentRequest accepts the consent request
// by responding to the hydra server.
func (c Client) AcceptConsentRequest(challenge string, remember bool, rememberFor int, grantScope []string, grantAccessTokenAudience []string) (AcceptConsentRequestResponse, error) {
	reqBody, err := json.Marshal(struct {
		Remember                 bool     `json:"remember"`
		RememberFor              int      `json:"remember_for"`
		GrantScope               []string `json:"grant_scope"`
		GrantAccessTokenAudience []string `json:"grant_access_token_audience"`
	}{
		Remember:                 remember,
		RememberFor:              rememberFor,
		GrantScope:               grantScope,
		GrantAccessTokenAudience: grantAccessTokenAudience,
	})
	if err != nil {
		return nil, fmt.Errorf("cloud not create request body: %w", err)
	}
	res, err := c.Put("consent", "accept", challenge, reqBody)
	if err != nil {
		return nil, fmt.Errorf("accepting consent request failed: %w", err)
	}
	var resBody acceptConsentRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w",err)
	}
	return resBody, nil
}

type getLogoutRequestResponse struct {
	Subject string `json:"subject"`
}

func (r getLogoutRequestResponse) GetSubject() string {
	return r.Subject
}

// GetLogoutRequest queries the logout request from hydra.
func (c Client) GetLogoutRequest(challenge string) (GetLogoutRequestResponse, error) {
	// query logout request information from hydra
	// using given login challenge
	res, err := c.Get("logout", challenge)
	if err != nil {
		return nil, fmt.Errorf("receiving logout request failed: %w", err)
	}
	defer res.Body.Close()
	var resBody getLogoutRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
	}
	return resBody, nil
}

type acceptLogoutRequestResponse struct {
	RedirectTo string `json:"redirect_to"`
}

func (r acceptLogoutRequestResponse) GetRedirectTo() string {
	return r.RedirectTo
}

// AcceptLogoutRequest accepts the logout request
// by responding to the hydra server.
func (c Client) AcceptLogoutRequest(challenge string) (AcceptLogoutRequestResponse, error) {
	res, err := c.Put("logout", "accept", challenge, nil)
	if err != nil {
		return nil, fmt.Errorf("accepting logout request failed: %w", err)
	}
	var resBody acceptConsentRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w",err)
	}
	return resBody, nil
}
