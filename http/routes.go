package http

func (s *Server) loadRoutes() {
	s.Router.Get("/api/v1/users", s.handleUsersList())
	s.Router.Post("/api/v1/users", s.handleUsersCreate())
	s.Router.Get("/api/v1/users/{userID}", s.handleUsersGet())
	s.Router.Patch("/api/v1/users/{userID}", s.handleUsersUpdate())
	s.Router.Delete("/api/v1/users/{userID}", s.handleUsersDelete())

	s.Router.Get("/api/v1/clients", s.handleClientsList())
	s.Router.Post("/api/v1/clients", s.handleClientsCreate())
	s.Router.Get("/api/v1/clients/{clientID}", s.handleClientsGet())
	s.Router.Patch("/api/v1/clients/{clientID}", s.handleClientsUpdate())
	s.Router.Delete("/api/v1/clients/{clientID}", s.handleClientsDelete())
	s.Router.Get("/api/v1/clients/{clientID}/api-keys", s.handleClientsAPIKeysGet())
	s.Router.Post("/api/v1/clients/{clientID}/api-keys", s.handleClientsAPIKeysCreate())

	s.Router.Post("/api/v1/auth/login", s.handleAuthLogin())
}
