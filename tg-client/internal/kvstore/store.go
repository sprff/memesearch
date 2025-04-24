package kvstore

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Store[T any] struct {
	db *sqlx.DB
}

func New[T any](path string) (Store[T], error) {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return Store[T]{}, fmt.Errorf("can't open: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return Store[T]{}, fmt.Errorf("can't ping: %w", err)
	}

	query := `
    CREATE TABLE IF NOT EXISTS storage (
        id TEXT PRIMARY KEY,
        value TEXT
	)`

	_, err = db.Exec(query)
	if err != nil {
		return Store[T]{}, fmt.Errorf("can't prepare kv: %w", err)
	}

	return Store[T]{
		db: db,
	}, nil

}

func (s *Store[T]) Get(key string) (res T, ok bool) {
	var b string
	err := s.db.Get(&b, "SELECT value FROM storage WHERE id=$1", key)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, false
		}
		panic(fmt.Sprintf("can't select: %v", err))
	}

	err = json.Unmarshal([]byte(b), &res)
	if err != nil {
		panic(fmt.Sprintf("can't decode: %v", err))
	}
	return res, true
}

func (s *Store[T]) Set(key string, value T) error {
	b, err := json.Marshal(value)
	if err != nil {
		panic(fmt.Sprintf("can't encode: %v", err))
	}

	_, ok := s.Get(key)
	if ok {
		_, err = s.db.Exec("UPDATE storage SET value=? WHERE id=?", string(b), key)
	} else {
		_, err = s.db.Exec("INSERT INTO storage (id, value) VALUES (?, ?)", key, string(b))
	}
	if err != nil {
		return fmt.Errorf("can't exec: %w", err)
	}
	return nil
}
