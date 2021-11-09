package conduit

import "time"

type Article struct {
	ID             uint      `json:"id,omitempty"`
	Title          string    `json:"title,omitempty"`
	Body           string    `json:"body,omitempty"`
	Description    string    `json:"description,omitempty"`
	Author         User      `json:"author,omitempty"`
	Favorited      []User    `json:"favorited,omitempty"`
	FavoritesCount uint      `json:"favoritesCount,omitempty"`
	Slug           string    `json:"slug,omitempty"`
	Tags           []Tag     `json:"tagList,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}

type Tag struct {
	Name string
}

type Favorite struct {
	UserID    uint
	ArticleID uint
}
