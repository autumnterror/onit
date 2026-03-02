package service

import "github.com/autumnterror/onit/internal/repo"

type Service struct {
	repo repo.ProductRepository
}

func NewService(
	repo repo.ProductRepository,
) *Service {
	return &Service{
		repo: repo,
	}
}
