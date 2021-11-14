package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/jmoiron/sqlx"
)

func createTag(ctx context.Context, tx *sqlx.Tx, tag *conduit.Tag) error {
	query := "INSERT INTO tags (name) VALUES ($1) RETURNING id"

	err := tx.QueryRowxContext(ctx, query, tag.Name).Scan(&tag.ID)

	if err != nil {
		return err
	}

	return nil
}

func findTags(ctx context.Context, tx *sqlx.Tx, filter conduit.TagFilter) ([]*conduit.Tag, error) {
	// Build WHERE clause.
	where, args := []string{}, []interface{}{}
	argPosition := 0

	if v := filter.Name; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("name = $%d", argPosition)), append(args, *v)
	}

	query := "SELECT * from tags WHERE " + strings.Join(where, " AND ") + " ORDER BY id ASC"
	rows, err := tx.QueryxContext(ctx, query, args...)

	if err != nil {
		log.Println(err)
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

func findTagByName(ctx context.Context, tx *sqlx.Tx, name string) (*conduit.Tag, error) {
	ts, err := findTags(ctx, tx, conduit.TagFilter{Name: &name})

	if err != nil {
		return nil, err
	} else if len(ts) == 0 {
		return nil, conduit.ErrNotFound
	}

	return ts[0], nil
}
