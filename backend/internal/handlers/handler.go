package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"lostemp/internal/apperror"
	"lostemp/internal/db"
	"lostemp/internal/middleware"
)

type meta struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func pagination(c *fiber.Ctx) (page, pageSize int) {
	page = c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	pageSize = c.QueryInt("page_size", 20)
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func parseID(c *fiber.Ctx, name string) (uint, error) {
	id, err := strconv.ParseUint(c.Params(name), 10, 64)
	if err != nil || id == 0 {
		return 0, apperror.BadRequest("id inválido")
	}
	return uint(id), nil
}

func currentUser(c *fiber.Ctx) *db.User {
	return &db.User{ID: middleware.UserID(c), Role: middleware.Role(c)}
}

func listResponse(data any, page, pageSize int, total int64) fiber.Map {
	return fiber.Map{"data": data, "meta": meta{Page: page, PageSize: pageSize, Total: total}}
}
