package http

import (
	"net/http"
)

func (s *Server) handleIndex() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, 200, "Hello, World!")
	})
}
