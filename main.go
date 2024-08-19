package main

import (
	"log"

	"github.com/RajaSunrise/fiber-chatbot/models"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	bbai := models.NewBLACKBOXAI()

	app.Post("/ask", func(c *fiber.Ctx) error {
		var promt models.Prompt
		if err := c.BodyParser(&promt); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Terjadi error nih, coba cek data kamu lagi")
		}

		response, err := bbai.Ask(promt.Prompt, false, false, "", false)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"response": response,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
