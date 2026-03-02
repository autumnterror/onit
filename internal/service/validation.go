package service

import (
	"errors"
	"github.com/autumnterror/onit/internal/domain"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
	"time"
)

func idValidation(id string) error {
	if uid.Validate(id) {
		return nil
	}
	return errors.New("неправильный id")
}

func stringEmpty(s string) bool {
	if s == "" {
		return true
	}
	return false
}

func productValidation(u *domain.Product) error {
	if err := idValidation(u.ID); err != nil {
		return err
	}
	if stringEmpty(u.Title) {
		u.Title = "Без названия"
	}
	if stringEmpty(u.Description) {
		u.Description = "Без описания"
	}

	if stringEmpty(u.Photo) {
		u.Photo = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR601q-F45T4q-lBZDLKPTdrTeCIHC5KFyLrg&s"
	}

	if u.Price == 0 {
		u.Price = 1000000
	}

	if u.CreatedAt == 0 {
		u.CreatedAt = time.Now().Unix()
	}
	if u.UpdatedAt == 0 {
		u.UpdatedAt = time.Now().Unix()
	}

	return nil
}
