package conduit

import (
	"context"
	"time"
)

type Article struct {
	ID             uint      `json:"-"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	Description    string    `json:"description"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount uint      `json:"favoritesCount" db:"favorites_count"`
	Slug           string    `json:"slug"`
	AuthorID       uint      `json:"-" db:"author_id"`
	Author         *User     `json:"-"`
	AuthorProfile  *Profile  `json:"author"`
	Tags           []*Tag    `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

func (a *Article) SetAuthorProfile(currentUser *User) {
	a.AuthorProfile = &Profile{
		Username: a.Author.Username,
		Bio:      a.Author.Bio,
		Image:    a.Author.Image,
	}

	a.AuthorProfile.Following = currentUser.IsFollowing(a.Author)
}

type ArticleFilter struct {
	ID             *uint
	Title          *string
	Description    *string
	AuthorID       *uint
	AuthorUsername *string
	Tag            *string
	Slug           *string

	Limit  int
	Offset int
}

type ArticlePatch struct {
	Title       *string
	Body        *string
	Description *string
	Tags        []Tag
}

type Favorite struct {
	UserID    uint
	ArticleID uint
}

type ArticleService interface {
	CreateArticle(context.Context, *Article) error
	ArticleBySlug(context.Context, string) (*Article, error)
	Articles(context.Context, ArticleFilter) ([]*Article, error)
	ArticleFeed(context.Context, *User, ArticleFilter) ([]*Article, error)
	UpdateArticle(context.Context, *Article, ArticlePatch) error
	DeleteArticle(context.Context, uint) error
}

func (a *Article) AddTags(_tags ...string) {
	for _, t := range _tags {
		a.Tags = append(a.Tags, &Tag{Name: t})
	}
}
