package middleware

import (
	"github.com/gofiber/fiber/v2"

	"lostemp/internal/apperror"
)

// RequireRoles garante que o usuário autenticado possua um dos papéis
// informados. Uso: middleware.RequireRoles(db.RoleAdmin, db.RoleProfessor).
func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := Role(c)
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return apperror.Forbidden("sem permissão para esta ação")
	}
}
