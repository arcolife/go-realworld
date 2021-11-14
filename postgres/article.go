package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/jmoiron/sqlx"
)

type ArticleService struct {
	db *DB
}

func NewArticleService(db *DB) *ArticleService {
	return &ArticleService{db}
}

func (as *ArticleService) CreateArticle(ctx context.Context, article *conduit.Article) error {
	tx, err := as.db.BeginTxx(ctx, nil)

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

func (as *ArticleService) Articles(ctx context.Context, filter conduit.ArticleFilter) ([]*conduit.Article, error) {
	tx, err := as.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	articles, err := findArticles(ctx, tx, filter)

	if err != nil {
		return nil, err
	}

	for _, article := range articles {
		if err := attachArticleAssociations(ctx, tx, article); err != nil {
			return nil, err
		}
	}

	return articles, nil
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

	err = setArticleTags(ctx, tx, article, tags)

	if err != nil {
		return err
	}

	return nil
}

func findArticles(ctx context.Context, tx *sqlx.Tx, filter conduit.ArticleFilter) ([]*conduit.Article, error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0 // used to set correct postgres argument enums i.e $1, $2

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("author_id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.AuthorID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Slug; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("slug = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Title; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("title = $%d", argPosition)), append(args, *v)
	}

	whereClause := ""

	if len(where) != 0 {
		whereClause = " WHERE " + strings.Join(where, " AND ")
	}

	query := "SELECT * from articles" + whereClause + " ORDER BY id ASC"

	rows, err := tx.QueryxContext(ctx, query, args...)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	articles := make([]*conduit.Article, 0)

	for rows.Next() {
		var article conduit.Article

		if err := rows.StructScan(&article); err != nil {
			return nil, fmt.Errorf("error retrieving articles: %w", err)
		}

		articles = append(articles, &article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error parsing articles: %w", err)
	}

	return articles, nil
}

func setArticleTags(ctx context.Context, tx *sqlx.Tx, article *conduit.Article, tags []string) error {
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

func attachArticleAssociations(ctx context.Context, tx *sqlx.Tx, article *conduit.Article) error {
	tags, err := findArticleTags(ctx, tx, article)

	if err != nil {
		return fmt.Errorf("cannot find article tags: %w", err)
	}

	article.Tags = tags

	user, err := findUserByID(ctx, tx, article.AuthorID)

	if err != nil {
		return fmt.Errorf("cannot find article author: %w", err)
	}

	article.Author = user

	return nil

}

func findArticleTags(ctx context.Context, tx *sqlx.Tx, article *conduit.Article) ([]*conduit.Tag, error) {
	query := `
	SELECT * from tags WHERE id IN (
		SELECT tag_id FROM article_tags WHERE article_id = $1
	)
	`

	rows, err := tx.QueryxContext(ctx, query, article.ID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*conduit.Tag, 0)

	for rows.Next() {
		var tag conduit.Tag

		if err := rows.StructScan(&tag); err != nil {
			log.Println(err)
			return nil, err
		}

		tags = append(tags, &tag)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return tags, nil
}
