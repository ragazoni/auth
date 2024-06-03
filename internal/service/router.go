package service

import "github.com/gofiber/fiber/v3"

func SetRouters(r fiber.Router) {

	users := r.Group("/users")
	users.Post("/", addUser)
	users.Post("/login", loginUser)
	users.Get("/", getAll)
	users.Get("/:id", getById)
	users.Put("/:id", updateUser)
	users.Delete("/:id", deleteUser)
	users.Get("/email/:email", findByEmail)
	users.Get("/check/:email", checkAdmin)
	users.Post("/logout", logout)
}
