package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("GET Hello")
	})

	app.Post("/hello", func(c *fiber.Ctx) error {
		return c.SendString("POST Hello")
	})

	//Parameters and Optional Parameters
	// app.Post("/hello/:name/:surname?", func(c *fiber.Ctx) error {
	// 	name := c.Params("name")
	// 	surname := c.Params("surname")
	// 	return c.SendString("POST Hello, " + name + surname)
	// })

	//ParamsInt
	app.Post("/hello/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.ErrBadRequest
		}
		return c.SendString(fmt.Sprintf("ID: %v", id))
	})
	
	app.Listen(":8000")
}