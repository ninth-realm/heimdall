package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mattmeyers/level"
)

type Server struct {
	Logger level.Logger

	Router chi.Router
}

// NewServer builds a new server object with the default middleware and router
// already configured. This is the typical way that `Server`s should be created
// to ensure they have everything needed to function.
func NewServer() *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	s := &Server{Router: r}
	s.loadRoutes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// ListenAndServe spins up the server to listen on the provided address. The
// allowed addresses follow the same rules as `http.ListenAndServe`.
func (s *Server) ListenAndServe(addr string) error {
	if s.Logger == nil {
		s.Logger, _ = level.NewBasicLogger(level.Info, nil)
	}

	s.Logger.Info("Server listening on %s", addr)
	return http.ListenAndServe(addr, s)
}

// The max size in bytes of a request body. 5 mB should be plenty.
const requestBodyLimit = 5_242_880

// Reads a request body into the provided destination. It is expected that the
// `dst` parameter is a pointer and that the request body is JSON.
func (s *Server) decode(r *http.Request, dst any) error {
	return json.NewDecoder(io.LimitReader(r.Body, requestBodyLimit)).Decode(dst)
}

// Writes a response to the client. All responses are written in an envelope with
// the provided data under the `response` key. At this point, all responses are
// encoded as JSON.
//
// This method does not return any errors. If an error occurs during the encoding
// process, the error is logged and an internal server error is returned to the
// client.
func (s *Server) respond(w http.ResponseWriter, r *http.Request, status int, data any) {
	type envelope struct {
		Response any `json:"response"`
	}

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(envelope{Response: data})
	if err != nil {
		s.respondWithError(w, r, http.StatusInternalServerError, err)
	}
}

// Writes an error response to the client. If the provided status code is a 5xx
// code, then the error is logged, and the appropriate status text is used in
// the response. We do not want to leak any internal details to the client.
func (s *Server) respondWithError(w http.ResponseWriter, r *http.Request, status int, data error) {
	type envelope struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}

	if status >= 500 {
		s.logError(r, data)
		data = errors.New(http.StatusText(status))
	}

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(envelope{Code: status, Error: data.Error()})
	if err != nil {
		// If we get here, something is seriously wrong. Log and move on.
		s.logError(r, data)
	}

}

func (s *Server) logError(r *http.Request, err error) {
	s.Logger.Error("%s: %v", r.Context().Value(middleware.RequestIDKey), err)
}
