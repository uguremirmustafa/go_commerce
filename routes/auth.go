package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/uguremirmustafa/go_commerce/database"
	"github.com/uguremirmustafa/go_commerce/helpers"
	"github.com/uguremirmustafa/go_commerce/models"
	validations "github.com/uguremirmustafa/go_commerce/utils"
	"golang.org/x/crypto/bcrypt"
)

func Singup(c *fiber.Ctx) error {
	// Get the email/password off req body
	body := new(models.SignupRequest)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// validate request model
	errors := validations.ValidateSignupRequestStruct(*body)
	if errors != nil {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := database.DB.Create(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": result.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user created",
		"user":    user,
	})
}

func Login(c *fiber.Ctx) error {
	// extract the credentials from the request body
	body := new(models.LoginRequest)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := validations.ValidateLoginRequestStruct(*body)
	if errors != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "login failed"})
	}

	user := new(models.User)

	res := database.DB.Where("email = ?", body.Email).First(&user)

	if res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "wrong password",
		})
	}

	// create jwt claims
	claims := jwt.MapClaims{
		"ID":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// generate encoded token and send it as response
	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    t,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: "Lax",
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"success": "Authorization cookie is set",
	})
}

func ValidateAuthorization(c *fiber.Ctx) error {
	// Get the token from the request headers or cookies
	token := c.Cookies("Authorization")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token not found",
		})
	}

	// Parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
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

	expiretDate := helpers.UnixToTime(int64(exp))

	// Return the claims as a response
	return c.JSON(fiber.Map{
		"userID":      userID,
		"email":       email,
		"expiretDate": expiretDate,
	})
}
