package http

import "net/http"

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

		s.respond(w, r, http.StatusOK, token)
	})
}
