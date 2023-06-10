package auth

import (
	"context"
	"errors"

	"github.com/ninth-realm/heimdall/crypto"
	"github.com/ninth-realm/heimdall/store"
)

type Service struct {
	Repo store.Repository
}

func (s Service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.Repo.GetUserByEmail(username, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return "", err
	}

	hash, err := s.Repo.GetUserPasswordHash(user.ID, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return "", err
	}

	correctPassword, err := crypto.ValidatePassword(password, hash)
	if err != nil {
		return "", err
	} else if !correctPassword {
		return "", errors.New("incorrect password")
	}

	token, err := generateJWT(JWTSettings{
		Issuer:     "heimdall",
		Lifespan:   900,
		SigningKey: "so-secret-wow",
		Algorithm:  HMAC256Algorithm,
	})
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
