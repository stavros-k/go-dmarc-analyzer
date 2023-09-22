package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/stavros-k/go-dmarc-analyzer/internal/database"
	"github.com/stavros-k/go-dmarc-analyzer/internal/routes"
)

type APIServer struct {
	host  string
	port  int
	store database.Storage
}

func NewAPIServer(host string, port int, store database.Storage) *APIServer {
	return &APIServer{
		host:  host,
		port:  port,
		store: store,
	}
}

func (s *APIServer) RegisterRoutesAndStart() error {
	app := fiber.New()

	// Register routes
	app.Get("/health", routes.HandleHealth)

	return app.Listen(fmt.Sprintf("%s:%d", s.host, s.port))
}
