package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/rafael-bit/whatz/internal/models"
)

type Client struct {
	ID       string
	UserID   string
	Username string
	RoomID   string
	Conn     *Connection
	Hub      *Hub
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *models.Message
	register   chan *Client
	unregister chan *Client
	rooms      map[string]map[*Client]bool
	mutex      sync.RWMutex
}

type Connection struct {
	Send chan []byte
}

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	log.Printf("🚀 Hub WebSocket iniciado")

	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true

			// Adicionar cliente à sala
			if client.RoomID != "" {
				if h.rooms[client.RoomID] == nil {
					h.rooms[client.RoomID] = make(map[*Client]bool)
				}
				h.rooms[client.RoomID][client] = true
			}
			h.mutex.Unlock()

			log.Printf("✅ Cliente %s (%s) conectado na sala %s", client.Username, client.ID, client.RoomID)

			// Enviar mensagem de sistema sobre novo usuário
			h.broadcastToRoom(client.RoomID, &WSMessage{
				Type: "user_joined",
				Payload: map[string]string{
					"user_id":  client.UserID,
					"username": client.Username,
				},
			})

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Conn.Send)

				// Remover cliente da sala
				if client.RoomID != "" && h.rooms[client.RoomID] != nil {
					delete(h.rooms[client.RoomID], client)
					if len(h.rooms[client.RoomID]) == 0 {
						delete(h.rooms, client.RoomID)
					}
				}
			}
			h.mutex.Unlock()

			log.Printf("❌ Cliente %s (%s) desconectado da sala %s", client.Username, client.ID, client.RoomID)

			// Enviar mensagem de sistema sobre usuário que saiu
			if client.RoomID != "" {
				h.broadcastToRoom(client.RoomID, &WSMessage{
					Type: "user_left",
					Payload: map[string]string{
						"user_id":  client.UserID,
						"username": client.Username,
					},
				})
			}

		case message := <-h.broadcast:
			h.mutex.RLock()
			roomClients := h.rooms[message.RoomID]
			h.mutex.RUnlock()

			if roomClients != nil {
				wsMessage := &WSMessage{
					Type:    "new_message",
					Payload: message,
				}

				data, err := json.Marshal(wsMessage)
				if err != nil {
					log.Printf("❌ Erro ao serializar mensagem: %v", err)
					continue
				}

				for client := range roomClients {
					select {
					case client.Conn.Send <- data:
					default:
						close(client.Conn.Send)
						delete(roomClients, client)
					}
				}

				log.Printf("📤 Mensagem enviada para %d clientes na sala %s", len(roomClients), message.RoomID)
			}
		}
	}
}

func (h *Hub) broadcastToRoom(roomID string, message *WSMessage) {
	h.mutex.RLock()
	roomClients := h.rooms[roomID]
	h.mutex.RUnlock()

	if roomClients != nil {
		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("❌ Erro ao serializar mensagem de sistema: %v", err)
			return
		}

		for client := range roomClients {
			select {
			case client.Conn.Send <- data:
			default:
				close(client.Conn.Send)
				delete(roomClients, client)
			}
		}
	}
}

func (h *Hub) GetRoomClients(roomID string) []*Client {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	roomClients := h.rooms[roomID]
	if roomClients == nil {
		return []*Client{}
	}

	clients := make([]*Client, 0, len(roomClients))
	for client := range roomClients {
		clients = append(clients, client)
	}

	return clients
}

func (h *Hub) GetOnlineUsers(roomID string) []map[string]string {
	clients := h.GetRoomClients(roomID)
	users := make([]map[string]string, 0, len(clients))

	for _, client := range clients {
		users = append(users, map[string]string{
			"user_id":  client.UserID,
			"username": client.Username,
		})
	}

	return users
}

func (h *Hub) SendTypingIndicator(roomID, userID, username string, isTyping bool) {
	message := &WSMessage{
		Type: "typing_indicator",
		Payload: map[string]interface{}{
			"user_id":   userID,
			"username":  username,
			"is_typing": isTyping,
		},
	}

	h.broadcastToRoom(roomID, message)
}
