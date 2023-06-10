package sqlite

import (
	"github.com/gofrs/uuid/v5"
	"github.com/ninth-realm/heimdall/store"
)

func (db DB) GetUserPasswordHash(userID uuid.UUID, opts store.QueryOptions) (string, error) {
	const query = `
		SELECT
			hash
		FROM
			password
		WHERE
			user_id = ?
	`

	var hash string
	err := db.querier(opts.Txn).GetContext(opts.Context(), &hash, query, userID)
	if err != nil {
		return "", err
	}

	return hash, nil
}
