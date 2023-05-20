package http

func (s *Server) loadRoutes() {
	s.Router.Get("/api/v1/users", s.handleUsersList())
	s.Router.Post("/api/v1/users", s.handleUsersCreate())
	s.Router.Get("/api/v1/users/{userID}", s.handleUsersGet())
	s.Router.Patch("/api/v1/users/{userID}", s.handleUsersUpdate())
	s.Router.Delete("/api/v1/users/{userID}", s.handleUsersDelete())
}
