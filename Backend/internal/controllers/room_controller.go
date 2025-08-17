package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rafael-bit/whatz/internal/models"
	"github.com/rafael-bit/whatz/internal/services"
)

// CreateRoomRequest representa a requisição para criar uma sala
// @Description Dados necessários para criar uma nova sala
type CreateRoomRequest struct {
	// @Description Nome da sala (3-100 caracteres)
	// @Example "Sala Geral"
	Name string `json:"name" validate:"required,min=3,max=100"`
	// @Description Descrição da sala (máximo 500 caracteres)
	// @Example "Sala para conversas gerais"
	Description string `json:"description" validate:"max=500"`
	// @Description Tipo da sala (public ou private)
	// @Example "public"
	Type string `json:"type" validate:"oneof=public private"`
	// @Description ID do usuário que criou a sala
	// @Example "123e4567-e89b-12d3-a456-426614174000"
	CreatedBy string `json:"created_by" validate:"required"`
}

// UpdateRoomRequest representa a requisição para atualizar uma sala
// @Description Dados necessários para atualizar uma sala existente
type UpdateRoomRequest struct {
	// @Description Nome da sala (3-100 caracteres)
	// @Example "Sala Geral"
	Name string `json:"name" validate:"required,min=3,max=100"`
	// @Description Descrição da sala (máximo 500 caracteres)
	// @Example "Sala para conversas gerais"
	Description string `json:"description" validate:"max=500"`
	// @Description Tipo da sala (public ou private)
	// @Example "public"
	Type string `json:"type" validate:"oneof=public private"`
	// @Description Lista de tags de acesso
	// @Example ["dev", "admin"]
	AccessTags []string `json:"access_tags"`
}

// CreateRoomWithAccessRequest representa a requisição para criar uma sala com controle de acesso
// @Description Dados necessários para criar uma sala com controle de acesso (rota administrativa)
type CreateRoomWithAccessRequest struct {
	// @Description Nome da sala (3-100 caracteres)
	// @Example "Sala VIP"
	Name string `json:"name" validate:"required,min=3,max=100"`
	// @Description Descrição da sala (máximo 500 caracteres)
	// @Example "Sala exclusiva para membros VIP"
	Description string `json:"description" validate:"max=500"`
	// @Description Tipo da sala (public ou private)
	// @Example "private"
	Type string `json:"type" validate:"oneof=public private"`
	// @Description Lista de tags de acesso
	// @Example ["dev", "admin"]
	AccessTags []string `json:"access_tags"`
	// @Description ID do usuário que criou a sala
	// @Example "123e4567-e89b-12d3-a456-426614174000"
	CreatedBy string `json:"created_by" validate:"required"`
}

type RoomController struct {
	roomService    *services.RoomService
	userService    *services.UserService
	messageService *services.MessageService
}

func NewRoomController(roomService *services.RoomService, userService *services.UserService, messageService *services.MessageService) *RoomController {
	return &RoomController{
		roomService:    roomService,
		userService:    userService,
		messageService: messageService,
	}
}

// Create godoc
// @Summary Criar nova sala
// @Description Cria uma nova sala de chat
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body CreateRoomRequest true "Dados da sala"
// @Success 201 {object} map[string]interface{} "Sala criada com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rooms [post]
func (c *RoomController) Create(ctx *fiber.Ctx) error {
	var req CreateRoomRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	room := models.NewRoom(req.Name, req.Description, req.Type, req.CreatedBy)
	if err := c.roomService.Create(room); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao criar sala",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sala criada com sucesso",
		"room":    room,
	})
}

// GetByID godoc
// @Summary Buscar sala por ID
// @Description Retorna uma sala específica pelo ID
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "ID da sala"
// @Success 200 {object} map[string]interface{} "Sala encontrada"
// @Failure 404 {object} map[string]interface{} "Sala não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rooms/{id} [get]
func (c *RoomController) GetByID(ctx *fiber.Ctx) error {
	roomID := ctx.Params("id")

	room, err := c.roomService.GetByID(roomID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	if room == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Sala não encontrada",
		})
	}

	return ctx.JSON(fiber.Map{
		"room": room,
	})
}

