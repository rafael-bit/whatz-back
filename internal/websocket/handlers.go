package websocket

import (
	"encoding/json"
	"log"
	"time"

	fiberws "github.com/gofiber/websocket/v2"
	"github.com/rafael-bit/whatz/internal/models"
	"github.com/rafael-bit/whatz/internal/repository"
)

type Handler struct {
	hub         *Hub
	userRepo    *repository.UserRepository
	messageRepo *repository.MessageRepository
	roomRepo    *repository.RoomRepository
}

func NewHandler(hub *Hub, userRepo *repository.UserRepository, messageRepo *repository.MessageRepository, roomRepo *repository.RoomRepository) *Handler {
	return &Handler{
		hub:         hub,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
	}
}

func (h *Handler) HandleWebSocket(c *fiberws.Conn) {
	start := time.Now()
	log.Printf("📡 Nova conexão WebSocket estabelecida: %s", c.RemoteAddr())

	// Extrair parâmetros da query string
	userID := c.Query("user_id")
	roomID := c.Query("room_id")

	if userID == "" || roomID == "" {
		log.Printf("❌ Parâmetros obrigatórios não fornecidos: user_id=%s, room_id=%s", userID, roomID)
		c.Close()
		return
	}

	// Buscar usuário
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		log.Printf("❌ Erro ao buscar usuário: %v", err)
		c.Close()
		return
	}

	if user == nil {
		log.Printf("❌ Usuário não encontrado: %s", userID)
		c.Close()
		return
	}

	// Buscar sala
	room, err := h.roomRepo.GetByID(roomID)
	if err != nil {
		log.Printf("❌ Erro ao buscar sala: %v", err)
		c.Close()
		return
	}

	if room == nil {
		log.Printf("❌ Sala não encontrada: %s", roomID)
		c.Close()
		return
	}

	// Criar cliente
	client := &Client{
		ID:       c.RemoteAddr().String(),
		UserID:   user.ID,
		Username: user.Username,
		RoomID:   room.ID,
		Conn: &Connection{
			Send: make(chan []byte, 256),
		},
		Hub: h.hub,
	}

	// Registrar cliente no hub
	h.hub.register <- client

	// Atualizar status do usuário para online
	h.userRepo.UpdateStatus(user.ID, "online")

	// Goroutine para enviar mensagens para o cliente
	go func() {
		defer func() {
			h.hub.unregister <- client
			c.Close()
		}()

		for {
			select {
			case message, ok := <-client.Conn.Send:
				if !ok {
					return
				}
				if err := c.WriteMessage(fiberws.TextMessage, message); err != nil {
					log.Printf("❌ Erro ao enviar mensagem: %v", err)
					return
				}
			}
		}
	}()

	// Enviar mensagens de boas-vindas e histórico
	h.sendWelcomeMessage(client, room)
	h.sendMessageHistory(client, room.ID)

	log.Printf("✅ Cliente %s conectado com sucesso em %v", user.Username, time.Since(start))

	// Loop principal para receber mensagens
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("❌ Erro ao ler mensagem: %v", err)
			break
		}

		if messageType == fiberws.TextMessage {
			h.handleMessage(client, message)
		}
	}

	// Atualizar status do usuário para offline
	h.userRepo.UpdateStatus(user.ID, "offline")
	log.Printf("❌ Cliente %s desconectado", user.Username)
}

func (h *Handler) handleMessage(client *Client, message []byte) {
	start := time.Now()
	log.Printf("📨 Mensagem recebida de %s: %s", client.Username, string(message))

	var wsMessage WSMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		log.Printf("❌ Erro ao deserializar mensagem: %v", err)
		return
	}

	switch wsMessage.Type {
	case "send_message":
		h.handleSendMessage(client, wsMessage.Payload)
	case "typing_start":
		h.handleTypingStart(client)
	case "typing_stop":
		h.handleTypingStop(client)
	default:
		log.Printf("⚠️ Tipo de mensagem desconhecido: %s", wsMessage.Type)
	}

	log.Printf("✅ Mensagem processada em %v", time.Since(start))
}

func (h *Handler) handleSendMessage(client *Client, payload interface{}) {
	start := time.Now()

	// Converter payload para map
	payloadMap, ok := payload.(map[string]interface{})
	if !ok {
		log.Printf("❌ Payload inválido para mensagem")
		return
	}

	content, ok := payloadMap["content"].(string)
	if !ok || content == "" {
		log.Printf("❌ Conteúdo da mensagem inválido")
		return
	}

	// Criar nova mensagem
	message := models.NewMessage(content, client.UserID, client.Username, "", "text", client.RoomID)

	// Salvar no banco de dados
	if err := h.messageRepo.Create(message); err != nil {
		log.Printf("❌ Erro ao salvar mensagem: %v", err)
		return
	}

	// Broadcast para todos os clientes na sala
	h.hub.broadcast <- message

	log.Printf("✅ Mensagem enviada com sucesso em %v", time.Since(start))
}

func (h *Handler) handleTypingStart(client *Client) {
	h.hub.SendTypingIndicator(client.RoomID, client.UserID, client.Username, true)
}

func (h *Handler) handleTypingStop(client *Client) {
	h.hub.SendTypingIndicator(client.RoomID, client.UserID, client.Username, false)
}

func (h *Handler) sendWelcomeMessage(client *Client, room *models.Room) {
	welcomeMessage := &WSMessage{
		Type: "welcome",
		Payload: map[string]interface{}{
			"room": map[string]interface{}{
				"id":          room.ID,
				"name":        room.Name,
				"description": room.Description,
			},
			"online_users": h.hub.GetOnlineUsers(room.ID),
		},
	}

	data, err := json.Marshal(welcomeMessage)
	if err != nil {
		log.Printf("❌ Erro ao serializar mensagem de boas-vindas: %v", err)
		return
	}

	client.Conn.Send <- data
}

func (h *Handler) sendMessageHistory(client *Client, roomID string) {
	// Buscar últimas 50 mensagens
	messages, err := h.messageRepo.GetRecentMessages(roomID, 50)
	if err != nil {
		log.Printf("❌ Erro ao buscar histórico de mensagens: %v", err)
		return
	}

	historyMessage := &WSMessage{
		Type:    "message_history",
		Payload: messages,
	}

	data, err := json.Marshal(historyMessage)
	if err != nil {
		log.Printf("❌ Erro ao serializar histórico de mensagens: %v", err)
		return
	}

	client.Conn.Send <- data
}
