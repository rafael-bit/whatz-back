package services

import (
	"github.com/rafael-bit/whatz/internal/models"
	"github.com/rafael-bit/whatz/internal/repository"
)

type RoomService struct {
	roomRepo *repository.RoomRepository
}

func NewRoomService(roomRepo *repository.RoomRepository) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
	}
}

func (s *RoomService) Create(room *models.Room) error {
	return s.roomRepo.Create(room)
}

func (s *RoomService) GetByID(id string) (*models.Room, error) {
	return s.roomRepo.GetByID(id)
}

func (s *RoomService) GetAll() ([]*models.Room, error) {
	return s.roomRepo.GetAll()
}

func (s *RoomService) GetPublicRooms() ([]*models.Room, error) {
	return s.roomRepo.GetPublicRooms()
}

func (s *RoomService) GetByCreator(createdBy string) ([]*models.Room, error) {
	return s.roomRepo.GetByCreator(createdBy)
}

func (s *RoomService) Update(room *models.Room) error {
	return s.roomRepo.Update(room)
}

func (s *RoomService) Delete(id string) error {
	return s.roomRepo.Delete(id)
}

func (s *RoomService) GetRoomsByAccessTags(userTags []string) ([]*models.Room, error) {
	return s.roomRepo.GetRoomsByAccessTags(userTags)
}
