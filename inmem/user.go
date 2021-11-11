package inmem

import (
	"context"
	"errors"
	"time"

	"github.com/0xdod/go-realworld/conduit"
)

type UserService struct {
	users []*conduit.User
}

// type UserService2 map[string]*conduit.User

func (us *UserService) CreateUser(_ context.Context, user *conduit.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	us.users = append(us.users, user)
	return nil
}

func (us *UserService) Authenticate(ctx context.Context, email, password string) (*conduit.User, error) {
	user, _ := us.UserByEmail(ctx, email)
	if user.VerifyPassword(password) {
		return user, nil
	}
	return nil, errors.New("invalid credentials")
}

func (us *UserService) UserByEmail(_ context.Context, email string) (*conduit.User, error) {
	for _, user := range us.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("not found")
}

func (us *UserService) Users(_ context.Context, uf conduit.UserFilter) ([]*conduit.User, error) {
	return nil, nil
}

func (us *UserService) UserByID(_ context.Context, id uint) (*conduit.User, error) {
	for _, user := range us.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("not found")
}
