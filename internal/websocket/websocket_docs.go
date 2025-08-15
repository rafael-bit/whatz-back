// Package websocket contém a implementação do WebSocket para chat em tempo real
//
// WebSocket Endpoint: ws://localhost:8080/ws
//
// Protocolo de Mensagens:
//
// 1. Conectar ao WebSocket:
//   - URL: ws://localhost:8080/ws
//   - Query Parameters:
//   - user_id: ID do usuário
//   - room_id: ID da sala (opcional)
//
// 2. Tipos de Mensagens:
//
//	a) Mensagem de Texto:
//	{
//	  "type": "message",
//	  "content": "Olá, mundo!",
//	  "room_id": "room-uuid",
//	  "user_id": "user-uuid"
//	}
//
//	b) Mensagem de Sistema:
//	{
//	  "type": "system",
//	  "content": "Usuário entrou na sala",
//	  "room_id": "room-uuid"
//	}
//
//	c) Mensagem de Status:
//	{
//	  "type": "status",
//	  "user_id": "user-uuid",
//	  "status": "online|offline|away"
//	}
//
// 3. Eventos Recebidos:
//
//	a) Mensagem Recebida:
//	{
//	  "type": "message",
//	  "id": "message-uuid",
//	  "content": "Conteúdo da mensagem",
//	  "user_id": "user-uuid",
//	  "username": "Nome do usuário",
//	  "avatar": "URL do avatar",
//	  "room_id": "room-uuid",
//	  "created_at": "2024-01-01T12:00:00Z"
//	}
//
//	b) Usuário Conectado:
//	{
//	  "type": "user_connected",
//	  "user_id": "user-uuid",
//	  "username": "Nome do usuário",
//	  "room_id": "room-uuid"
//	}
//
//	c) Usuário Desconectado:
//	{
//	  "type": "user_disconnected",
//	  "user_id": "user-uuid",
//	  "room_id": "room-uuid"
//	}
//
//	d) Erro:
//	{
//	  "type": "error",
//	  "message": "Descrição do erro"
//	}
//
// 4. Códigos de Status:
//   - 1000: Conexão normal
//   - 1001: Cliente desconectando
//   - 1002: Erro de protocolo
//   - 1003: Tipo de dados não suportado
//   - 1006: Conexão anormal
//   - 1009: Mensagem muito grande
//   - 1011: Erro interno do servidor
//
// 5. Exemplo de Uso com JavaScript:
//
//	const ws = new WebSocket('ws://localhost:8080/ws?user_id=123&room_id=456');
//
//	ws.onopen = function() {
//	  console.log('Conectado ao WebSocket');
//	};
//
//	ws.onmessage = function(event) {
//	  const data = JSON.parse(event.data);
//	  console.log('Mensagem recebida:', data);
//	};
//
//	ws.onclose = function() {
//	  console.log('Desconectado do WebSocket');
//	};
//
//	// Enviar mensagem
//	ws.send(JSON.stringify({
//	  type: 'message',
//	  content: 'Olá, mundo!',
//	  room_id: '456',
//	  user_id: '123'
//	}));
//
// swagger:meta
package websocket
