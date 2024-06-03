package service

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func checkAdmin(c fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "missing email"})
	}

	isAdmin, err := IsAdmin(email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"isadmin": isAdmin})
}
