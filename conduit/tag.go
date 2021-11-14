package conduit

import (
	"context"
	"strconv"
)

type Tag struct {
	ID   uint   `json:"-"`
	Name string `json:"name"`
}

func (t Tag) MarshalJSON() ([]byte, error) {
	jsonValue := strconv.Quote(t.Name)
	return []byte(jsonValue), nil
}

type TagFilter struct {
	Name *string
}

type TagService interface {
	AddTag(context.Context, *Tag, *Article) error
}
