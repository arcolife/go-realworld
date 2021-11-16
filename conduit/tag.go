package conduit

import (
	"context"
	"strconv"
)

type Tag struct {
	ID   uint
	Name string
}

func (t Tag) MarshalJSON() ([]byte, error) {
	jsonValue := strconv.Quote(t.Name)
	return []byte(jsonValue), nil
}

type TagFilter struct {
	Name *string
}

type TagService interface {
	Tags(context.Context, TagFilter) ([]*Tag, error)
}
