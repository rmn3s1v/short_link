package repository

type Repository interface {
	Save(url, short_url string) error
	GetByURL(url string) (string, error)
	GetByShortURL(short_url string) (string, error)
}
