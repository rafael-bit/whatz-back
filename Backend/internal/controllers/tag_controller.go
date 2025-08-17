package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafael-bit/whatz/internal/models"
	"github.com/rafael-bit/whatz/internal/services"
)

type TagController struct {
	tagService *services.TagService
}

func NewTagController(tagService *services.TagService) *TagController {
	return &TagController{
		tagService: tagService,
	}
}

// Create godoc
// @Summary Criar nova tag
// @Description Cria uma nova tag no sistema
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body map[string]string true "Dados da tag"
// @Success 201 {object} map[string]interface{} "Tag criada com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /tags [post]
func (c *TagController) Create(ctx *fiber.Ctx) error {
	var req struct {
		Name string `json:"name" validate:"required"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	// Check if tag already exists
	existingTag, err := c.tagService.GetByName(req.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	if existingTag != nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Tag já existe",
		})
	}

	tag := models.NewTag(req.Name)
	if err := c.tagService.Create(tag); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao criar tag",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tag criada com sucesso",
		"tag":     tag,
	})
}

// GetAll godoc
// @Summary Listar todas as tags
// @Description Retorna uma lista de todas as tags do sistema
// @Tags tags
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Lista de tags"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /tags [get]
func (c *TagController) GetAll(ctx *fiber.Ctx) error {
	tags, err := c.tagService.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	// Ensure tags is never null
	if tags == nil {
		tags = []*models.Tag{}
	}

	return ctx.JSON(fiber.Map{
		"tags":  tags,
		"count": len(tags),
	})
}

// Delete godoc
// @Summary Deletar tag
// @Description Remove uma tag do sistema
// @Tags tags
// @Accept json
// @Produce json
// @Param id path string true "ID da tag"
// @Success 200 {object} map[string]interface{} "Tag deletada com sucesso"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /tags/{id} [delete]
func (c *TagController) Delete(ctx *fiber.Ctx) error {
	tagID := ctx.Params("id")

	if err := c.tagService.Delete(tagID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao deletar tag",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Tag deletada com sucesso",
	})
}
