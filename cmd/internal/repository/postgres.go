package repository

import "database/sql"

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo{
	return &PostgresRepo{
		db: db,
	}
}

func (p *PostgresRepo) Save(url, short string) error{
	_, err := p.db.Exec(`
		INSERT INTO urls (original_url, short_url)
		VALUES ($1, $2)
		ON CONFLICT (original_url) DO NOTHING
	`, url, short)

	return err
}

func (p *PostgresRepo) GetByURL(url string) (string, error){
	var short string
	err := p.db.QueryRow(`
		SELECT short_url FROM urls WHERE original_url=$1
	`, url).Scan(&short)

	return short, err
}

func (p *PostgresRepo) GetByShortURL(short string) (string, error){
	var url string
	err := p.db.QueryRow(`
		SELECT original_url FROM urls WHERE short_url=$1
	`, short).Scan(&url)

	return url, err
}
