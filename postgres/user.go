package postgres

import (
	"context"
	"errors"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret = []byte("sample-secret")

type UserService struct {
}

func NewUserService() *UserService {
	return nil
}

func (us *UserService) CreateUser(ctx context.Context, user *conduit.User) error {
	return nil
}

func (us *UserService) Authenticate(ctx context.Context, email, password string) (*conduit.User, error) {
	user, err := us.UserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if !user.VerifyPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})

	tokenString, err := token.SignedString(hmacSampleSecret)

	if err != nil {
		return nil, err
	}

	user.Token = tokenString

	// update user

	return user, nil
}

func (us *UserService) UserByEmail(ctx context.Context, email string) (*conduit.User, error) {
	return nil, nil
}

func (us *UserService) Users(ctx context.Context, uf conduit.UserFilter) ([]*conduit.User, error) {
	return nil, nil
}

func (us *UserService) UserByID(ctx context.Context, id uint) (*conduit.User, error) {
	return nil, nil
}
