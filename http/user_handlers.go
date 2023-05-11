package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/ninth-realm/heimdall/store"
)

func (s *Server) handleUsersList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := s.UserService.ListUsers(r.Context())
		if err != nil {
			s.respondWithError(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, users)
	})
}

func (s *Server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		FirstName nonEmptyString `json:"firstName"`
		LastName  nonEmptyString `json:"lastName"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody request
		err := s.decode(r, &requestBody)
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := s.UserService.CreateUser(r.Context(), store.NewUser{
			FirstName: string(requestBody.FirstName),
			LastName:  string(requestBody.LastName),
		})
		if err != nil {
			s.respondWithError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, user)
	})
}

func (s *Server) handleUsersGet() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(chi.URLParamFromCtx(r.Context(), "userID"))
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := s.UserService.GetUser(r.Context(), id)
		if err != nil {
			s.respondWithError(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, user)
	})
}

func (s *Server) handleUsersUpdate() http.HandlerFunc {
	type request struct {
		FirstName *nonEmptyString `json:"firstName"`
		LastName  *nonEmptyString `json:"lastName"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(chi.URLParamFromCtx(r.Context(), "userID"))
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

		user, err := s.UserService.UpdateUser(r.Context(), id, store.UserPatch{
			FirstName: (*string)(requestBody.FirstName),
			LastName:  (*string)(requestBody.LastName),
		})
		if err != nil {
			s.respondWithError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, user)
	})
}

func (s *Server) handleUsersDelete() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(chi.URLParamFromCtx(r.Context(), "userID"))
		if err != nil {
			s.respondWithError(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.UserService.DeleteUser(r.Context(), id)
		if err != nil {
			s.respondWithError(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusNoContent, nil)
	})
}
