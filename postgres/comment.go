package postgres

import (
	"context"
	"fmt"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/jmoiron/sqlx"
)

type CommentService struct {
	*DB
}

func NewCommentService(db *DB) *CommentService {
	return &CommentService{db}
}

var _ conduit.CommentService = (*CommentService)(nil)

func (cs *CommentService) CreateComment(ctx context.Context, comment *conduit.Comment) error {
	tx, err := cs.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := createComment(ctx, tx, comment); err != nil {
		return err
	}

	return tx.Commit()
}

func (cs *CommentService) Comment(ctx context.Context, id uint) error {
	panic("not implemented") // TODO: Implement
}

func (cs *CommentService) Comments(ctx context.Context, cf conduit.CommentFilter) ([]*conduit.Comment, error) {
	tx, err := cs.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	comments, err := findComments(ctx, tx, cf)

	if err != nil {
		return comments, err
	}

	return comments, tx.Commit()
}

func (cs *CommentService) DeleteComment(ctx context.Context, id uint) error {
	panic("not implemented") // TODO: Implement
}

func createComment(ctx context.Context, tx *sqlx.Tx, comment *conduit.Comment) error {
	query := `
	INSERT INTO comments (body, article_id, author_id) VALUES ($1, $2, $3) 
	RETURNING id, created_at
	`

	args := []interface{}{
		comment.Body,
		comment.ArticleID,
		comment.AuthorID,
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func findComments(ctx context.Context, tx *sqlx.Tx, filter conduit.CommentFilter) ([]*conduit.Comment, error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0 // used to set correct postgres argument enums i.e $1, $2

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.AuthorID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("author_id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.ArticleID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("article_id = $%d", argPosition)), append(args, *v)
	}

	query := "SELECT * from comments" + formatWhereClause(where) + " ORDER BY created_at DESC " +
		formatLimitOffset(filter.Limit, filter.Offset)

	comments := make([]*conduit.Comment, 0)

	if err := findMany(ctx, tx, &comments, query, args...); err != nil {
		return comments, err
	}

	for _, c := range comments {
		if err := attachCommentAuthor(ctx, tx, c); err != nil {
			return comments, err
		}
	}

	return comments, nil
}

func attachCommentAuthor(ctx context.Context, tx *sqlx.Tx, comment *conduit.Comment) error {
	user, err := findUserByID(ctx, tx, comment.AuthorID)

	if err != nil {
		return err
	}

	comment.Author = user
	return nil
}
