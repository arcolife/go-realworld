package conduit

import (
	"context"
	"encoding/json"
	"time"
)

type Article struct {
	ID             uint      `json:"-"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	Description    string    `json:"description"`
	AuthorID       uint      `json:"-" db:"author_id"`
	Author         *User     `json:"author"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount uint      `json:"favoritesCount" db:"favorites_count"`
	Slug           string    `json:"slug"`
	Tags           []*Tag    `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

func (a Article) MarshalJSON() ([]byte, error) {
	type Author struct {
		Username  string `json:"username"`
		Bio       string `json:"bio"`
		Image     string `json:"image"`
		Following bool   `json:"following"`
	}
	var author Author

	if a.Author != nil {
		author = Author{
			Username: a.Author.Username,
			Bio:      a.Author.Bio,
			Image:    a.Author.Image,
		}
	}

	type ArticleAlias Article

	aux := struct {
		ArticleAlias
		Author `json:"author"`
	}{
		ArticleAlias: ArticleAlias(a),
		Author:       author,
	}

	return json.Marshal(aux)
}

type ArticleFilter struct {
	ID             *uint
	Title          *string
	Description    *string
	AuthorID       *uint
	AuthorUsername *string
	Tag            *string
	Slug           *string
}

type ArticlePatch struct {
	Title       *string
	Description *string
	Tags        []Tag
}

type Favorite struct {
	UserID uint
	PostID uint
}

type ArticleService interface {
	CreateArticle(context.Context, *Article) error
	ArticleByID(context.Context, uint) error
	Articles(context.Context, ArticleFilter) ([]*Article, error)
	UpdateArticle(context.Context, *Article, ArticlePatch) error
	DeleteArticle(context.Context, uint) error
}

func (a *Article) AddTags(_tags ...string) {
	for _, t := range _tags {
		a.Tags = append(a.Tags, &Tag{Name: t})
	}
}
