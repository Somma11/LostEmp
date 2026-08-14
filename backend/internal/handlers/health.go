package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(d *gorm.DB) *HealthHandler {
	return &HealthHandler{db: d}
}

// Health verifica se a API e o banco estão respondendo.
// @Summary Health check
// @Description Retorna status da API e conectividade com o banco.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 503 {object} map[string]any
// @Router /health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status": "unavailable",
			"db":     "error",
			"time":   time.Now(),
		})
	}
	if err := sqlDB.Ping(); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status": "unavailable",
			"db":     "error",
			"time":   time.Now(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "ok",
		"db":     "ok",
		"time":   time.Now(),
	})
}
