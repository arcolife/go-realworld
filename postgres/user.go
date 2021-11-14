package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"

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

func (us *UserService) UserByID(ctx context.Context, id uint) (*conduit.User, error) {
	return nil, nil
}

func (us *UserService) UserByEmail(ctx context.Context, email string) (*conduit.User, error) {
	tx, err := us.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	user, err := findUserByEmail(ctx, tx, email)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Users(ctx context.Context, uf conduit.UserFilter) ([]*conduit.User, error) {
	return nil, nil
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

func (us *UserService) UpdateUser(ctx context.Context, user *conduit.User, patch conduit.UserPatch) error {
	tx, err := us.db.BeginTxx(ctx, nil)

	if err != nil {
		log.Println(err)
		return conduit.ErrInternal
	}

	defer tx.Rollback()

	if err := updateUser(ctx, tx, user, patch); err != nil {
		log.Println(err)
		return conduit.ErrInternal
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return conduit.ErrInternal
	}

	return nil
}

func (us *UserService) DeleteUser(ctx context.Context, id uint) error {
	return nil
}

func createUser(ctx context.Context, tx *sqlx.Tx, user *conduit.User) error {
	// Execute insertion query.
	query := `
	INSERT INTO users (email, username, bio, image, password_hash)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`
	args := []interface{}{user.Email, user.Username, user.Bio, user.Image, user.PasswordHash}
	err := tx.QueryRowxContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

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

	return nil
}

// FormatLimitOffset returns a SQL string for a given limit & offset.
// Clauses are only added if limit and/or offset are greater than zero.
func FormatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	} else if offset > 0 {
		return fmt.Sprintf("OFFSET %d", offset)
	}
	return ""
}

// findUserByEmail is a helper function to fetch a user by email.
func findUserByEmail(ctx context.Context, tx *sqlx.Tx, email string) (*conduit.User, error) {
	us, err := findUsers(ctx, tx, conduit.UserFilter{Email: &email})

	if err != nil {
		return nil, err
	} else if len(us) == 0 {
		return nil, conduit.ErrNotFound
	}
	return us[0], nil
}

func findUserByID(ctx context.Context, tx *sqlx.Tx, id uint) (*conduit.User, error) {
	us, err := findUsers(ctx, tx, conduit.UserFilter{ID: &id})

	if err != nil {
		return nil, err
	} else if len(us) == 0 {
		return nil, conduit.ErrNotFound
	}
	return us[0], nil
}

func findUsers(ctx context.Context, tx *sqlx.Tx, filter conduit.UserFilter) ([]*conduit.User, error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Email; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("email = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Username; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("username = $%d", argPosition)), append(args, *v)
	}

	// Execute query to fetch user rows.
	query := "SELECT * from users WHERE " + strings.Join(where, " AND ") + " ORDER BY id ASC" +
		FormatLimitOffset(filter.Limit, filter.Offset)
	rows, err := tx.QueryxContext(ctx, query, args...)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	// Deserialize rows into User objects.
	users := make([]*conduit.User, 0)

	for rows.Next() {
		var user conduit.User

		if err := rows.StructScan(&user); err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

func updateUser(ctx context.Context, tx *sqlx.Tx, user *conduit.User, patch conduit.UserPatch) error {
	if v := patch.Bio; v != nil {
		user.Bio = *v
	}

	if v := patch.Email; v != nil {
		user.Email = *v
	}

	if v := patch.PasswordHash; v != nil {
		user.PasswordHash = *v
	}

	if v := patch.Image; v != nil {
		user.Image = *v
	}

	if v := patch.Username; v != nil {
		user.Image = *v
	}

	args := []interface{}{
		user.Username,
		user.Email,
		user.Bio,
		user.Image,
		user.PasswordHash,
		user.ID,
	}

	query := `
	UPDATE users 
	SET username = $1, email = $2, bio = $3, image = $4, password_hash = $5, updated_at = NOW()
	WHERE id = $6
	RETURNING updated_at`

	if err := tx.QueryRowxContext(ctx, query, args...).Scan(&user.UpdatedAt); err != nil {
		log.Printf("error updating record: %v", err)
		return conduit.ErrInternal
	}

	return nil
}
