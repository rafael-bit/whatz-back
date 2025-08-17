# 📚 Whatz Chat API - Documentação Completa

## 🎯 Visão Geral

A API do Whatz Chat é um sistema RESTful completo para chat em tempo real, desenvolvido em Go com WebSocket para comunicação instantânea.

## 🔗 Base URL

```
http://localhost:8080/api/v1
```

## 🔐 Autenticação

A API utiliza autenticação baseada em JWT. Inclua o token no header:

```
Authorization: Bearer <seu-token-jwt>
```

## 📊 Códigos de Resposta

| Código | Descrição |
|--------|-----------|
| 200 | Sucesso |
| 201 | Criado com sucesso |
| 400 | Dados inválidos |
| 401 | Não autorizado |
| 404 | Não encontrado |
| 500 | Erro interno do servidor |

## 👥 Usuários

### Listar Todos os Usuários

```http
GET /users
```

**Resposta:**
```json
{
  "success": true,
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "username": "joao123",
      "email": "joao@example.com",
      "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123",
      "status": "online",
      "role": "user",
      "tags": "vip,premium",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Criar Usuário

```http
POST /users
Content-Type: application/json
```

**Body:**
```json
{
  "username": "joao123",
  "email": "joao@example.com",
  "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123"
}
```

**Resposta:**
```json
{
  "success": true,
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "joao123",
    "email": "joao@example.com",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123",
    "status": "online",
    "role": "user",
    "tags": "",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Buscar Usuário por ID

```http
GET /users/{id}
```

### Atualizar Usuário

```http
PUT /users/{id}
Content-Type: application/json
```

**Body:**
```json
{
  "username": "joao123_updated",
  "email": "joao.updated@example.com",
  "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123_updated",
  "status": "online"
}
```

### Deletar Usuário

```http
DELETE /users/{id}
```

## 🏠 Salas

### Listar Salas

```http
GET /rooms?user_id={user_id}
```

**Parâmetros:**
- `user_id` (opcional): ID do usuário para filtrar salas acessíveis

**Resposta:**
```json
{
  "success": true,
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Sala Geral",
      "description": "Sala para conversas gerais",
      "type": "public",
      "access_tags": "",
      "created_by": "user-id",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Listar Salas Públicas

```http
GET /rooms/public
```

### Criar Sala

```http
POST /rooms
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Sala Geral",
  "description": "Sala para conversas gerais",
  "type": "public",
  "created_by": "user-id"
}
```

### Buscar Sala por ID

```http
GET /rooms/{id}
```

### Atualizar Sala

```http
PUT /rooms/{id}
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Sala Geral Atualizada",
  "description": "Sala atualizada para conversas gerais",
  "type": "public"
}
```

### Deletar Sala

```http
DELETE /rooms/{id}
```

### Mensagens da Sala

```http
GET /rooms/{id}/messages?limit=50&offset=0
```

**Parâmetros:**
- `limit` (opcional): Limite de mensagens (padrão: 50)
- `offset` (opcional): Offset para paginação (padrão: 0)

**Resposta:**
```json
{
  "success": true,
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "content": "Olá, mundo!",
      "user_id": "user-id",
      "username": "joao123",
      "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123",
      "type": "text",
      "room_id": "room-id",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## 🔐 Administração

### Criar Sala com Controle de Acesso

```http
POST /admin/rooms
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Sala VIP",
  "description": "Sala exclusiva para membros VIP",
  "type": "private",
  "access_tags": ["vip", "premium"],
  "created_by": "admin-user-id"
}
```

### Usuários por Role

```http
GET /admin/users/role/{role}
```

### Atualizar Role do Usuário

```http
PUT /admin/users/{id}/role
Content-Type: application/json
```

**Body:**
```json
{
  "role": "admin"
}
```

### Atualizar Tags do Usuário

```http
PUT /admin/users/{id}/tags
Content-Type: application/json
```

**Body:**
```json
{
  "tags": ["vip", "premium"]
}
```

## 🏷️ Tags

### Listar Tags

```http
GET /tags
```

### Criar Tag

```http
POST /tags
Content-Type: application/json
```

**Body:**
```json
{
  "name": "vip"
}
```

### Deletar Tag

```http
DELETE /tags/{id}
```

## 🌐 WebSocket

### Conexão

```
WS ws://localhost:8080/ws?user_id={user_id}&room_id={room_id}
```

### Eventos

#### Cliente → Servidor

**Enviar Mensagem:**
```json
{
  "type": "send_message",
  "payload": {
    "content": "Olá, mundo!"
  }
}
```

**Iniciar Digitação:**
```json
{
  "type": "typing_start",
  "payload": {}
}
```

**Parar Digitação:**
```json
{
  "type": "typing_stop",
  "payload": {}
}
```

#### Servidor → Cliente

**Nova Mensagem:**
```json
{
  "type": "new_message",
  "payload": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "content": "Olá, mundo!",
    "user_id": "user-id",
    "username": "joao123",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123",
    "type": "text",
    "room_id": "room-id",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Indicador de Digitação:**
```json
{
  "type": "typing_indicator",
  "payload": {
    "user_id": "user-id",
    "username": "joao123",
    "is_typing": true
  }
}
```

**Usuário Entrou:**
```json
{
  "type": "user_joined",
  "payload": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "joao123",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123"
  }
}
```

**Usuário Saiu:**
```json
{
  "type": "user_left",
  "payload": "user-id"
}
```

**Histórico de Mensagens:**
```json
{
  "type": "message_history",
  "payload": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "content": "Mensagem antiga",
      "user_id": "user-id",
      "username": "joao123",
      "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123",
      "type": "text",
      "room_id": "room-id",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**Boas-vindas:**
```json
{
  "type": "welcome",
  "payload": {
    "message": "Bem-vindo à sala!",
    "online_users": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "username": "joao123",
        "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123"
      }
    ]
  }
}
```

## 🎯 Exemplos de Uso

### JavaScript/TypeScript

```javascript
// Criar usuário
const createUser = async () => {
  const response = await fetch('http://localhost:8080/api/v1/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      username: 'joao123',
      email: 'joao@example.com',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=joao123'
    })
  });
  
  const data = await response.json();
  return data;
};

// WebSocket
const connectWebSocket = (userId, roomId) => {
  const ws = new WebSocket(`ws://localhost:8080/ws?user_id=${userId}&room_id=${roomId}`);
  
  ws.onopen = () => {
    console.log('Conectado ao WebSocket');
  };
  
  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Mensagem recebida:', data);
  };
  
  ws.onclose = () => {
    console.log('Desconectado do WebSocket');
  };
  
  return ws;
};
```

### cURL

```bash
# Criar usuário
curl -X POST "http://localhost:8080/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "joao123",
    "email": "joao@example.com",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123"
  }'

