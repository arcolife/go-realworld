package postgres

import (
	"context"

	"github.com/0xdod/go-realworld/conduit"
)

type UserService struct {
}

func NewUserService() *UserService {
	return nil
}

func (us *UserService) CreateUser(ctx context.Context, user *conduit.User) error {
	return nil
}

// func (us *UserService) Authenticate(_ context.Context, user *conduit.User) (*conduit.User, error) {
// 	user.CreatedAt = time.Now()
// 	user.UpdatedAt = time.Now()
// 	us.users = append(us.users, user)
// 	return nil
// }

func (us *UserService) UserByEmail(ctx context.Context, email string) (*conduit.User, error) {
	return nil, nil
}

func (us *UserService) Users(ctx context.Context, uf conduit.UserFilter) ([]*conduit.User, error) {
	return nil, nil
}

func (us *UserService) UserByID(ctx context.Context, id uint) (*conduit.User, error) {
	return nil, nil
}
