package http

func (s *Server) loadRoutes() {
	s.Router.With(s.authenticateRoute).Get("/api/v1/users", s.handleUsersList())
	s.Router.With(s.authenticateRoute).Post("/api/v1/users", s.handleUsersCreate())
	s.Router.With(s.authenticateRoute).Get("/api/v1/users/{userID}", s.handleUsersGet())
	s.Router.With(s.authenticateRoute).Patch("/api/v1/users/{userID}", s.handleUsersUpdate())
	s.Router.With(s.authenticateRoute).Delete("/api/v1/users/{userID}", s.handleUsersDelete())

	s.Router.With(s.authenticateRoute).Get("/api/v1/clients", s.handleClientsList())
	s.Router.With(s.authenticateRoute).Post("/api/v1/clients", s.handleClientsCreate())
	s.Router.With(s.authenticateRoute).Get("/api/v1/clients/{clientID}", s.handleClientsGet())
	s.Router.With(s.authenticateRoute).Patch("/api/v1/clients/{clientID}", s.handleClientsUpdate())
	s.Router.With(s.authenticateRoute).Delete("/api/v1/clients/{clientID}", s.handleClientsDelete())
	s.Router.With(s.authenticateRoute).Get("/api/v1/clients/{clientID}/api-keys", s.handleClientsAPIKeysGet())
	s.Router.With(s.authenticateRoute).Post("/api/v1/clients/{clientID}/api-keys", s.handleClientsAPIKeysCreate())
	s.Router.With(s.authenticateRoute).Delete("/api/v1/clients/{clientID}/api-keys/{keyID}", s.handleClientsAPIKeysDelete())

	s.Router.Post("/api/v1/auth/login", s.handleAuthLogin())
	s.Router.Post("/api/v1/auth/logout", s.handleAuthLogout())
	s.Router.With(s.authenticateRoute).Post("/api/v1/auth/introspect", s.handleAuthIntrospect())
}
