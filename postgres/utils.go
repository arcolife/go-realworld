package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

func formatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	} else if offset > 0 {
		return fmt.Sprintf("OFFSET %d", offset)
	}
	return ""
}

func formatWhereClause(where []string) string {
	if len(where) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(where, " AND ")
}

func findMany(ctx context.Context, tx *sqlx.Tx, ss interface{}, query string, args ...interface{}) error {
	rows, err := tx.QueryxContext(ctx, query, args...)

	if err != nil {
		return err
	}

	defer rows.Close()

	sPtrVal, err := asSlicePtr(ss) // get the reflect.Value of the pointer to slice
	if err != nil {
		return err
	}
	sVal := sPtrVal.Elem()                           // get the relfect.Value of the slice pointed to by ss
	newSlice := reflect.MakeSlice(sVal.Type(), 0, 0) // new slice

	for rows.Next() {
		typ := sVal.Type().Elem().Elem() // to get conduit.Article from []*conduit.Article
		newVal := reflect.New(typ)
		if err := rows.StructScan(newVal.Interface()); err != nil {
			return nil
		}
		newSlice = reflect.Append(newSlice, newVal)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	sPtrVal.Elem().Set(newSlice)

	return nil
}

func asSlicePtr(v interface{}) (reflect.Value, error) {
	typ := reflect.TypeOf(v)

	if typ.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("expecting a pointer to a slice")
	}

	if typ.Elem().Kind() != reflect.Slice {
		return reflect.Value{}, errors.New("expecting  a pointer to a slice")
	}

	return reflect.ValueOf(v), nil
}
