package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New" // для логов и ошибок

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

	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveUrl возвращает index созданной записи
func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		// юзаем библиотеку go-sqlite3
		// смотрим, что внутри ошибки от sqlite3
		// делаем это, чтобы возвращает тот же текст ошибки клиенту, если поменяем БД
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUrlExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// поддерживается не всеми БД
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: faild to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var urlResult string
	if err := stmt.QueryRow(alias).Scan(&urlResult); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrUrlNotFound)
		}
		return "", fmt.Errorf("%s: execute statement %w", op, err)
	}

	return urlResult, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.sqlite.DeleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.Exec(alias); err != nil {
		return fmt.Errorf("%s: execute statement %w", op, err)
	}

	return nil
}
