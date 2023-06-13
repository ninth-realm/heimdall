package auth

import (
	"context"
	"errors"

	"github.com/ninth-realm/heimdall/crypto"
	"github.com/ninth-realm/heimdall/store"
)

type Service struct {
	Repo store.Repository

	JWTSettings JWTSettings
}

func (s Service) Login(ctx context.Context, username, password string) (Token, error) {
	user, err := s.Repo.GetUserByEmail(username, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return Token{}, err
	}

	hash, err := s.Repo.GetUserPasswordHash(user.ID, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return Token{}, err
	}

	correctPassword, err := crypto.ValidatePassword(password, hash)
	if err != nil {
		return Token{}, err
	} else if !correctPassword {
		return Token{}, errors.New("incorrect password")
	}

	token, err := generateJWT(s.JWTSettings)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}
