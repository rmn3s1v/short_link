package repository

import (
	"context"
	"database/sql"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func InitPostgres(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id BIGSERIAL PRIMARY KEY,
			original_url TEXT NOT NULL UNIQUE,
			short_url VARCHAR(10) NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func (p *PostgresRepo) Save(url, short string) error {
	_, err := p.db.Exec(`
		INSERT INTO urls (original_url, short_url)
		VALUES ($1, $2)
		ON CONFLICT (original_url) DO NOTHING
	`, url, short)

	return err
}

func (p *PostgresRepo) GetByURL(url string) (string, error) {
	var short string
	err := p.db.QueryRow(`
		SELECT short_url FROM urls WHERE original_url=$1
	`, url).Scan(&short)
	if err == sql.ErrNoRows {
		return "", nil
	}

	return short, err
}

func (p *PostgresRepo) GetByShortURL(short string) (string, error) {
	var url string
	err := p.db.QueryRow(`
		SELECT original_url FROM urls WHERE short_url=$1
	`, short).Scan(&url)
	if err == sql.ErrNoRows {
		return "", nil
	}

	return url, err
}
