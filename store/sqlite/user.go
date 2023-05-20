package sqlite

import (
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/ninth-realm/heimdall/store"
)

func (db DB) ListUsers(opts store.QueryOptions) ([]store.User, error) {
	const query = `
		SELECT
			id,
			first_name,
			last_name,
			created_at,
			updated_at
		FROM
			` + "`user`"

	var users []store.User
	err := db.querier(opts.Txn).SelectContext(opts.Context(), &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db DB) GetUserById(id uuid.UUID, opts store.QueryOptions) (store.User, error) {
	const query = `
		SELECT
			id,
			first_name,
			last_name,
			created_at,
			updated_at
		FROM
			` + "`user`" + `
		WHERE
			id = ?
	`

	var user store.User
	err := db.querier(opts.Txn).GetContext(opts.Context(), &user, query, id)

	if err != nil {
		return store.User{}, err
	}

	return user, nil
}

func (db DB) InsertUser(user store.NewUser, opts store.QueryOptions) (uuid.UUID, error) {
	const query = `
		INSERT INTO` + "`user`" + `
			(id, first_name, last_name)
		VALUES
			(?, ?, ?)
	`

	id := db.UUIDGenerator.GenerateUUID()
	_, err := db.querier(opts.Txn).ExecContext(
		opts.Context(),
		query,
		id,
		user.FirstName,
		user.LastName,
	)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (db DB) SaveUser(user store.User, opts store.QueryOptions) error {
	const query = `
		UPDATE` + "`user`" + `
		SET
			first_name = ?,
			last_name = ?
		WHERE
			id = ?
	`

	_, err := db.querier(opts.Txn).ExecContext(
		opts.Context(),
		query,
		user.FirstName,
		user.LastName,
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteUser(id uuid.UUID, opts store.QueryOptions) error {
	const query = `
		DELETE FROM` + "`user`" + `
		WHERE
			id = ?
	`

	res, err := db.querier(opts.Txn).ExecContext(opts.Ctx, query, id)
	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); n == 0 {
		return errors.New("user not found")
	} else if err != nil {
		return err
	}

	return nil

}
