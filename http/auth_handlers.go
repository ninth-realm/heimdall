package http

import (
	"net/http"

	"github.com/ninth-realm/heimdall/auth"
)

func (s *Server) handleAuthLogin() http.HandlerFunc {
	type request struct {
		Username nonEmptyString `json:"username"`
		Password nonEmptyString `json:"password"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody request
		err := s.decode(r, &requestBody)
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		token, err := s.AuthService.Login(
			r.Context(),
			requestBody.Username.toString(),
			requestBody.Password.toString(),
		)
		if err != nil {
			s.respondWithError(w, r, http.StatusUnauthorized, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     SessionCookieName,
			Value:    token.AccessToken,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   token.Lifespan,
		})

		s.respond(w, r, http.StatusNoContent, nil)
	})
}

func (s *Server) handleAuthLogout() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(SessionCookieName)
		if err != nil {
			// No session cookie means nothing to do.
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		err = s.AuthService.Logout(r.Context(), token.Value)
		if err != nil {
			s.respondWithError(w, r, http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     SessionCookieName,
			Value:    "",
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			MaxAge:   -1,
		})

		s.respond(w, r, http.StatusNoContent, nil)
	})
}

func (s *Server) handleAuthIntrospect() http.HandlerFunc {
	type request struct {
		Token string `json:"token"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body request
		if err := s.decode(r, &body); err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		token, err := s.AuthService.IntrospectToken(r.Context(), body.Token)
		if err != nil {
			token = auth.TokenInfo{Active: false}
			s.respond(w, r, http.StatusUnauthorized, token)
			return
		}

		s.respond(w, r, http.StatusOK, token)
	})
}
