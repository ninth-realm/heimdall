package client

import (
	"context"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/ninth-realm/heimdall/store"
)

type Service struct {
	Repo store.Repository
}

func (s Service) ListClients(ctx context.Context) ([]store.Client, error) {
	return s.Repo.ListClients(store.QueryOptions{Ctx: ctx})
}

func (s Service) GetClient(ctx context.Context, id uuid.UUID) (store.Client, error) {
	return s.Repo.GetClientById(id, store.QueryOptions{Ctx: ctx})
}

func (s Service) CreateClient(ctx context.Context, client store.NewClient) (store.Client, error) {
	client = cleanNewClient(client)
	return store.RunUnitOfWork(ctx, s.Repo, func(txn *sqlx.Tx) (store.Client, error) {
		id, err := s.Repo.InsertClient(client, store.QueryOptions{Ctx: ctx, Txn: txn})
		if err != nil {
			return store.Client{}, err
		}

		return s.Repo.GetClientById(id, store.QueryOptions{Ctx: ctx, Txn: txn})
	})
}

func cleanNewClient(client store.NewClient) store.NewClient {
	client.Name = strings.TrimSpace(client.Name)

	return client
}

func (s Service) UpdateClient(ctx context.Context, id uuid.UUID, patch store.ClientPatch) (store.Client, error) {
	client, err := s.Repo.GetClientById(id, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return store.Client{}, err
	}

	client = patch.ApplyTo(client)

	err = s.Repo.SaveClient(client, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return store.Client{}, err
	}

	return s.Repo.GetClientById(id, store.QueryOptions{Ctx: ctx})
}

func (s Service) DeleteClient(ctx context.Context, id uuid.UUID) error {
	return s.Repo.DeleteClient(id, store.QueryOptions{Ctx: ctx})
}

func (s Service) ListClientAPIKeys(ctx context.Context, clientID uuid.UUID) ([]store.APIKey, error) {
	return s.Repo.ListClientAPIKeys(clientID, store.QueryOptions{Ctx: ctx})
}
