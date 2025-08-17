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

	// Check if user is already connected and disconnect previous connection
	h.hub.mutex.RLock()
	var existingClient *Client
	for client := range h.hub.clients {
		if client.UserID == user.ID {
			existingClient = client
			break
		}
	}
	h.hub.mutex.RUnlock()

	if existingClient != nil {
		log.Printf("⚠️ User %s is already connected, disconnecting previous connection", user.Username)

		// Disconnect existing client gracefully
		h.hub.unregister <- existingClient

		// Wait for disconnection to complete
		time.Sleep(500 * time.Millisecond)

		// Verify that the client was removed
		h.hub.mutex.RLock()
		stillConnected := false
		for client := range h.hub.clients {
			if client.UserID == user.ID {
				stillConnected = true
				log.Printf("⚠️ User %s still connected, forcing disconnection", user.Username)
				h.hub.unregister <- client
			}
		}
		h.hub.mutex.RUnlock()

		if stillConnected {
			// Wait a bit more to ensure cleanup
			time.Sleep(300 * time.Millisecond)
		}
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

				// Usar mutex para proteger escrita na conexão
				client.writeMutex.Lock()
				err := c.WriteMessage(fiberws.TextMessage, message)
				client.writeMutex.Unlock()

				if err != nil {
					log.Printf("❌ Erro ao enviar mensagem: %v", err)
					return
				}
			}
		}
	}()

	// Send welcome message and history
	log.Printf("📤 Sending welcome message to %s", user.Username)
	h.sendWelcomeMessage(client, room)

	log.Printf("📤 Sending message history to %s", user.Username)
	h.sendMessageHistory(client, room.ID)

	log.Printf("✅ Client %s connected successfully in %v", user.Username, time.Since(start))

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
		log.Printf("❌ Error serializing welcome message: %v", err)
		return
	}

	select {
	case client.Conn.Send <- data:
		log.Printf("✅ Welcome message sent to %s", client.Username)
	default:
		log.Printf("❌ Failed to send welcome message to %s (channel full)", client.Username)
	}
}

func (h *Handler) sendMessageHistory(client *Client, roomID string) {
	// Get last 50 messages
	messages, err := h.messageRepo.GetRecentMessages(roomID, 50)
	if err != nil {
		log.Printf("❌ Error fetching message history: %v", err)
		return
	}

	historyMessage := &WSMessage{
		Type:    "message_history",
		Payload: messages,
	}

	data, err := json.Marshal(historyMessage)
	if err != nil {
		log.Printf("❌ Error serializing message history: %v", err)
		return
	}

	select {
	case client.Conn.Send <- data:
		log.Printf("✅ Message history sent to %s (%d messages)", client.Username, len(messages))
	default:
		log.Printf("❌ Failed to send message history to %s (channel full)", client.Username)
	}
}
