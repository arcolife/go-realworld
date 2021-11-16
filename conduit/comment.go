package conduit

import (
	"context"
	"time"
)

type Comment struct {
	ID            uint      `json:"-"`
	ArticleID     uint      `json:"-" db:"article_id"`
	Article       *Article  `json:"-"`
	AuthorID      uint      `json:"-" db:"author_id"`
	Author        *User     `json:"-"`
	AuthorProfile *Profile  `json:"author"`
	Body          string    `json:"body"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type CommentFilter struct {
	ID        *uint
	ArticleID *uint

	Limit  int
	Offset int
}

type CommentService interface {
	CreateComment(context.Context, *Comment) error
	Comment(context.Context, uint) error
	Comments(context.Context, CommentFilter) ([]*Comment, error)
	DeleteComment(context.Context, uint) error
}
