package http

import (
	"errors"
	"net/http"
)

const APIKeyHeaderName = "X-API-Key"

const SessionCookieName = "heimdall_sessionToken"

var authErr = errors.New("missing or invalid auth token")

func (s *Server) authenticateRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.DisableAuth {
			next.ServeHTTP(w, r)
			return
		}

		err := s.authenticateAPIKey(r)
		if err == nil {
			next.ServeHTTP(w, r)
			return
		}

		err = s.authenticateSessionToken(r)
		if err == nil {
			next.ServeHTTP(w, r)
			return
		}

		s.respondWithError(w, r, http.StatusUnauthorized, authErr)
	})
}

func (s *Server) authenticateAPIKey(r *http.Request) error {
	token := r.Header.Get(APIKeyHeaderName)
	if token == "" {
		return authErr
	}

	return s.AuthService.ValidateAPIKey(r.Context(), token)
}

func (s *Server) authenticateSessionToken(r *http.Request) error {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return authErr
	}

	_, err = s.AuthService.IntrospectToken(r.Context(), cookie.Value)
	return err
}
