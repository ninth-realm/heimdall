package http

import (
	"fmt"
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
			MaxAge:   token.Lifespan,
		})

		s.respond(w, r, http.StatusOK, token)
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
			fmt.Println(err.Error())
			return
		}

		s.respond(w, r, http.StatusOK, token)
	})
}
