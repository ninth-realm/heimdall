package sqlite

import (
	"github.com/gofrs/uuid/v5"
	"github.com/ninth-realm/heimdall/store"
)

func (db DB) InsertEmail(email store.NewEmail, opts store.QueryOptions) (uuid.UUID, error) {
	const query = `
		INSERT INTO email
			(id, user_id, email)
		VALUES
			(?, ?, ?)
	`

	id := db.UUIDGenerator.GenerateUUID()
	_, err := db.querier(opts.Txn).ExecContext(
		opts.Context(),
		query,
		id,
		email.UserID,
		email.Email,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
