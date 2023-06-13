package store

import "github.com/gofrs/uuid/v5"

type AuthRepository interface {
	GetUserPasswordHash(userID uuid.UUID, opts QueryOptions) (string, error)
}
