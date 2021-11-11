package postgres

import (
	"context"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/jmoiron/sqlx"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db}
}

func (us *UserService) CreateUser(ctx context.Context, user *conduit.User) error {
	tx, err := us.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := createUser(ctx, tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

func (us *UserService) Authenticate(ctx context.Context, email, password string) (*conduit.User, error) {
	user, err := us.UserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if !user.VerifyPassword(password) {
		return nil, conduit.ErrUnAuthorized
	}

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

func createUser(ctx context.Context, tx *sqlx.Tx, user *conduit.User) error {
	// Execute insertion query.
	query := `
	INSERT INTO users (email, username, bio, image, password_hash)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
	`
	args := []interface{}{user.Email, user.Username, user.Bio, user.Image, user.PasswordHash}
	err := tx.QueryRowxContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return conduit.ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return conduit.ErrDuplicateUsername
		default:
			return err
		}
	}

	query = `UPDATE users SET updated_at = $1 where id = $2`

	err = tx.QueryRowxContext(ctx, query, user.CreatedAt, user.ID).Err()

	if err != nil {
		return err
	}

	return nil
}
