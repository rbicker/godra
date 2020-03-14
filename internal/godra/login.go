package godra

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/rbicker/godra/internal/nogo"
)

// read content from template file located in the environment
// variable with the given "envName". The string will be wrapped
// in a template definition with the given "tmplName".
// Errors will be logged.
// On error or if environment variable is not defined, the string
// given as "def" will be used.
func readTemplateFromFile(tmplName string, envName string, def string) string {
	content := def
	if p, ok := os.LookupEnv(envName); ok {
		file, err := os.Open(p)
		if err != nil {
			log.Printf("error while trying to open custom template file '%s' under '%s': %s\n", tmplName, p, err)
		} else {
			defer file.Close()
			b, err := ioutil.ReadAll(file)
			if err != nil {
				log.Printf("error while trying to read custom template file '%s' under '%s': %s\n", tmplName, p, err)
			} else {
				content = string(b)
			}
		}
	}
	return fmt.Sprintf(`{{ define "%s" }}%s{{ end }}`, tmplName, content)
}

// renderLoginForm renders the login form.
func renderLoginForm(w http.ResponseWriter, challenge string, alert string) {
	n, err := nogo.Get("/assets/templates/login.html")
	if err != nil {
		log.Printf("error while opening login html file: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	header := readTemplateFromFile("header", "CUSTOM_HEADER_PATH", "<h2>Login</h2>")
	footer := readTemplateFromFile("footer", "CUSTOM_FOOTER_PATH", "")
	t, err := template.New("login").Parse(header + footer + string(n.Content))
	if err != nil {
		log.Printf("error while parsing login html template: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var stylesheet string
	if ss, ok := os.LookupEnv("CUSTOM_STYLESHEET_PATH"); ok {
		stylesheet = fmt.Sprintf(`<link rel="stylesheet" type="text/css" href="%s">`, ss)
	}
	inputs := struct {
		Challenge  string
		Alert      string
		Stylesheet string
	}{
		Challenge:  challenge,
		Alert:      alert,
		Stylesheet: stylesheet,
	}
	t.Execute(w, inputs)
}

// GetLoginHandler returns the handler for the /login route.
// The GET request implements the hydra login flow
// and expects to receive a login_challenge as
// query parameter. It shows a login page if necessary
// or redirects the browser back to hydra.
// The POST request expects to receive a username
// and a password in the body. If the login with
// the given credentials is successful, a request
// to hydra's accept endpoint will be sent. Otherwise,
// the reject endpoint will be used.
func (srv Server) GetLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGet(w, r, srv)
		case "POST":
			handlePost(w, r, srv)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// get function for the login route
// It returns a login page if "skip" is not equal true for the login challenge.
// If skip is true and the user still exists, the login challenge gets accepted.
// If the user does not exist, the login challenge gets rejected.
func handleGet(w http.ResponseWriter, r *http.Request, srv Server) {
	c := r.URL.Query().Get("login_challenge")
	if c == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := srv.hydraclient.GetLoginRequest(c)
	if err != nil {
		log.Printf("error while querying login request from hydra: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// is skip is false, we need to show a login form
	if !body.GetSkip() {
		renderLoginForm(w, c, "")
		return
	}
	// when skip is set, only verify if subject is a valid userid
	// for the case the user was deleted
	_, err = srv.Database().FindUserByID(body.GetSubject())
	if err != nil {
		reject(
			w,
			r,
			srv,
			c,
			"user_notfound",
			fmt.Sprintf("unable to find user with id: %v", body.GetSubject()),
		)
		return
	}
	// TODO: implement disabled users
	accept(w, r, srv, c, body.GetSubject())
}

// handle login request
// The function receives a username and password, sent by a form.
// It verifies the login and does either an accept or a reject.
func handlePost(w http.ResponseWriter, r *http.Request, srv Server) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form in login post request: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	submit, challenge, username, password := r.FormValue("submit"), r.FormValue("challenge"), r.FormValue("username"), r.FormValue("password")
	if submit == "cancel" {
		reject(w, r, srv, challenge, "cancelled", "login was cancelled by the user")
		return
	}
	if username == "" || password == "" {
		renderLoginForm(w, challenge, "Username or Password not set.")
		return
	}
	u, err := srv.Database().FindUserByName(username)
	if err != nil {
		renderLoginForm(w, challenge, fmt.Sprintf("User '%s' not found.", username))
		return
	}
	err = u.ValidatePassword(password)
	if err != nil {
		renderLoginForm(w, challenge, fmt.Sprintf("Invalid password for user '%s'.", username))
		return
	}
	accept(w, r, srv, challenge, u.ID.Hex())

}

// accept the logon request
func accept(w http.ResponseWriter, r *http.Request, srv Server, challenge string, userID string) {
	body, err := srv.hydraclient.AcceptLoginRequest(challenge, true, 7200, userID)
	if err != nil {
		log.Printf("error while accepting login request: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	http.Redirect(w, r, body.GetRedirectTo(), http.StatusTemporaryRedirect)
}

// reject the logon request
func reject(w http.ResponseWriter, r *http.Request, srv Server, challenge string, errorID string, errorDescription string) {
	body, err := srv.hydraclient.RejectLoginRequest(challenge, errorID, errorDescription)
	if err != nil {
		log.Printf("error while rejecting login request: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	http.Redirect(w, r, body.GetRedirectTo(), http.StatusTemporaryRedirect)
}
