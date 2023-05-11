package sqlite

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/ninth-realm/heimdall/store"
)

func (db DB) ListUsers(ctx context.Context) ([]store.User, error) {

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
	err := db.Conn.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db DB) GetUserById(ctx context.Context, id uuid.UUID) (store.User, error) {
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
	err := db.Conn.GetContext(ctx, &user, query, id)
	if err != nil {
		return store.User{}, err
	}

	return user, nil
}

func (db DB) InsertUser(ctx context.Context, user store.NewUser) (uuid.UUID, error) {
	const query = `
		INSERT INTO` + "`user`" + `
			(id, first_name, last_name)
		VALUES
			(?, ?, ?)
	`

	id := db.UUIDGenerator.GenerateUUID()
	_, err := db.Conn.ExecContext(ctx, query, id, user.FirstName, user.LastName)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (db DB) SaveUser(ctx context.Context, user store.User) error {
	const query = `
		UPDATE` + "`user`" + `
		SET
			first_name = ?,
			last_name = ?
		WHERE
			id = ?
	`

	_, err := db.Conn.ExecContext(ctx, query, user.FirstName, user.LastName, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM` + "`user`" + `
		WHERE
			id = ?
	`

	res, err := db.Conn.ExecContext(ctx, query, id)
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
