package sqlite

import (
	"github.com/gofrs/uuid/v5"
	"github.com/ninth-realm/heimdall/store"
)

func (db DB) InsertPassword(password store.NewPassword, opts store.QueryOptions) (uuid.UUID, error) {
	const query = `
		INSERT INTO password
			(id, user_id, hash)
		VALUES
			(?, ?, ?)
	`

	id := db.UUIDGenerator.GenerateUUID()
	_, err := db.querier(opts.Txn).ExecContext(
		opts.Context(),
		query,
		id,
		password.UserID,
		password.Hash,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
