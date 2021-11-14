package conduit

import (
	"context"
	"encoding/json"
	"time"
)

type Article struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	Description    string    `json:"description"`
	AuthorID       uint      `json:"-"`
	Author         *User     `json:"author"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount uint      `json:"favoritesCount"`
	Slug           string    `json:"slug"`
	Tags           []*Tag    `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
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
	ID          *uint
	Title       *string
	Description *string
	AuthorID    *uint
	Slug        *string
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
