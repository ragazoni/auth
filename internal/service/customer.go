package service

import (
	"auth/internal/db"
	"auth/internal/models"
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

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
	//	body := new(models.Customer)
	var data map[string]string

	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	var user models.Customer
	err := userCollection.FindOne(context.Background(), bson.M{"cpf": data["cpf"]}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.Hex(),
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not login",
		})
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}
