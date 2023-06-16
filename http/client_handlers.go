package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/ninth-realm/heimdall/store"
)

func (s *Server) handleClientsList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clients, err := s.ClientService.ListClients(r.Context())
		if err != nil {
			s.respondWithError(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, clients)
	})
}

func (s *Server) handleClientsCreate() http.HandlerFunc {
	type request struct {
		Name    nonEmptyString `json:"name"`
		Enabled bool           `json:"enabled"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody request
		err := s.decode(r, &requestBody)
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		client, err := s.ClientService.CreateClient(r.Context(), store.NewClient{
			Name:    requestBody.Name.toString(),
			Enabled: requestBody.Enabled,
		})
		if err != nil {
			s.respondWithError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, client)
	})
}

func (s *Server) handleClientsGet() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(chi.URLParamFromCtx(r.Context(), "clientID"))
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		client, err := s.ClientService.GetClient(r.Context(), id)
		if err != nil {
			s.respondWithError(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, client)
	})
}

func (s *Server) handleClientsUpdate() http.HandlerFunc {
	type request struct {
		Name    *nonEmptyString `json:"name"`
		Enabled *bool           `json:"enabled"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(chi.URLParamFromCtx(r.Context(), "clientID"))
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		var requestBody request
		err = s.decode(r, &requestBody)
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		client, err := s.ClientService.UpdateClient(r.Context(), id, store.ClientPatch{
			Name:    (*string)(requestBody.Name),
			Enabled: requestBody.Enabled,
		})
		if err != nil {
			s.respondWithError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, client)
	})
}

func (s *Server) handleClientsDelete() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(chi.URLParamFromCtx(r.Context(), "clientID"))
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.ClientService.DeleteClient(r.Context(), id)
		if err != nil {
			s.respondWithError(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusNoContent, nil)
	})
}

func (s *Server) handleClientsAPIKeysGet() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(chi.URLParamFromCtx(r.Context(), "clientID"))
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		keys, err := s.ClientService.ListClientAPIKeys(r.Context(), id)
		if err != nil {
			s.respondWithError(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, keys)
	})
}

func (s *Server) handleClientsAPIKeysCreate() http.HandlerFunc {
	type request struct {
		Description *string `json:"description"`
	}

	type response struct {
		Key string `json:"key"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(chi.URLParamFromCtx(r.Context(), "clientID"))
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		var body request
		if err = s.decode(r, &body); err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		keys, err := s.ClientService.GenerateAPIKey(r.Context(), store.NewAPIKey{ClientID: id, Description: body.Description})
		if err != nil {
			s.respondWithError(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, keys)
	})
}
