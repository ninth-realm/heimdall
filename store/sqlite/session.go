package sqlite

import (
	"errors"
	"time"

	"github.com/ninth-realm/heimdall/store"
)

func (db DB) GetSession(token string, opts store.QueryOptions) (store.Session, error) {
	const query = `
		SELECT
			token,
			user_id,
			created_at,
			expires_at
		FROM
			session
        WHERE
            token = ?
	`

	var session store.Session
	err := db.querier(opts.Txn).GetContext(opts.Context(), &session, query, token)
	if err != nil {
		return store.Session{}, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return store.Session{}, errors.New("session expired")
	}

	return session, nil
}
func (db DB) SaveSession(session store.Session, opts store.QueryOptions) error {
	const query = `
        INSERT INTO session
            (token, user_id, created_at, expires_at)
        VALUES
            (?, ?, ?, ?)
    `

	_, err := db.querier(opts.Txn).ExecContext(
		opts.Context(),
		query,
		session.Token,
		session.UserId,
		session.CreatedAt,
		session.ExpiresAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteSession(token string, opts store.QueryOptions) error {
	const query = `
        DELETE FROM
            session
        WHERE
            token = ?
    `

	res, err := db.querier(opts.Txn).ExecContext(opts.Context(), query, token)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return store.NotFoundError{ResourceType: "session", ResourceID: token}
	}

	return nil
}
