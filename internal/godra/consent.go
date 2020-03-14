package godra

import (
	"log"
	"net/http"
)

// GetConsentHandler handles the consent flow
// as godra is intended for internal use,
// all concent requests are accepted automatically
func (srv Server) GetConsentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		c := r.URL.Query().Get("consent_challenge")
		if c == "" {
			log.Printf("received empty consent_challenge at: %v\n", r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		body, err := srv.hydraclient.GetConsentRequest(c)
		if err != nil {
			log.Printf("error while querying consent request: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bodyAccept, err := srv.hydraclient.AcceptConsentRequest(c, true, 7200, body.GetRequestedScope(), body.GetRequestedAccessTokenAudience())
		http.Redirect(w, r, bodyAccept.GetRedirectTo(), http.StatusTemporaryRedirect)
	}
}
