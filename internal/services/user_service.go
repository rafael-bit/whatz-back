package services

import (
	"time"

	"github.com/rafael-bit/whatz/internal/models"
	"github.com/rafael-bit/whatz/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(user *models.User) error {
	return s.userRepo.Create(user)
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *UserService) GetAll() ([]*models.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) Update(user *models.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.UpdateStatus(user.ID, user.Status)
}

func (s *UserService) Delete(id string) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) UpdateTags(id, tags string) error {
	return s.userRepo.UpdateTags(id, tags)
}

func (s *UserService) UpdateRole(id, role string) error {
	return s.userRepo.UpdateRole(id, role)
}

func (s *UserService) GetByRole(role string) ([]*models.User, error) {
	return s.userRepo.GetByRole(role)
}

func (s *UserService) UpdateStatus(id, status string) error {
	return s.userRepo.UpdateStatus(id, status)
}
