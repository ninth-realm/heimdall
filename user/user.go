package user

import (
	"context"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/ninth-realm/heimdall/store"
)

type Service struct {
	Repo store.Repository
}

func (s Service) ListUsers(ctx context.Context) ([]store.User, error) {
	return s.Repo.ListUsers(ctx)
}

func (s Service) GetUser(ctx context.Context, id uuid.UUID) (store.User, error) {
	return s.Repo.GetUserById(ctx, id)
}

func (s Service) CreateUser(ctx context.Context, user store.NewUser) (store.User, error) {
	user = cleanNewUser(user)

	id, err := s.Repo.InsertUser(ctx, user)
	if err != nil {
		return store.User{}, err
	}

	return s.Repo.GetUserById(ctx, id)
}

func cleanNewUser(user store.NewUser) store.NewUser {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	return user
}

func (s Service) UpdateUser(ctx context.Context, id uuid.UUID, patch store.UserPatch) (store.User, error) {
	user, err := s.Repo.GetUserById(ctx, id)
	if err != nil {
		return store.User{}, err
	}

	user = patch.ApplyTo(user)

	err = s.Repo.SaveUser(ctx, user)
	if err != nil {
		return store.User{}, err
	}

	return s.Repo.GetUserById(ctx, id)
}

func (s Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.Repo.DeleteUser(ctx, id)
}