// GetAll godoc
// @Summary Listar salas
// @Description Retorna lista de salas baseada no usuário e suas permissões
// @Tags rooms
// @Accept json
// @Produce json
// @Param user_id query string false "ID do usuário para filtrar salas acessíveis"
// @Success 200 {object} map[string]interface{} "Lista de salas"
// @Failure 404 {object} map[string]interface{} "Usuário não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rooms [get]
func (c *RoomController) GetAll(ctx *fiber.Ctx) error {
	userID := ctx.Query("user_id")
	if userID == "" {
		rooms, err := c.roomService.GetPublicRooms()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		}

		// Ensure rooms is never null
		if rooms == nil {
			rooms = []*models.Room{}
		}

		return ctx.JSON(fiber.Map{
			"rooms": rooms,
			"count": len(rooms),
			"type":  "public_only",
		})
	}

	// Special case for admin access
	if userID == "admin" {
		rooms, err := c.roomService.GetAll()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		}

		// Ensure rooms is never null
		if rooms == nil {
			rooms = []*models.Room{}
		}

		return ctx.JSON(fiber.Map{
			"rooms": rooms,
			"count": len(rooms),
			"type":  "admin_all",
		})
	}

	user, err := c.userService.GetByID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	if user == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Usuário não encontrado",
		})
	}

	var userTags []string
	if err := json.Unmarshal([]byte(user.Tags), &userTags); err != nil {
		userTags = []string{}
	}

	if user.Role == "admin" {
		rooms, err := c.roomService.GetAll()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		}

		// Ensure rooms is never null
		if rooms == nil {
			rooms = []*models.Room{}
		}

		return ctx.JSON(fiber.Map{
			"rooms": rooms,
			"count": len(rooms),
			"type":  "admin_all",
		})
	}

	rooms, err := c.roomService.GetRoomsByAccessTags(userTags)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	// Ensure rooms is never null
	if rooms == nil {
		rooms = []*models.Room{}
	}

	return ctx.JSON(fiber.Map{
		"rooms":     rooms,
		"count":     len(rooms),
		"type":      "user_accessible",
		"user_tags": userTags,
	})
}

// GetPublicRooms godoc
// @Summary Listar salas públicas
// @Description Retorna apenas as salas públicas do sistema
// @Tags rooms
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Lista de salas públicas"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rooms/public [get]
func (c *RoomController) GetPublicRooms(ctx *fiber.Ctx) error {
	rooms, err := c.roomService.GetPublicRooms()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	// Ensure rooms is never null
	if rooms == nil {
		rooms = []*models.Room{}
	}

	return ctx.JSON(fiber.Map{
		"rooms": rooms,
		"count": len(rooms),
	})
}

// GetMessages godoc
// @Summary Buscar mensagens de uma sala
// @Description Retorna as mensagens de uma sala específica com paginação
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "ID da sala"
// @Param limit query int false "Limite de mensagens (padrão: 50)"
// @Param offset query int false "Offset para paginação (padrão: 0)"
// @Success 200 {object} map[string]interface{} "Mensagens da sala"
// @Failure 404 {object} map[string]interface{} "Sala não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rooms/{id}/messages [get]
func (c *RoomController) GetMessages(ctx *fiber.Ctx) error {
	roomID := ctx.Params("id")
	limitStr := ctx.Query("limit", "50")
	offsetStr := ctx.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	room, err := c.roomService.GetByID(roomID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	if room == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Sala não encontrada",
		})
	}

	messages, err := c.messageService.GetByRoom(roomID, limit, offset)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	total, err := c.messageService.GetMessageCount(roomID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	return ctx.JSON(fiber.Map{
		"messages": messages,
		"room":     room,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	})
}

// Update godoc
// @Summary Atualizar sala
// @Description Atualiza os dados de uma sala existente
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "ID da sala"
// @Param room body UpdateRoomRequest true "Dados atualizados da sala"
// @Success 200 {object} map[string]interface{} "Sala atualizada com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 404 {object} map[string]interface{} "Sala não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rooms/{id} [put]
func (c *RoomController) Update(ctx *fiber.Ctx) error {
	roomID := ctx.Params("id")

	var req UpdateRoomRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	room, err := c.roomService.GetByID(roomID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	if room == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Sala não encontrada",
		})
	}

	// Only update fields that are provided
	if req.Name != "" {
		room.Name = req.Name
	}
	if req.Description != "" {
		room.Description = req.Description
	}
	if req.Type != "" {
		room.Type = req.Type
	}
	room.UpdatedAt = time.Now()

	// Update access tags if provided
	if req.AccessTags != nil {
		accessTagsJSON, err := json.Marshal(req.AccessTags)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		}
		room.AccessTags = string(accessTagsJSON)
	}

	if err := c.roomService.Update(room); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao atualizar sala",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Sala atualizada com sucesso",
		"room":    room,
	})
}

// Delete godoc
// @Summary Deletar sala
// @Description Remove uma sala do sistema
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "ID da sala"
// @Success 200 {object} map[string]interface{} "Sala deletada com sucesso"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rooms/{id} [delete]
func (c *RoomController) Delete(ctx *fiber.Ctx) error {
	roomID := ctx.Params("id")

	if err := c.roomService.Delete(roomID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao deletar sala",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Sala deletada com sucesso",
	})
}

// CreateWithAccess godoc
// @Summary Criar sala com controle de acesso
// @Description Cria uma nova sala com controle de acesso por tags (rota administrativa)
// @Tags admin
// @Accept json
// @Produce json
// @Param room body CreateRoomWithAccessRequest true "Dados da sala com controle de acesso"
// @Success 201 {object} map[string]interface{} "Sala criada com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /admin/rooms [post]
func (c *RoomController) CreateWithAccess(ctx *fiber.Ctx) error {
	var req CreateRoomWithAccessRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	accessTagsJSON, err := json.Marshal(req.AccessTags)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	room := &models.Room{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		AccessTags:  string(accessTagsJSON),
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := c.roomService.Create(room); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao criar sala",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sala criada com sucesso",
		"room":    room,
	})
}
