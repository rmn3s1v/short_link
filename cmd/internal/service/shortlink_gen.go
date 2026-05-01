package service

import (
	"errors"
	"net/url"
	"short-link/cmd/internal/repository"
	"short-link/cmd/internal/utils"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Shorten(original string) (string, error) {
	if _, err := url.ParseRequestURI(original); err != nil {
		return "", errors.New("invalid url")
	}

	existing, err := s.repo.GetByURL(original)
	if err != nil {
		return "", err
	}
	if existing != "" {
		return existing, nil
	}

	short := utils.Generate(original)
	if err := s.repo.Save(original, short); err != nil {
		return "", err
	}

	return short, nil
}

func (s *Service) Resolve(short string) (string, error) {
	return s.repo.GetByShortURL(short)
}
