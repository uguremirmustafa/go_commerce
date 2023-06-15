package routes

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/uguremirmustafa/go_commerce/database"
	"github.com/uguremirmustafa/go_commerce/models"
	validations "github.com/uguremirmustafa/go_commerce/utils"
)

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
	}

	errors := validations.ValidateUserStruct(*user)
	if errors != nil {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "could not create the user"})
	}

	return c.Status(http.StatusOK).JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}
	database.DB.Find(&users)
	return c.Status(http.StatusOK).JSON(users)
}

func findUser(id int, user *models.User) error {
	database.DB.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}
	return nil
}

func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var user models.User

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Please ensure that :id is an integer")
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Please ensure that :id is an integer")
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// if updateData.FirstName != "" {
	// 	user.FirstName = updateData.FirstName
	// }
	// if updateData.LastName != "" {
	// 	user.LastName = updateData.LastName
	// }

	database.DB.Save(&user)

	return c.Status(201).JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Please ensure that :id is an integer")
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(http.StatusOK).SendString("Successfully deleted user")
}
