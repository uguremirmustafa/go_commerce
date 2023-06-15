package middlewares

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/uguremirmustafa/go_commerce/helpers"
)

func RequireAuth(c *fiber.Ctx) error {
	cookie := c.Cookies("Authorization")
	fmt.Println(cookie)

	parsedToken, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key used for signing
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	// Check if the token is valid
	if !parsedToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	// Access the claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse claims",
		})
	}

	// Access individual claim values
	userID := claims["ID"].(float64)
	email := claims["email"].(string)
	exp := claims["exp"].(float64)

	expireDate := helpers.UnixToTime(int64(exp))

	// Set the claims in the context's locals
	c.Locals("userID", userID)
	c.Locals("email", email)
	c.Locals("exp", exp)
	c.Locals("expireDate", expireDate)

	return c.Next()
}
