package godra

import (
	"log"
	"net/http"
)

// GetLogoutHandler handles the logout flow.
// All logout requests are accepted automatically.
func (srv Server) GetLogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		c := r.URL.Query().Get("logout_challenge")
		if c == "" {
			log.Printf("received empty logout_challenge at: %v\n", r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err := srv.hydraclient.GetLogoutRequest(c)
		if err != nil {
			log.Printf("error while querying logout request: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bodyAccept, err := srv.hydraclient.AcceptLogoutRequest(c)
		http.Redirect(w, r, bodyAccept.GetRedirectTo(), http.StatusTemporaryRedirect)
	}
}
