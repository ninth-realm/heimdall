package store

import (
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type NewUser struct {
	FirstName string
	LastName  string
	Email     string
	Password  *string
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
	ListUsers(opts QueryOptions) ([]User, error)
	GetUserById(id uuid.UUID, opts QueryOptions) (User, error)
	GetUserByEmail(email string, opts QueryOptions) (User, error)
	InsertUser(user NewUser, opts QueryOptions) (uuid.UUID, error)
	SaveUser(user User, opts QueryOptions) error
	DeleteUser(id uuid.UUID, opts QueryOptions) error
}
