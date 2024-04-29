package database

import (
	"fmt"
	"reflect"
)

// ExtractStructSlice returns the field `db` columns and the values as 2d arrays
// Can take a slice of struct or a slice of pointers to structs
func ExtractStructSlice[T any](a []T) ([]string, [][]any, error) {
	var cols []string
        rows := make([][]any, len(a))

	t := reflect.TypeOf(a)

	t = t.Elem()

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("ok")
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		cols = append(cols, f.Tag.Get("db"))
	}

        row := make([]any, t.NumField())
        for i := range len(a) {
                for j := 0; j < t.NumField(); j++ {
                        row[j] = reflect.ValueOf(a[i]).FieldByName([]int{j})
                }
                rows[i] = row
        }

	return cols, rows, nil
}
