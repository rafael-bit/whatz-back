package services

import (
	"github.com/rafael-bit/whatz/internal/models"
	"github.com/rafael-bit/whatz/internal/repository"
)

type TagService struct {
	tagRepo *repository.TagRepository
}

func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

func (s *TagService) Create(tag *models.Tag) error {
	return s.tagRepo.Create(tag)
}

func (s *TagService) GetAll() ([]*models.Tag, error) {
	return s.tagRepo.GetAll()
}

func (s *TagService) GetByName(name string) (*models.Tag, error) {
	return s.tagRepo.GetByName(name)
}

func (s *TagService) Delete(id string) error {
	return s.tagRepo.Delete(id)
}
