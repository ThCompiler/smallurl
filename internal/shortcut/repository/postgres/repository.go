package postgres

import (
	"context"
	"smallurl/internal/shortcut/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

const (
	getLongURLQuery = "SELECT long_url from urls where short_url = $1"

	addURLQuery = "INSERT INTO urls (short_url, long_url) VALUES ($1, $2)"
)

type Repository struct {
	db AdapterInterface
}

func NewRepository(db AdapterInterface) *Repository {
	return &Repository{
		db: db,
	}
}

func (rp *Repository) GetLongURL(shortURL string) (string, error) {
	var longURL string
	if err := rp.db.QueryRow(context.Background(), getLongURLQuery, shortURL).Scan(&longURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.Wrapf(repository.ErrorURLNotFound, "short url %s not found", shortURL)
		}

		return "", errors.Wrapf(err, "try get long url for short url %s", shortURL)
	}

	return longURL, nil
}

func (rp *Repository) AddURL(shortURL, longURL string) error {
	if _, err := rp.db.Exec(context.Background(), addURLQuery, shortURL, longURL); err != nil {
		return errors.Wrapf(checkPgConflictError(err), "try add long url %s with short url %s", longURL, shortURL)
	}

	return nil
}

const (
	uniqueConflictCode   = "23505"
	uniqueConstraintName = "urls_short_url_key"
)

func checkPgConflictError(err error) error {
	var e *pgconn.PgError

	if !errors.As(err, &e) {
		return err
	}

	if e.Code == uniqueConflictCode && e.ConstraintName == uniqueConstraintName {
		return repository.ErrorURLAlreadyExists
	}

	return err
}
