package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"lostemp/internal/apperror"
)

// ErrorHandler centraliza a resposta de erros em JSON consistente.
func ErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return c.Status(appErr.Status).JSON(fiber.Map{"error": appErr})
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"error": fiber.Map{"status": 404, "code": "not_found", "message": "registro não encontrado"},
		})
	}

	status := fiber.StatusInternalServerError
	message := "erro interno"
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		status = fiberErr.Code
		message = fiberErr.Message
	}
	return c.Status(status).JSON(fiber.Map{
		"error": fiber.Map{"status": status, "code": "error", "message": message},
	})
}
