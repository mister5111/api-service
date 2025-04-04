package sqlite

import (
	"api-service/src/storage"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type AliasTableSqlite struct {
	Alias string `json:"alias"`
	Url   string `json:"url"`
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqlite3.ErrNo(sqliteErr.ExtendedCode) == sqlite3.ErrConstraint {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) ShowAlias(alias string) (AliasTableSqlite, error) {
	const op = "storage.ShowAlias"

	stmt, err := s.db.Prepare("SELECT alias, url FROM url WHERE alias = ?")
	if err != nil {
		return AliasTableSqlite{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	results := AliasTableSqlite{}
	err = stmt.QueryRow(alias).Scan(&results.Alias, &results.Url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AliasTableSqlite{}, storage.ErrALIASNotFound
		}
		return AliasTableSqlite{}, fmt.Errorf("%s: %w", op, err)
	}

	return results, nil
}

func (s *Storage) ShowAll() ([]AliasTableSqlite, error) {
	const op = "storage.ShowAll"

	stmt, err := s.db.Prepare("SELECT alias, url FROM url")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrALIASNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var results []AliasTableSqlite

	for rows.Next() {
		var row AliasTableSqlite
		if err := rows.Scan(&row.Alias, &row.Url); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return results, nil
}

func (s *Storage) Delete(alias string) error {
	const op = "storage.sqlite.Delete"

	stmt, err := s.db.Prepare("SELECT alias FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var rowsAlias string
	err = stmt.QueryRow(alias).Scan(&rowsAlias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrALIASNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	delete, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = delete.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
