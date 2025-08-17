package controllers

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rafael-bit/whatz/internal/models"
	"github.com/rafael-bit/whatz/internal/services"
)

// CreateUserRequest representa a requisição para criar um usuário
// @Description Dados necessários para criar um novo usuário
type CreateUserRequest struct {
	// @Description Nome de usuário (3-50 caracteres)
	// @Example "joao123"
	Username string `json:"username" validate:"required,min=3,max=50"`
	// @Description Email do usuário
	// @Example "joao@example.com"
	Email string `json:"email" validate:"required,email"`
	// @Description URL do avatar do usuário
	// @Example "https://example.com/avatar.jpg"
	Avatar string `json:"avatar"`
	// @Description Função do usuário (user ou admin)
	// @Example "user"
	Role string `json:"role"`
	// @Description Lista de tags do usuário
	// @Example ["vip", "premium"]
	Tags []string `json:"tags"`
}

// UpdateUserRequest representa a requisição para atualizar um usuário
// @Description Dados necessários para atualizar um usuário existente
type UpdateUserRequest struct {
	// @Description Nome de usuário (3-50 caracteres)
	// @Example "joao123"
	Username string `json:"username" validate:"required,min=3,max=50"`
	// @Description Email do usuário
	// @Example "joao@example.com"
	Email string `json:"email" validate:"required,email"`
	// @Description URL do avatar do usuário
	// @Example "https://example.com/avatar.jpg"
	Avatar string `json:"avatar"`
	// @Description Status do usuário (online, offline, away)
	// @Example "online"
	Status string `json:"status" validate:"oneof=online offline away"`
	// @Description Função do usuário (user ou admin)
	// @Example "user"
	Role string `json:"role" validate:"oneof=user admin"`
	// @Description Lista de tags do usuário
	// @Example ["dev", "admin"]
	Tags []string `json:"tags"`
}

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Create godoc
// @Summary Criar novo usuário
// @Description Cria um novo usuário no sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "Dados do usuário"
// @Success 201 {object} map[string]interface{} "Usuário criado com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /users [post]
func (c *UserController) Create(ctx *fiber.Ctx) error {
	var req CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	user := models.NewUser(req.Username, req.Email, req.Avatar)
	user.Role = req.Role

	// Convert tags to JSON string
	if len(req.Tags) > 0 {
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		}
		user.Tags = string(tagsJSON)
	}

	if err := c.userService.Create(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao criar usuário",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Usuário criado com sucesso",
		"user":    user,
	})
}

// GetByID godoc
// @Summary Buscar usuário por ID
// @Description Retorna um usuário específico pelo ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} map[string]interface{} "Usuário encontrado"
// @Failure 404 {object} map[string]interface{} "Usuário não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /users/{id} [get]
func (c *UserController) GetByID(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

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

	return ctx.JSON(fiber.Map{
		"user": user,
	})
}

// GetAll godoc
// @Summary Listar todos os usuários
// @Description Retorna uma lista de todos os usuários do sistema
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Lista de usuários"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /users [get]
func (c *UserController) GetAll(ctx *fiber.Ctx) error {
	users, err := c.userService.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	// Ensure users is never null
	if users == nil {
		users = []*models.User{}
	}

	return ctx.JSON(fiber.Map{
		"users": users,
		"count": len(users),
	})
}

// Update godoc
// @Summary Atualizar usuário
// @Description Atualiza os dados de um usuário existente
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param user body UpdateUserRequest true "Dados atualizados do usuário"
// @Success 200 {object} map[string]interface{} "Usuário atualizado com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 404 {object} map[string]interface{} "Usuário não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /users/{id} [put]
func (c *UserController) Update(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	var req UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
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

	user.Username = req.Username
	user.Email = req.Email
	user.Avatar = req.Avatar
	user.Status = req.Status
	user.Role = req.Role
	user.UpdatedAt = time.Now()

	// Convert tags to JSON string
	if len(req.Tags) > 0 {
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro interno do servidor",
			})
		}
		user.Tags = string(tagsJSON)
	}

	if err := c.userService.Update(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao atualizar usuário",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Usuário atualizado com sucesso",
		"user":    user,
	})
}

// Delete godoc
// @Summary Deletar usuário
// @Description Remove um usuário do sistema
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {object} map[string]interface{} "Usuário deletado com sucesso"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /users/{id} [delete]
func (c *UserController) Delete(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	if err := c.userService.Delete(userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao deletar usuário",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Usuário deletado com sucesso",
	})
}

// UpdateTags godoc
// @Summary Atualizar tags do usuário
// @Description Atualiza as tags de um usuário (rota administrativa)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param tags body object true "Lista de tags"
// @Success 200 {object} map[string]interface{} "Tags atualizadas com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /admin/users/{id}/tags [put]
func (c *UserController) UpdateTags(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	var req struct {
		Tags []string `json:"tags" validate:"required"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	if err := c.userService.UpdateTags(userID, string(tagsJSON)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao atualizar tags",
		})
	}

	updatedUser, _ := c.userService.GetByID(userID)

	return ctx.JSON(fiber.Map{
		"message": "Tags atualizadas com sucesso",
		"user":    updatedUser,
	})
}

// UpdateRole godoc
// @Summary Atualizar role do usuário
// @Description Atualiza a role de um usuário (rota administrativa)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param role body object true "Nova role do usuário"
// @Success 200 {object} map[string]interface{} "Role atualizada com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /admin/users/{id}/role [put]
func (c *UserController) UpdateRole(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	var req struct {
		Role string `json:"role" validate:"required,oneof=admin user"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	if err := c.userService.UpdateRole(userID, req.Role); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao atualizar role",
		})
	}

	updatedUser, _ := c.userService.GetByID(userID)

	return ctx.JSON(fiber.Map{
		"message": "Role atualizada com sucesso",
		"user":    updatedUser,
	})
}

// GetByRole godoc
// @Summary Buscar usuários por role
// @Description Retorna todos os usuários com uma role específica
// @Tags admin
// @Accept json
// @Produce json
// @Param role path string true "Role dos usuários"
// @Success 200 {object} map[string]interface{} "Lista de usuários por role"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /admin/users/role/{role} [get]
func (c *UserController) GetByRole(ctx *fiber.Ctx) error {
	role := ctx.Params("role")

	users, err := c.userService.GetByRole(role)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	return ctx.JSON(fiber.Map{
		"users": users,
		"count": len(users),
		"role":  role,
	})
}
