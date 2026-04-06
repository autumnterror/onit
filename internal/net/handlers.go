package net

import (
	"errors"
	"github.com/autumnterror/onit/internal/domain"
	"github.com/autumnterror/onit/internal/service"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (e *Echo) GetAll(c echo.Context) error {
	all, err := e.repo.GetAll()
	if err != nil {
		if errors.Is(err, service.ErrServer) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, domain.ToHttpMany(all))
}

func (e *Echo) GetById(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "нет id")
	}

	all, err := e.repo.GetById(id)
	if err != nil {
		if errors.Is(err, service.ErrServer) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, all.ToHttp())
}

func (e *Echo) Delete(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "нет id")
	}

	err := e.repo.Delete(id)
	if err != nil {
		if errors.Is(err, service.ErrServer) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (e *Echo) Create(c echo.Context) error {
	var req domain.ProductHttp
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "плохой продукт")
	}

	id := uid.New()
	req.Id = id
	err := e.repo.Create(req.ToDomain())
	if err != nil {
		if errors.Is(err, service.ErrServer) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	all, err := e.repo.GetById(id)
	if err != nil {
		if errors.Is(err, service.ErrServer) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, all)
}

func (e *Echo) Update(c echo.Context) error {
	var req domain.ProductHttp
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "плохой продукт")
	}

	err := e.repo.Update(req.ToDomain())
	if err != nil {
		if errors.Is(err, service.ErrServer) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	all, err := e.repo.GetById(req.Id)
	if err != nil {
		if errors.Is(err, service.ErrServer) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, all)
}

func (e *Echo) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, "health123")
}
