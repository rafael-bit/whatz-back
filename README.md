# Whatz - Chat em Tempo Real 🚀

Um sistema de chat moderno e responsivo desenvolvido em Go com WebSocket para comunicação em tempo real.

## 📚 Documentação

- **Swagger UI**: http://localhost:8080/swagger/
- **Documentação Completa**: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

## 🎯 Características

- **Backend robusto** em Go com Fiber
- **WebSocket** para comunicação em tempo real
- **Banco de dados SQLite** para persistência
- **API REST** completa para gerenciamento
- **Sistema de salas** públicas e privadas
- **Indicadores de digitação** em tempo real
- **Histórico de mensagens** com paginação
- **Logs detalhados** com tracing de performance
- **Dados de seed** para testes rápidos

## 🛠️ Stack Tecnológica

- **Backend**: Go 1.25+
- **Framework**: Fiber v2
- **WebSocket**: Fiber WebSocket
- **Banco de Dados**: SQLite
- **Logs**: Log padrão do Go com timestamps
- **Validação**: Validação manual de dados

## 📋 Pré-requisitos

- Go 1.25 ou superior
- Git

## 🚀 Instalação e Execução

### 1. Clone o repositório
```bash
git clone <repository-url>
cd whatz
```

### 2. Configure as variáveis de ambiente
```bash
cp env.example .env
# Edite o arquivo .env conforme necessário
```

### 3. Instale as dependências
```bash
go mod tidy
```

### 4. Execute o servidor
```bash
go run cmd/server/main.go
```

O servidor estará disponível em `http://localhost:8080`

## 📊 Endpoints da API

### Health Check
- `GET /` - Status do servidor
- `GET /health` - Health check detalhado

### Usuários
- `POST /api/v1/users` - Criar usuário
- `GET /api/v1/users` - Listar usuários
- `GET /api/v1/users/:id` - Buscar usuário
- `PUT /api/v1/users/:id` - Atualizar usuário
- `DELETE /api/v1/users/:id` - Deletar usuário

### Salas
- `POST /api/v1/rooms` - Criar sala
- `GET /api/v1/rooms` - Listar salas
- `GET /api/v1/rooms/public` - Salas públicas
- `GET /api/v1/rooms/:id` - Buscar sala
- `GET /api/v1/rooms/:id/messages` - Mensagens da sala
- `PUT /api/v1/rooms/:id` - Atualizar sala
- `DELETE /api/v1/rooms/:id` - Deletar sala

### WebSocket
- `WS /ws?user_id=X&room_id=Y` - Conexão WebSocket

## 💬 Protocolo WebSocket

### Tipos de Mensagens

#### Enviadas pelo Cliente
```json
{
  "type": "send_message",
  "payload": {
    "content": "Olá pessoal!"
  }
}
```

```json
{
  "type": "typing_start"
}
```

```json
{
  "type": "typing_stop"
}
```

#### Recebidas pelo Cliente
```json
{
  "type": "welcome",
  "payload": {
    "room": {
      "id": "room-id",
      "name": "Geral",
      "description": "Sala geral"
    },
    "online_users": [
      {
        "user_id": "user-id",
        "username": "alice"
      }
    ]
  }
}
```

