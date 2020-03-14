package hydraclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client represents a hyra client.
type Client struct {
	hydraPrivateURL string
}

// error response for all kind of flows.
type errorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorDebug       string `json:"error_debug"`
	StatusCode       int    `json:"status_code"`
}

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

// GetLoginRequest queries the login request from hydra.
func (c Client) GetLoginRequest(challenge string) (GetLoginRequestResponse, error) {
	params := url.Values{}
	params.Add("login_challenge", challenge)
	res, err := http.Get(fmt.Sprintf("%v/oauth2/auth/requests/login?%v", c.hydraPrivateURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err = determineError(res); err != nil {
		return nil, err
	}
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
		return nil, err
	}
	client := http.Client{
		Timeout: time.Second * 5,
	}
	params := url.Values{}
	params.Add("login_challenge", challenge)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%v/oauth2/auth/requests/login/accept?%v", c.hydraPrivateURL, params.Encode()),
		bytes.NewBuffer(reqBody),
	)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err = determineError(res); err != nil {
		return nil, err
	}
	var resBody acceptLoginRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
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
		return nil, err
	}
	client := http.Client{
		Timeout: time.Second * 5,
	}
	params := url.Values{}
	params.Add("login_challenge", challenge)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%v/oauth2/auth/requests/login/reject?%v", c.hydraPrivateURL, params.Encode()),
		bytes.NewBuffer(reqBody),
	)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err = determineError(res); err != nil {
		return nil, err
	}
	var resBody rejectLoginRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
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
	params := url.Values{}
	params.Add("consent_challenge", challenge)
	res, err := http.Get(fmt.Sprintf("%v/oauth2/auth/requests/consent?%v", c.hydraPrivateURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err = determineError(res); err != nil {
		return nil, err
	}
	var resBody getConsentRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
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
		return nil, err
	}
	client := http.Client{
		Timeout: time.Second * 5,
	}
	params := url.Values{}
	params.Add("consent_challenge", challenge)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%v/oauth2/auth/requests/consent/accept?%v", c.hydraPrivateURL, params.Encode()),
		bytes.NewBuffer(reqBody),
	)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err = determineError(res); err != nil {
		return nil, err
	}
	var body acceptConsentRequestResponse

	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	return body, nil
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
	params := url.Values{}
	params.Add("consent_challenge", challenge)
	res, err := http.Get(fmt.Sprintf("%v/oauth2/auth/requests/logout?%v", c.hydraPrivateURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err = determineError(res); err != nil {
		return nil, err
	}
	var resBody GetLogoutRequestResponse
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
	client := http.Client{
		Timeout: time.Second * 5,
	}
	params := url.Values{}
	params.Add("logout_challenge", challenge)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%v/oauth2/auth/requests/login/accept?%v", c.hydraPrivateURL, params.Encode()),
		nil,
	)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err = determineError(res); err != nil {
		return nil, err
	}
	var resBody acceptLogoutRequestResponse
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
	}
	return resBody, nil
}
