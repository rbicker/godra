package godra

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rbicker/godra/internal/db"
	"github.com/rbicker/godra/internal/hydraclient"
	"github.com/rbicker/nogo"
)

// Server represents an api server.
type Server struct {
	port            int
	httpServer      *http.Server
	db              db.Database
	hydraPrivateURL string
	hydraclient     hydraclient.Client
}

// NewServer creates a new api server.
// It takes functional parameters to change default options
// such as the api port
// It returns the newly created server or an error if
// something went wrong.
func NewServer(opts ...func(*Server) error) (*Server, error) {
	// create server with default options
	var srv = Server{
		port:            5000,
		hydraPrivateURL: "http://127.0.0.1:4445",
	}
	// run functional options
	for _, op := range opts {
		err := op(&srv)
		if err != nil {
			return nil, fmt.Errorf("setting server option failed: %w", err)
		}
	}
	if srv.db == nil {
		db, err := db.NewMongoConnection()
		if err != nil {
			return nil, fmt.Errorf("creating new mongo db connection failed: %w", err)
		}
		srv.db = db
	}
	return &srv, nil
}

// Serve starts the http server.
func (srv Server) Serve() error {
	m := http.NewServeMux()
	if static, ok := os.LookupEnv("CUSTOM_STATIC_PATH"); ok {
		m.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(static))))
	}
	m.Handle("/public/", http.StripPrefix("/public/", http.FileServer(nogo.Dir("/assets/public"))))
	m.HandleFunc("/login", srv.GetLoginHandler())
	m.HandleFunc("/consent", srv.GetConsentHandler())
	m.HandleFunc("/logout", srv.GetLogoutHandler())
	srv.httpServer = &http.Server{Addr: fmt.Sprintf(":%v", srv.port), Handler: m}
	return srv.httpServer.ListenAndServe()
}

// Shutdown stops the http server gracefully.
func (srv Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.httpServer.Shutdown(ctx)
}

// Database returns the database connection.
func (srv Server) Database() db.Database {
	return srv.db
}

// HydraPrivateURL returns the private url of the hydra server.
func (srv Server) HydraPrivateURL() string {
	return srv.hydraPrivateURL
}

// SetPort changes the port on which the api server listens.
// The default port is 5000.
// SetPort returns an error if an invalid port is given.
func SetPort(port int) func(*Server) error {
	return func(srv *Server) error {
		if port <= 0 {
			return fmt.Errorf("invalid port number: %v", port)
		}
		srv.port = port
		return nil
	}
}

// SetDatabase sets the given db provider.
// A mongodb provider will be used by default
func SetDatabase(db db.Database) func(*Server) error {
	return func(srv *Server) error {
		srv.db = db
		return nil
	}
}

// SetHydraClient sets the given hydra client.
func SetHydraClient(client hydraclient.Client) func(*Server) error {
	return func(srv *Server) error {
		srv.hydraclient = client
		return nil
	}
}
