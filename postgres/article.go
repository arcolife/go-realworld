package postgres

import (
	"context"
	"errors"
	"log"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/jmoiron/sqlx"
)

type ArticleService struct {
	*DB
}

func NewArticleService(db *DB) *ArticleService {
	return &ArticleService{db}
}

func (as *ArticleService) CreateArticle(ctx context.Context, article *conduit.Article) error {
	tx, err := as.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := createArticle(ctx, tx, article); err != nil {
		return err
	}

	return tx.Commit()
}

func (as *ArticleService) ArticleByID(_ context.Context, _ uint) error {
	panic("not implemented") // TODO: Implement
}

func (as *ArticleService) Articles(_ context.Context, _ conduit.ArticleFilter) ([]*conduit.Article, error) {
	panic("not implemented") // TODO: Implement
}

func (as *ArticleService) UpdateArticle(_ context.Context, _ *conduit.Article, _ conduit.ArticlePatch) error {
	panic("not implemented") // TODO: Implement
}

func (as *ArticleService) DeleteArticle(_ context.Context, _ uint) error {
	panic("not implemented") // TODO: Implement
}

func createArticle(ctx context.Context, tx *sqlx.Tx, article *conduit.Article) error {

	query := `
	INSERT INTO articles (title, body, description, author_id, favorites_count, slug) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at
	`

	args := []interface{}{
		article.Title,
		article.Body,
		article.Description,
		article.Author.ID,
		article.FavoritesCount,
		article.Slug,
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&article.ID, &article.CreatedAt, &article.UpdatedAt)

	if err != nil {
		return err
	}

	tags := make([]string, len(article.Tags))
	for i, tag := range article.Tags {
		tags[i] = tag.Name
	}

	err = setTags(ctx, tx, article, tags)

	if err != nil {
		return err
	}

	return nil
}

func findArticles(ctx context.Context, tx *sqlx.Tx, filter conduit.ArticleFilter) ([]*conduit.Article, error) {
	// where, args := []string{}, []interface{}{}
	return nil, nil
}

func setTags(ctx context.Context, tx *sqlx.Tx, article *conduit.Article, tags []string) error {

	for _, v := range tags {
		tag, err := findTagByName(ctx, tx, v)

		if err != nil {
			switch {
			case errors.Is(err, conduit.ErrNotFound):
				tag = &conduit.Tag{Name: v}
				err = createTag(ctx, tx, tag)
				if err != nil {
					return err
				}
			default:
				log.Println(err)
				return err
			}
		}

		err = associateArticleWithTag(ctx, tx, article, tag)

		if err != nil {
			return err
		}
	}

	return nil
}

func associateArticleWithTag(ctx context.Context, tx *sqlx.Tx, article *conduit.Article, tag *conduit.Tag) error {
	query := "INSERT INTO article_tags (article_id, tag_id) VALUES ($1, $2)"
	_, err := tx.ExecContext(ctx, query, article.ID, tag.ID)

	if err != nil {
		return err
	}

	return nil
}
