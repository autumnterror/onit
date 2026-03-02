package service

import (
	"github.com/autumnterror/onit/internal/domain"
	"github.com/autumnterror/utils_go/pkg/log"
)

type Repo interface {
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id string) error
	GetById(id string) (*domain.Product, error)
	GetAll() ([]domain.Product, error)
}

func (s *Service) Create(product *domain.Product) error {
	if err := productValidation(product); err != nil {
		return err
	}

	err := s.repo.Create(product)
	if err != nil {
		log.Red(err)
		return ErrServer
	}
	return nil
}

func (s *Service) Update(product *domain.Product) error {
	if err := productValidation(product); err != nil {
		return err
	}

	err := s.repo.Update(product)
	if err != nil {
		log.Red(err)
		return ErrServer
	}
	return nil
}
func (s *Service) Delete(id string) error {
	if err := idValidation(id); err != nil {
		return err
	}

	err := s.repo.Delete(id)
	if err != nil {
		log.Red(err)
		return ErrServer
	}
	return nil
}

func (s *Service) GetById(id string) (*domain.Product, error) {
	if err := idValidation(id); err != nil {
		return nil, err
	}

	one, err := s.repo.GetByID(id)
	if err != nil {
		log.Red(err)
		return nil, ErrServer
	}
	return one, nil
}

func (s *Service) GetAll() ([]domain.Product, error) {
	all, err := s.repo.GetAll()
	if err != nil {
		log.Red(err)
		return nil, ErrServer
	}
	return all, nil
}
