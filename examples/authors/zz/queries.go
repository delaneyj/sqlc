package zz

import (
	"fmt"
	"zombiezen.com/go/sqlite"
)

type GetAuthorRes struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func GetAuthor(
	tx *sqlite.Conn,
	id int64,
) (
	res *GetAuthorRes,
	err error,
) {
	// Prepare statement into cache
	stmt := tx.Prep(`SELECT id, name, bio FROM authors
WHERE id = ? LIMIT 1`)
	defer stmt.Reset()

	// Bind parameters
	stmt.BindInt64(1, id)

	// Execute query
	if hasRow, err := stmt.Step(); err != nil {
		return res, err
	} else if hasRow {

		row := GetAuthorRes{
			Id:   stmt.ColumnInt64(1),
			Name: stmt.ColumnText(2),
			Bio:  stmt.ColumnText(3),
		}
		res = &row
	}

	return res, nil
}

type ListAuthorsRes struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func ListAuthors(
	tx *sqlite.Conn,
) (
	res []ListAuthorsRes,
	err error,
) {
	// Prepare statement into cache
	stmt := tx.Prep(`SELECT id, name, bio FROM authors
ORDER BY name`)
	defer stmt.Reset()

	// Execute query
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return res, fmt.Errorf("failed to execute listauthors SQL: %w", err)
		} else if !hasRow {
			break
		}

		row := ListAuthorsRes{
			Id:   stmt.ColumnInt64(1),
			Name: stmt.ColumnText(2),
			Bio:  stmt.ColumnText(3),
		}

		res = append(res, row)
	}

	return res, nil
}

func CreateAuthor(
	tx *sqlite.Conn,
	name string,
	bio string,
) (
	err error,
) {
	// Prepare statement into cache
	stmt := tx.Prep(``)
	defer stmt.Reset()

	// Bind parameters
	stmt.BindText(1, name)
	stmt.BindText(2, bio)

	// Execute query
	if _, err := stmt.Step(); err != nil {
		return fmt.Errorf("failed to execute createauthor SQL: %w", err)
	}

	return nil
}

func DeleteAuthor(
	tx *sqlite.Conn,
	id int64,
) (
	err error,
) {
	// Prepare statement into cache
	stmt := tx.Prep(``)
	defer stmt.Reset()

	// Bind parameters
	stmt.BindInt64(1, id)

	// Execute query
	if _, err := stmt.Step(); err != nil {
		return fmt.Errorf("failed to execute deleteauthor SQL: %w", err)
	}

	return nil
}
