package hydraclient

// GetLoginRequestResponse represents a response from a
// GetLoginRequest.
type GetLoginRequestResponse interface {
	GetSkip() bool
	GetSubject() string
}

// AcceptLoginRequestResponse represents a response from a
// AcceptLoginRequest.
type AcceptLoginRequestResponse interface {
	GetRedirectTo() string
}

// RejectLoginRequestResponse represents a response from a
// RejectLoginRequest.
type RejectLoginRequestResponse interface {
	GetRedirectTo() string
}

// GetConsentRequestResponse represents a response from a
// GetConsentRequest.
type GetConsentRequestResponse interface {
	GetSubject() string
	GetRequestedScope() []string
	GetRequestedAccessTokenAudience() []string
}

// AcceptConsentRequestResponse represents a response from a
// AcceptConsentRequest.
type AcceptConsentRequestResponse interface {
	GetRedirectTo() string
}

// GetLogoutRequestResponse represents a response from a
// GetLogoutRequest.
type GetLogoutRequestResponse interface {
	GetSubject() string
}

// AcceptLogoutRequestResponse represents a response from a
// AcceptLogoutRequest.
type AcceptLogoutRequestResponse interface {
	GetRedirectTo() string
}

// HydraClient describe all client functions
// to interact with hydra.
type HydraClient interface {
	GetLoginRequest(challenge string) (GetLoginRequestResponse, error)
	AcceptLoginRequest(challenge string, remember bool, rememberFor int, subject string) (AcceptLoginRequestResponse, error)
	RejectLoginRequest(challenge string, errorID string, errorDescription string) (RejectLoginRequestResponse, error)
	GetConsentRequest(challenge string) (GetConsentRequestResponse, error)
	AcceptConsentRequest(challenge string, remember bool, rememberFor int, grantScope []string, grantAccessTokenAudience []string) (AcceptConsentRequestResponse, error)
	GetLogoutRequest(challenge string) (GetLogoutRequestResponse, error)
	AcceptLogoutRequest(challenge string) (AcceptLogoutRequestResponse, error)
}
