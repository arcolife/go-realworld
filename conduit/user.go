package conduit

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint      `json:"-"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Bio          string    `json:"bio"`
	Image        string    `json:"image"`
	Token        string    `json:"token"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"-" db:"updated_at"`
}

type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}


func (u *User) Profile() *Profile {
	return &Profile{
		Username: u.Username,
		Bio: u.Bio,
		Image: u.Image,
	}
}

var AnonymousUser User

type UserFilter struct {
	ID       *uint   `json:"id,omitempty"`
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`

	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type UserPatch struct {
	Email        *string `json:"email,omitempty"`
	Username     *string `json:"username,omitempty"`
	Image        *string `json:"image,omitempty"`
	Bio          *string `json:"bio,omitempty"`
	PasswordHash *string `json:"-" db:"password_hash"`
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
	Authenticate(ctx context.Context, email, password string) (*User, error)

	CreateUser(context.Context, *User) error

	UserByID(context.Context, uint) (*User, error)

	UserByEmail(context.Context, string) (*User, error)

	Users(context.Context, UserFilter) ([]*User, error)

	UpdateUser(context.Context, *User, UserPatch) error

	DeleteUser(context.Context, uint) error
}
