package main

import (
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const length = 6                       // Length of the short url
var urlStore = make(map[string]string) // In-memory storage for URL mappings

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("views/index.html", fiber.Map{})
	})

	app.Post("/", func(c *fiber.Ctx) error {
		url := c.FormValue("url")

		if url == "" {
			return c.Render("views/index.html", fiber.Map{"url": "URL cannot be empty"})
		}

		surl := GenerateShortURL()
		urlStore[surl] = url

		return c.Render("views/index.html", fiber.Map{"url": surl})
	})

	app.Get("/:shorturl", func(c *fiber.Ctx) error {
		surl := c.Params("shorturl")

		url, exists := urlStore[surl]
		if !exists {
			return c.Status(404).JSON(fiber.Map{"error": "URL not found"})
		}

		return c.Redirect(url, 301)
	})

	app.Listen(":5000")
}

func GenerateShortURL() string {
	result := make([]byte, length) // Create an empty slice with zeros

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))] // Randomly pick a character from the letters string
	}

	return string(result)
}
