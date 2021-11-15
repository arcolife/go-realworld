package postgres

import (
	"fmt"
	"strings"
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

// func findMany(ctx context.Context, tx *sqlx.Tx, ss interface{}, query string, args ...interface{}) error {
// 	rows, err := tx.QueryxContext(ctx, query, args...)

// 	if err != nil {
// 		return err
// 	}

// 	defer rows.Close()

// 	if reflect.TypeOf(ss).Kind() != reflect.Ptr {
// 		return fmt.Errorf("expecting a pointer to a slice")
// 	}

// 	if reflect.TypeOf(ss).Elem().Kind() != reflect.Slice {
// 		return fmt.Errorf("expecting  a pointer to a slice")

// 	}

// 	sPtrVal := reflect.ValueOf(ss) // pointer to slice value
// 	ssVal := sPtrVal.Elem()        // slice pointed to by ss
// 	typ := ssVal.Type().Elem()     // type of slice elements

// 	newSlice := reflect.MakeSlice(ssVal.Type(), 10, 10) // new slice

// 	for rows.Next() {
// 		newVal := reflect.New(typ)

// 		if err := rows.StructScan(newVal.Interface()); err != nil {
// 			return nil
// 		}
// 		newSlice = reflect.Append(newSlice, newVal)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return err
// 	}

// 	sPtrVal.Elem().Set(newSlice)

// 	return nil
// }
