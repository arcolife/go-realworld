package postgres

import (
	"context"

	"github.com/0xdod/go-realworld/conduit"
)

type CommentService struct {
	*DB
}

var _ conduit.CommentService = (*CommentService)(nil)

func (cs *CommentService) CreateComment(ctx context.Context, c *conduit.Comment) error {
	panic("not implemented") // TODO: Implement
}

func (cs *CommentService) Comment(ctx context.Context, id uint) error {
	panic("not implemented") // TODO: Implement
}

func (cs *CommentService) Comments(ctx context.Context, cf conduit.CommentFilter) ([]*conduit.Comment, error) {
	panic("not implemented") // TODO: Implement
}

func (cs *CommentService) DeleteComment(ctx context.Context, id uint) error {
	panic("not implemented") // TODO: Implement
}

func createComment() error { return nil }
