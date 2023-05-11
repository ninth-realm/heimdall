package store

import (
	"context"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type NewUser struct {
	FirstName string
	LastName  string
}

type UserPatch struct {
	FirstName *string
	LastName  *string
}

func (p UserPatch) ApplyTo(user User) User {
	if p.FirstName != nil {
		user.FirstName = strings.TrimSpace(*p.FirstName)
	}

	if p.LastName != nil {
		user.LastName = strings.TrimSpace(*p.LastName)
	}

	return user
}

type UserRepository interface {
	ListUsers(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (User, error)
	InsertUser(ctx context.Context, user NewUser) (uuid.UUID, error)
	SaveUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
