package sqlite

import (
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/ninth-realm/heimdall/store"
)

func (db DB) ListClients(opts store.QueryOptions) ([]store.Client, error) {
	const query = `
		SELECT
			id,
			name,
			enabled,
			created_at,
			updated_at
		FROM
			client
	`

	var clients []store.Client
	err := db.querier(opts.Txn).SelectContext(opts.Context(), &clients, query)
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (db DB) GetClientById(id uuid.UUID, opts store.QueryOptions) (store.Client, error) {
	const query = `
		SELECT
			id,
			name,
			enabled,
			created_at,
			updated_at
		FROM
			client
		WHERE
			id = ?
	`

	var client store.Client
	err := db.querier(opts.Txn).GetContext(opts.Context(), &client, query, id)

	if err != nil {
		return store.Client{}, err
	}

	return client, nil
}

func (db DB) InsertClient(client store.NewClient, opts store.QueryOptions) (uuid.UUID, error) {
	const query = `
		INSERT INTO client
			(id, name, enabled)
		VALUES
			(?, ?, ?)
	`

	id := db.UUIDGenerator.GenerateUUID()
	_, err := db.querier(opts.Txn).ExecContext(
		opts.Context(),
		query,
		id,
		client.Name,
		client.Enabled,
	)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (db DB) SaveClient(client store.Client, opts store.QueryOptions) error {
	const query = `
		UPDATE client
		SET
			name = ?,
			enabled = ?
		WHERE
			id = ?
	`

	_, err := db.querier(opts.Txn).ExecContext(
		opts.Context(),
		query,
		client.Name,
		client.Enabled,
		client.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteClient(id uuid.UUID, opts store.QueryOptions) error {
	const query = `
		DELETE FROM client
		WHERE
			id = ?
	`

	res, err := db.querier(opts.Txn).ExecContext(opts.Ctx, query, id)
	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); n == 0 {
		return errors.New("client not found")
	} else if err != nil {
		return err
	}

	return nil

}

func (db DB) ListClientAPIKeys(clientID uuid.UUID, opts store.QueryOptions) ([]store.APIKey, error) {
	const query = `
		SELECT
			id,
			client_id,
			description,
			prefix,
			hash,
			created_at,
			updated_at
		FROM
			api_key
		WHERE
			client_id = ?
	`

	keys := []store.APIKey{}
	err := db.querier(opts.Txn).SelectContext(opts.Context(), &keys, query, clientID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
