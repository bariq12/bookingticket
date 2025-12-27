package middleware

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bariq12/bookingticket/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Authprotected(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			log.Warn("Authorization header is missing")
			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Authorization header is required",
			})
		}
		tokenpart := strings.Split(authHeader, " ")
		if len(tokenpart) != 2 || tokenpart[0] != "Bearer" {
			log.Warn("Invalid Authorization header format")
			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Invalid Authorization header format",
			})
		}
		tokenstr := tokenpart[1]
		secret := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		if err != nil {
			log.Warn("Error parsing token: ", err)
			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Invalid token",
			})
		}
		Userid := token.Claims.(jwt.MapClaims)["id"]
		if err := db.First(&models.User{}).Where("id = ?", Userid).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found: ", err)
			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "User not found",
			})
		}
		ctx.Locals("userId", Userid)
		return ctx.Next()
	}
}
