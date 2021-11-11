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
	Token        string    `json:"token,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
}

type UserFilter struct {
	ID       *uint   `json:"id"`
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Token    *string `json:"token"`
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

type UserService interface {
	CreateUser(context.Context, *User) error
	//UserByID(context.Context, uint) (*User, error)
	UserByEmail(context.Context, string) (*User, error)
	//Users(context.Context, UserFilter) ([]*User, error)
	//Authenticate(email, password string) (*User, error)
}