# Listar salas
curl -X GET "http://localhost:8080/api/v1/rooms?user_id=user-id"

# Criar sala
curl -X POST "http://localhost:8080/api/v1/rooms" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sala Geral",
    "description": "Sala para conversas gerais",
    "type": "public",
    "created_by": "user-id"
  }'
```

## 🔧 Health Check

### Status da API

```http
GET /health
```

**Resposta:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "2.0.0"
}
```

### Informações da API

```http
GET /
```

**Resposta:**
```json
{
  "message": "Whatz Chat API",
  "version": "2.0.0",
  "status": "running",
  "documentation": "http://localhost:8080/swagger/index.html"
}
```

## 🚨 Tratamento de Erros

### Formato de Erro

```json
{
  "error": "Descrição do erro",
  "code": "ERROR_CODE",
  "details": {
    "field": "Detalhes específicos do campo"
  }
}
```

### Códigos de Erro Comuns

| Código | Descrição |
|--------|-----------|
| `VALIDATION_ERROR` | Erro de validação de dados |
| `USER_NOT_FOUND` | Usuário não encontrado |
| `ROOM_NOT_FOUND` | Sala não encontrada |
| `ACCESS_DENIED` | Acesso negado |
| `INVALID_TOKEN` | Token JWT inválido |
| `WEBSOCKET_ERROR` | Erro de conexão WebSocket |

## 📊 Rate Limiting

A API implementa rate limiting para proteger contra spam:

- **Limite**: 100 requisições por minuto por IP
- **Headers de resposta**:
  - `X-RateLimit-Limit`: Limite de requisições
  - `X-RateLimit-Remaining`: Requisições restantes
  - `X-RateLimit-Reset`: Tempo para reset

## 🔗 Links Úteis

- [Swagger UI](http://localhost:8080/swagger/index.html)
- [Repositório GitHub](https://github.com/rafael-bit/whatz-chat)
- [Frontend](http://localhost:3001)

---

**Whatz Chat API** - Sistema de Chat em Tempo Real 🚀
