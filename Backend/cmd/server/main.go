// @title Whatz Chat API
// @version 1.0.0
// @description API para sistema de chat em tempo real Whatz
// @termsOfService http://swagger.io/terms/

// @contact.name Rafael Bit
// @contact.url https://github.com/rafael-bit
// @contact.email rafael@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Digite "Bearer" seguido de um espaço e o token JWT.

package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	ws "github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	_ "github.com/rafael-bit/whatz/docs"
	"github.com/rafael-bit/whatz/internal/controllers"
	"github.com/rafael-bit/whatz/internal/database"
	"github.com/rafael-bit/whatz/internal/logger"
	"github.com/rafael-bit/whatz/internal/repository"
	"github.com/rafael-bit/whatz/internal/services"
	"github.com/rafael-bit/whatz/internal/websocket"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ Arquivo .env não encontrado, usando variáveis de ambiente padrão")
	}

	// Inicializar banco de dados
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./whatz.db"
	}

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar com banco de dados: %v", err)
	}
	defer db.Close()

	// Inicializar repositórios
	userRepo := repository.NewUserRepository(db.DB)
	roomRepo := repository.NewRoomRepository(db.DB)
	messageRepo := repository.NewMessageRepository(db.DB)
	tagRepo := repository.NewTagRepository(db.DB)

	// Inicializar serviços
	userService := services.NewUserService(userRepo)
	roomService := services.NewRoomService(roomRepo)
	messageService := services.NewMessageService(messageRepo)
	tagService := services.NewTagService(tagRepo)

	// Inicializar hub WebSocket
	hub := websocket.NewHub()
	go hub.Run()

	// Inicializar controllers
	userController := controllers.NewUserController(userService)
	roomController := controllers.NewRoomController(roomService, userService, messageService)
	tagController := controllers.NewTagController(tagService)
	wsHandler := websocket.NewHandler(hub, userRepo, messageRepo, roomRepo)

	// Configurar Fiber
	app := fiber.New(fiber.Config{
		AppName: "Whatz Chat API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.RequestLogger())

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "*"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigin,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Rotas de health check
	// @Summary Informações da API
	// @Description Retorna informações básicas da API
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]interface{} "Informações da API"
	// @Router / [get]
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Whatz Chat API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// @Summary Status de saúde
	// @Description Verifica se a API está funcionando
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]interface{} "Status de saúde"
	// @Router /health [get]
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// API v1
	api := app.Group("/api/v1")

	// Rotas de usuários
	users := api.Group("/users")
	users.Post("/", userController.Create)
	users.Get("/", userController.GetAll)
	users.Get("/:id", userController.GetByID)
	users.Put("/:id", userController.Update)
	users.Delete("/:id", userController.Delete)

	// Rotas de salas
	rooms := api.Group("/rooms")
	rooms.Post("/", roomController.Create)
	rooms.Get("/", roomController.GetAll)
	rooms.Get("/public", roomController.GetPublicRooms) // Deve vir antes de /:id
	rooms.Get("/:id", roomController.GetByID)
	rooms.Get("/:id/messages", roomController.GetMessages)
	rooms.Put("/:id", roomController.Update)
	rooms.Delete("/:id", roomController.Delete)

	// Rotas de tags
	tags := api.Group("/tags")
	tags.Post("/", tagController.Create)
	tags.Get("/", tagController.GetAll)
	tags.Delete("/:id", tagController.Delete)

	// Rotas administrativas
	admin := api.Group("/admin")
	admin.Put("/users/:id/tags", userController.UpdateTags)
	admin.Put("/users/:id/role", userController.UpdateRole)
	admin.Get("/users/role/:role", userController.GetByRole)
	admin.Post("/rooms", roomController.CreateWithAccess)

	// WebSocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		if ws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", ws.New(wsHandler.HandleWebSocket))

	// Middleware 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint não encontrado",
		})
	})

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor iniciado na porta %s", port)
	log.Printf("📚 Documentação: http://localhost:%s/api/v1", port)
	log.Printf("💬 WebSocket: ws://localhost:%s/ws", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}
