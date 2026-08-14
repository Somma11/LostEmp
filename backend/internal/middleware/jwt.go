package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"lostemp/internal/apperror"
)

const (
	ctxUserID = "user_id"
	ctxRole   = "role"
	ctxEmail  = "email"
)

// JWT autentica requisições com token Bearer (Authorization: Bearer <token>).
func JWT(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return apperror.Unauthorized("token de autenticação ausente")
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inválido")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return apperror.Unauthorized("token inválido ou expirado")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return apperror.Unauthorized("token inválido")
		}

		sub, err := claims.GetSubject()
		if err != nil {
			return apperror.Unauthorized("token inválido")
		}
		id, err := strconv.ParseUint(sub, 10, 64)
		if err != nil {
			return apperror.Unauthorized("token inválido")
		}

		role, _ := claims["role"].(string)
		email, _ := claims["email"].(string)

		c.Locals(ctxUserID, uint(id))
		c.Locals(ctxRole, role)
		c.Locals(ctxEmail, email)
		return c.Next()
	}
}

// UserID lê o usuário autenticado a partir dos locals do Fiber.
func UserID(c *fiber.Ctx) uint {
	id, _ := c.Locals(ctxUserID).(uint)
	return id
}

// Role lê o papel do usuário autenticado.
func Role(c *fiber.Ctx) string {
	role, _ := c.Locals(ctxRole).(string)
	return role
}
