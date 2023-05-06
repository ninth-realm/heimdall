package http

func (s *Server) loadRoutes() {
	s.Router.Get("/", s.handleIndex())
}