```json
{
  "type": "new_message",
  "payload": {
    "id": "message-id",
    "content": "Olá pessoal!",
    "user_id": "user-id",
    "username": "alice",
    "avatar": "avatar-url",
    "type": "text",
    "room_id": "room-id",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

```json
{
  "type": "typing_indicator",
  "payload": {
    "user_id": "user-id",
    "username": "alice",
    "is_typing": true
  }
}
```

```json
{
  "type": "user_joined",
  "payload": {
    "user_id": "user-id",
    "username": "bob"
  }
}
```

```json
{
  "type": "user_left",
  "payload": {
    "user_id": "user-id",
    "username": "bob"
  }
}
```

```json
{
  "type": "message_history",
  "payload": [
    {
      "id": "message-id",
      "content": "Mensagem anterior",
      "user_id": "user-id",
      "username": "alice",
      "avatar": "avatar-url",
      "type": "text",
      "room_id": "room-id",
      "created_at": "2024-01-01T11:00:00Z"
    }
  ]
}
```

## 🗄️ Estrutura do Banco de Dados

### Tabela: users
- `id` (TEXT, PRIMARY KEY)
- `username` (TEXT, UNIQUE, NOT NULL)
- `email` (TEXT, UNIQUE, NOT NULL)
- `avatar` (TEXT)
- `status` (TEXT, DEFAULT 'online')
- `created_at` (DATETIME, NOT NULL)
- `updated_at` (DATETIME, NOT NULL)

### Tabela: rooms
- `id` (TEXT, PRIMARY KEY)
- `name` (TEXT, NOT NULL)
- `description` (TEXT)
- `type` (TEXT, DEFAULT 'public')
- `created_by` (TEXT, NOT NULL, FOREIGN KEY)
- `created_at` (DATETIME, NOT NULL)
- `updated_at` (DATETIME, NOT NULL)

### Tabela: messages
- `id` (TEXT, PRIMARY KEY)
- `content` (TEXT, NOT NULL)
- `user_id` (TEXT, NOT NULL, FOREIGN KEY)
- `username` (TEXT, NOT NULL)
- `avatar` (TEXT)
- `type` (TEXT, DEFAULT 'text')
- `room_id` (TEXT, NOT NULL, FOREIGN KEY)
- `created_at` (DATETIME, NOT NULL)
- `updated_at` (DATETIME, NOT NULL)

## 🌱 Dados de Seed

O sistema inclui dados de exemplo:

### Usuários
- admin (admin@whatz.com)
- alice (alice@whatz.com)
- bob (bob@whatz.com)
- charlie (charlie@whatz.com)
- diana (diana@whatz.com)

### Salas
- Geral (pública)
- Desenvolvimento (pública)
- Design (pública)
- Marketing (pública)
- Off-topic (pública)

### Mensagens
- Mensagens de exemplo em cada sala

## 📁 Estrutura do Projeto

```
whatz/
├── cmd/
│   └── server/
│       └── main.go          # Ponto de entrada da aplicação
├── internal/
│   ├── database/
│   │   ├── database.go      # Configuração do banco de dados
│   │   └── seed.go          # Dados de seed
│   ├── handlers/
│   │   ├── user_handler.go  # Handlers HTTP para usuários
│   │   └── room_handler.go  # Handlers HTTP para salas
│   ├── logger/
│   │   └── logger.go        # Middleware de logging
│   ├── models/
│   │   ├── user.go          # Modelo de usuário
│   │   ├── room.go          # Modelo de sala
│   │   └── message.go       # Modelo de mensagem
│   ├── repository/
│   │   ├── user_repository.go    # Repositório de usuários
│   │   ├── room_repository.go    # Repositório de salas
│   │   └── message_repository.go # Repositório de mensagens
│   └── websocket/
│       ├── hub.go           # Hub central do WebSocket
│       └── handlers.go      # Handlers do WebSocket
├── go.mod                   # Dependências do Go
├── go.sum                   # Checksums das dependências
├── env.example              # Exemplo de variáveis de ambiente
└── README.md                # Documentação
```

## 🔧 Configuração

### Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `PORT` | Porta do servidor | `8080` |
| `DB_PATH` | Caminho do banco SQLite | `./whatz.db` |
| `LOG_LEVEL` | Nível de log | `info` |
| `CORS_ORIGIN` | Origem permitida para CORS | `*` |

## 🧪 Testando a API

### Criar um usuário
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=test"
  }'
```

### Listar usuários
```bash
curl http://localhost:8080/api/v1/users
```

### Criar uma sala
```bash
curl -X POST http://localhost:8080/api/v1/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Teste",
    "description": "Sala de teste",
    "type": "public",
    "created_by": "user-id"
  }'
```

### Listar salas
```bash
curl http://localhost:8080/api/v1/rooms
```

## 🚀 Deploy

### Docker (Opcional)

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Build para produção
```bash
go build -o whatz-server cmd/server/main.go
./whatz-server
```

## 📈 Performance

- **WebSocket**: Suporte a múltiplas conexões simultâneas
- **Banco de dados**: Índices otimizados para consultas frequentes
- **Logs**: Tracing de performance em todas as operações
- **CORS**: Configurado para desenvolvimento e produção

## 🔒 Segurança

- Validação de dados de entrada
- Sanitização de parâmetros
- Controle de acesso por sala
- Logs de auditoria

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 🆘 Suporte

Para suporte, abra uma issue no repositório ou entre em contato com a equipe de desenvolvimento.

---

**Whatz** - Conectando pessoas em tempo real! 💬✨
