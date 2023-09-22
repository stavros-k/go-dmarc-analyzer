package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stavros-k/go-dmarc-analyzer/internal/types"
)

func HandleHealth(c *fiber.Ctx) error {
	return c.JSON(&types.HealthResponse{
		Status: "OK",
	})

}
