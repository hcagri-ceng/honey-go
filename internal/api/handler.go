package api

import (
	"github.com/hcagri-ceng/honey-go/internal/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hcagri-ceng/honey-go/internal/models"
)

type Server struct {
	app  *fiber.App
	repo *storage.SQLiteRepo
}

func NewServer(repo *storage.SQLiteRepo) *Server {
	app := fiber.New()

	// Frontend'in API'ye erişebilmesi için CORS ayarlarını açıyoruz
	app.Use(cors.New())

	s := &Server{
		app:  app,
		repo: repo,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Statik klasörü (HTML/CSS/JS) root dizininden sun
	s.app.Static("/", "./public")

	// API endpoint'imiz
	s.app.Get("/api/events", s.handleGetEvents)
}

func (s *Server) handleGetEvents(c *fiber.Ctx) error {
	events, err := s.repo.GetEvents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Veriler getirilemedi",
		})
	}

	// Boş slice durumunda frontend'in patlamaması için boş dizi dönüyoruz
	if events == nil {
		events = []models.Event{} // models importunu unutma
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(events),
		"data":    events,
	})
}

// Start API sunucusunu başlatır
func (s *Server) Start(port string) error {
	return s.app.Listen(port)
}

// Shutdown API sunucusunu güvenlice kapatır
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
