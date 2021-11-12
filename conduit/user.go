package conduit

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint      `json:"id,omitempty"`
	Email        string    `json:"email,omitempty"`
	Username     string    `json:"username,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	Image        string    `json:"image,omitempty"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

var AnonymousUser User

type UserFilter struct {
	ID       *uint   `json:"id,omitempty"`
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`

	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

func (u *User) SetPassword(password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		// return better error message
		return err
	}

	u.PasswordHash = string(hashBytes)

	return nil
}

func (u User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

func (u *User) IsAnonymous() bool {
	return u == &AnonymousUser
}

type UserService interface {
	CreateUser(context.Context, *User) error

	//UserByID(context.Context, uint) (*User, error)

	UserByEmail(context.Context, string) (*User, error)

	//Users(context.Context, UserFilter) ([]*User, error)

	Authenticate(ctx context.Context, email, password string) (*User, error)
}
