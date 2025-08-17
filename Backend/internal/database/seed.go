package database

import (
	"log"
	"time"

	"github.com/rafael-bit/whatz/internal/models"
	"github.com/rafael-bit/whatz/internal/repository"
)

type Seeder struct {
	userRepo    *repository.UserRepository
	roomRepo    *repository.RoomRepository
	messageRepo *repository.MessageRepository
}

func NewSeeder(userRepo *repository.UserRepository, roomRepo *repository.RoomRepository, messageRepo *repository.MessageRepository) *Seeder {
	return &Seeder{
		userRepo:    userRepo,
		roomRepo:    roomRepo,
		messageRepo: messageRepo,
	}
}

func (s *Seeder) Seed() error {
	start := time.Now()
	log.Printf("🌱 Iniciando seed do banco de dados...")

	// Criar usuários de exemplo
	users, err := s.seedUsers()
	if err != nil {
		return err
	}

	// Criar salas de exemplo
	rooms, err := s.seedRooms(users)
	if err != nil {
		return err
	}

	// Criar mensagens de exemplo
	err = s.seedMessages(users, rooms)
	if err != nil {
		return err
	}

	log.Printf("✅ Seed concluído com sucesso em %v", time.Since(start))
	return nil
}

func (s *Seeder) seedUsers() ([]*models.User, error) {
	start := time.Now()
	log.Printf("👥 Criando usuários de exemplo...")

	users := []*models.User{
		models.NewUser("admin", "admin@whatz.com", "https://api.dicebear.com/7.x/avataaars/svg?seed=admin"),
		models.NewUser("alice", "alice@whatz.com", "https://api.dicebear.com/7.x/avataaars/svg?seed=alice"),
		models.NewUser("bob", "bob@whatz.com", "https://api.dicebear.com/7.x/avataaars/svg?seed=bob"),
		models.NewUser("charlie", "charlie@whatz.com", "https://api.dicebear.com/7.x/avataaars/svg?seed=charlie"),
		models.NewUser("diana", "diana@whatz.com", "https://api.dicebear.com/7.x/avataaars/svg?seed=diana"),
	}

	for _, user := range users {
		// Verificar se usuário já existe
		existingUser, err := s.userRepo.GetByUsername(user.Username)
		if err != nil {
			log.Printf("❌ Erro ao verificar usuário existente: %v", err)
			return nil, err
		}

		if existingUser == nil {
			if err := s.userRepo.Create(user); err != nil {
				log.Printf("❌ Erro ao criar usuário %s: %v", user.Username, err)
				return nil, err
			}
			log.Printf("✅ Usuário criado: %s", user.Username)
		} else {
			log.Printf("⚠️ Usuário já existe: %s", user.Username)
			user = existingUser
		}
	}

	log.Printf("✅ %d usuários processados em %v", len(users), time.Since(start))
	return users, nil
}

func (s *Seeder) seedRooms(users []*models.User) ([]*models.Room, error) {
	start := time.Now()
	log.Printf("🏠 Criando salas de exemplo...")

	rooms := []*models.Room{
		models.NewRoom("Geral", "Sala geral para discussões da equipe", "public", users[0].ID),
		models.NewRoom("Desenvolvimento", "Sala para discussões técnicas", "public", users[1].ID),
		models.NewRoom("Design", "Sala para discussões de design", "public", users[2].ID),
		models.NewRoom("Marketing", "Sala para estratégias de marketing", "public", users[3].ID),
		models.NewRoom("Off-topic", "Sala para conversas informais", "public", users[4].ID),
	}

	for _, room := range rooms {
		// Verificar se sala já existe
		existingRoom, err := s.roomRepo.GetByID(room.ID)
		if err != nil {
			log.Printf("❌ Erro ao verificar sala existente: %v", err)
			return nil, err
		}

		if existingRoom == nil {
			if err := s.roomRepo.Create(room); err != nil {
				log.Printf("❌ Erro ao criar sala %s: %v", room.Name, err)
				return nil, err
			}
			log.Printf("✅ Sala criada: %s", room.Name)
		} else {
			log.Printf("⚠️ Sala já existe: %s", room.Name)
			room = existingRoom
		}
	}

	log.Printf("✅ %d salas processadas em %v", len(rooms), time.Since(start))
	return rooms, nil
}

func (s *Seeder) seedMessages(users []*models.User, rooms []*models.Room) error {
	start := time.Now()
	log.Printf("💬 Criando mensagens de exemplo...")

	// Mensagens para a sala Geral
	generalMessages := []string{
		"Olá pessoal! Bem-vindos ao Whatz! 🚀",
		"Oi! Que legal essa nova plataforma de chat!",
		"Vamos começar a usar para nossas conversas da equipe",
		"Concordo! Muito mais organizado que outros chats",
		"Alguém já testou as funcionalidades?",
		"Sim! O WebSocket está funcionando perfeitamente",
		"E o histórico de mensagens também!",
		"Perfeito! Agora temos uma ferramenta profissional",
	}

	// Mensagens para a sala Desenvolvimento
	devMessages := []string{
		"Pessoal, vamos discutir a arquitetura do novo projeto",
		"Que tal usar Go com Fiber para o backend?",
		"Excelente escolha! Go é muito performático",
		"E para o frontend? React ou Vue?",
		"Eu sugiro React com TypeScript",
		"Concordo! TypeScript vai nos ajudar muito",
		"Vamos usar SQLite para desenvolvimento?",
		"Sim! Depois migramos para PostgreSQL em produção",
	}

	// Mensagens para a sala Design
	designMessages := []string{
		"Galera, precisamos definir o design system",
		"Que tal começarmos com as cores principais?",
		"Eu sugiro um tema escuro e claro",
		"Ótima ideia! Vamos usar CSS variables",
		"E para os componentes? Material Design?",
		"Não, vamos criar algo mais customizado",
		"Concordo! Vamos ter nossa própria identidade",
		"Perfeito! Vamos começar pelos botões e inputs",
	}

	// Criar mensagens para cada sala
	roomMessages := map[string][]string{
		rooms[0].ID: generalMessages,
		rooms[1].ID: devMessages,
		rooms[2].ID: designMessages,
	}

	messageCount := 0
	for roomID, messages := range roomMessages {
		for i, content := range messages {
			// Distribuir mensagens entre os usuários
			userIndex := i % len(users)
			user := users[userIndex]

			message := models.NewMessage(content, user.ID, user.Username, user.Avatar, "text", roomID)

			// Definir timestamp específico para ordenação
			message.CreatedAt = time.Now().Add(-time.Duration(len(messages)-i) * time.Minute)
			message.UpdatedAt = message.CreatedAt

			if err := s.messageRepo.Create(message); err != nil {
				log.Printf("❌ Erro ao criar mensagem: %v", err)
				return err
			}
			messageCount++
		}
	}

	log.Printf("✅ %d mensagens criadas em %v", messageCount, time.Since(start))
	return nil
}
