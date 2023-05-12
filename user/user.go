package user

import (
	"context"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/ninth-realm/heimdall/crypto"
	"github.com/ninth-realm/heimdall/store"
)

type Service struct {
	Repo store.Repository
}

func (s Service) ListUsers(ctx context.Context) ([]store.User, error) {
	return s.Repo.ListUsers(store.QueryOptions{Ctx: ctx})
}

func (s Service) GetUser(ctx context.Context, id uuid.UUID) (store.User, error) {
	return s.Repo.GetUserById(id, store.QueryOptions{Ctx: ctx})
}

func (s Service) CreateUser(ctx context.Context, user store.NewUser) (store.User, error) {
	user = cleanNewUser(user)
	return store.RunUnitOfWork(ctx, s.Repo, func(txn *sqlx.Tx) (store.User, error) {
		id, err := s.Repo.InsertUser(user, store.QueryOptions{Ctx: ctx, Txn: txn})
		if err != nil {
			return store.User{}, err
		}

		_, err = s.Repo.InsertEmail(
			store.NewEmail{
				UserID: id,
				Email:  user.Email,
			},
			store.QueryOptions{Ctx: ctx, Txn: txn},
		)
		if err != nil {
			return store.User{}, err
		}

		if user.Password != nil {
			hash, err := crypto.GetPasswordHash(*user.Password, crypto.DefaultParams)
			if err != nil {
				return store.User{}, err
			}

			_, err = s.Repo.InsertPassword(
				store.NewPassword{
					Hash:   hash,
					UserID: id,
				},
				store.QueryOptions{Ctx: ctx, Txn: txn},
			)

			if err != nil {
				return store.User{}, err
			}
		}

		return s.Repo.GetUserById(id, store.QueryOptions{Ctx: ctx, Txn: txn})
	})
}

func cleanNewUser(user store.NewUser) store.NewUser {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	return user
}

func (s Service) UpdateUser(ctx context.Context, id uuid.UUID, patch store.UserPatch) (store.User, error) {
	user, err := s.Repo.GetUserById(id, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return store.User{}, err
	}

	user = patch.ApplyTo(user)

	err = s.Repo.SaveUser(user, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return store.User{}, err
	}

	return s.Repo.GetUserById(id, store.QueryOptions{Ctx: ctx})
}

func (s Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.Repo.DeleteUser(id, store.QueryOptions{Ctx: ctx})
}
