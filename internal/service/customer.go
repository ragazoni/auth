package service

import (
	"auth/internal/db"
	"auth/internal/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func addUser(c fiber.Ctx) error {
	//body := new(models.Customer)
	var user models.Customer

	if err := c.Bind().Body(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(password)

	_, err := db.Insert("users", user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot register user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func getAll(c fiber.Ctx) error {
	var document []models.Customer

	err := db.Find("users", &document)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	return c.JSON(document)
}

func getById(c fiber.Ctx) error {
	var document models.Customer

	err := db.FindById("users", c.Params("id"), &document)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	return c.JSON(document)
}

func updateUser(c fiber.Ctx) error {
	body := new(models.Customer)

	if err := c.Bind().Body(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON("invalid json")
	}

	var result models.Customer

	err := db.UpdateById("users", c.Params("id"), body, &result)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	return c.Status(http.StatusCreated).JSON(result)
}

func deleteUser(c fiber.Ctx) error {
	errr := db.DeletebyId("users", c.Params("id"))
	if errr != nil {
		return c.Status(http.StatusInternalServerError).JSON(errr.Error())
	}
	return c.Status(http.StatusNoContent).SendString("")
}

func findByEmail(c fiber.Ctx) error {
	email := c.Params("email")
	result, err := db.FindByEmail("users", email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error find usr by email",
			"error":   err.Error(),
		})

	}

	if result == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "user not found",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user found",
		"user":    result,
	})

}

var jwtSecret = []byte("secret")

func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func loginUser(c fiber.Ctx) error {
	data := new(models.Customer)
	//var data map[string]string

	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	//var user models.Customer
	user, err := db.FindByEmail("users", data.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if user == nil || !user.CheckPassword(data.Password) {
		return c.Status(http.StatusUnauthorized).JSON("invalid email or password")
	}

	token, err := GenerateJWT(user.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON("could not generate token")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
}

func logout(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token is required",
		})
	}
	models.InvalidTokens.Lock()
	models.InvalidTokens.Tokens[token] = true
	models.InvalidTokens.Unlock()

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})

}
