# 🚀 Whatz Chat - Backend API

Sistema de chat em tempo real desenvolvido em Go com WebSocket, GORM e Gin Framework.

## 📋 Índice

- [Funcionalidades](#-funcionalidades)
- [Tecnologias](#-tecnologias)
- [Estrutura do Projeto](#-estrutura-do-projeto)
- [Instalação](#-instalação)
- [Configuração](#-configuração)
- [Uso](#-uso)
- [API Endpoints](#-api-endpoints)
- [WebSocket](#-websocket)
- [Documentação](#-documentação)
- [Desenvolvimento](#-desenvolvimento)
- [Testes](#-testes)
- [Deploy](#-deploy)

## ✨ Funcionalidades

### 🎯 Principais
- **Chat em Tempo Real**: Comunicação instantânea via WebSocket
- **Sistema de Salas**: Salas públicas e privadas
- **Controle de Acesso**: Baseado em tags e roles
- **Gestão de Usuários**: CRUD completo
- **Interface Administrativa**: Painel de administração
- **Persistência de Dados**: SQLite com GORM
- **API RESTful**: Endpoints bem documentados

### 🔐 Segurança
- **Autenticação JWT**: Tokens seguros
- **Controle de Acesso**: Roles (admin/user)
- **Validação de Dados**: Middleware de validação
- **CORS**: Configurado para frontend
- **Rate Limiting**: Proteção contra spam

### 📊 Monitoramento
- **Logs Estruturados**: Logging centralizado
- **Métricas**: Health checks
- **Documentação**: Swagger/OpenAPI
- **Error Handling**: Tratamento de erros robusto

## 🛠 Tecnologias

### Core
- **Go 1.21+**: Linguagem principal
- **Gin**: Framework web
- **GORM**: ORM para banco de dados
- **SQLite**: Banco de dados
- **WebSocket**: Comunicação em tempo real

### Utilitários
- **Swaggo**: Documentação automática
- **JWT-Go**: Autenticação
- **Validator**: Validação de dados
- **Cors**: Cross-origin resource sharing

### Desenvolvimento
- **Air**: Hot reload
- **Go Modules**: Gerenciamento de dependências
- **Make**: Automação de tarefas

## 📁 Estrutura do Projeto

```
Backend/
├── cmd/
│   └── server/
│       └── main.go              # Ponto de entrada
├── internal/
│   ├── controllers/             # Controladores HTTP
│   │   ├── user_controller.go
│   │   ├── room_controller.go
│   │   ├── tag_controller.go
│   │   └── message_controller.go
│   ├── services/                # Lógica de negócio
│   │   ├── user_service.go
│   │   ├── room_service.go
│   │   ├── tag_service.go
│   │   └── message_service.go
│   ├── repository/              # Camada de dados
│   │   ├── user_repository.go
│   │   ├── room_repository.go
│   │   ├── tag_repository.go
│   │   └── message_repository.go
│   ├── models/                  # Modelos de dados
│   │   ├── user.go
│   │   ├── room.go
│   │   ├── tag.go
│   │   └── message.go
│   ├── database/                # Configuração do banco
│   │   ├── database.go
│   │   └── seed.go
│   ├── websocket/               # WebSocket handlers
│   │   ├── hub.go
│   │   └── handlers.go
│   ├── handlers/                # Handlers HTTP
│   └── logger/                  # Sistema de logs
│       └── logger.go
├── docs/                        # Documentação Swagger
│   ├── docs.go
│   └── swagger.json
├── scripts/                     # Scripts utilitários
│   ├── clear_db.go
│   └── generate-docs.sh
├── go.mod                       # Dependências Go
├── go.sum                       # Checksums
├── env.example                  # Exemplo de variáveis
└── README.md                    # Este arquivo
```

## 🚀 Instalação

### Pré-requisitos
- Go 1.21 ou superior
- Git
- Make (opcional)

### 1. Clone o repositório
```bash
git clone https://github.com/rafael-bit/whatz-chat.git
cd whatz-chat/Backend
```

### 2. Instale as dependências
```bash
go mod download
```

### 3. Configure as variáveis de ambiente
```bash
cp env.example .env
# Edite o arquivo .env com suas configurações
```

### 4. Execute o projeto
```bash
# Desenvolvimento (com hot reload)
make dev

# Ou diretamente
go run cmd/server/main.go
```

## ⚙️ Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` baseado no `env.example`:

```env
# Servidor
PORT=8080
HOST=localhost

# Banco de Dados
DB_PATH=./whatz.db

# JWT
JWT_SECRET=your-secret-key-here

# CORS
ALLOWED_ORIGINS=http://localhost:3001,http://localhost:3000

# Logs
LOG_LEVEL=info
LOG_FILE=./logs/app.log

# WebSocket
WS_PATH=/ws
```

### Banco de Dados

O sistema usa SQLite por padrão. O banco será criado automaticamente na primeira execução.

```bash
# Limpar banco de dados
make clear-db

# Executar seeds
make seed
```

## 🎯 Uso

### Desenvolvimento

```bash
# Iniciar servidor de desenvolvimento
make dev

# Ou com air (hot reload)
air

# Executar testes
make test

# Gerar documentação
make docs
```

### Produção

```bash
# Build da aplicação
make build

# Executar
./whatz-chat
```

### Comandos Make Disponíveis

```bash
make help          # Mostra todos os comandos
make dev           # Desenvolvimento com hot reload
make build         # Build da aplicação
make run           # Executar aplicação
make test          # Executar testes
make docs          # Gerar documentação Swagger
make clear-db      # Limpar banco de dados
make seed          # Executar seeds
make clean         # Limpar arquivos temporários
```

## 📡 API Endpoints

### 👥 Usuários

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/users` | Listar todos os usuários |
| `POST` | `/api/v1/users` | Criar novo usuário |
| `GET` | `/api/v1/users/{id}` | Buscar usuário por ID |
| `PUT` | `/api/v1/users/{id}` | Atualizar usuário |
| `DELETE` | `/api/v1/users/{id}` | Deletar usuário |

### 🏠 Salas

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/rooms` | Listar salas |
| `POST` | `/api/v1/rooms` | Criar nova sala |
| `GET` | `/api/v1/rooms/{id}` | Buscar sala por ID |
| `PUT` | `/api/v1/rooms/{id}` | Atualizar sala |
| `DELETE` | `/api/v1/rooms/{id}` | Deletar sala |
| `GET` | `/api/v1/rooms/{id}/messages` | Mensagens da sala |

### 🔐 Administração

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/api/v1/admin/rooms` | Criar sala com controle de acesso |
| `GET` | `/api/v1/admin/users/role/{role}` | Usuários por role |
| `PUT` | `/api/v1/admin/users/{id}/role` | Atualizar role do usuário |
| `PUT` | `/api/v1/admin/users/{id}/tags` | Atualizar tags do usuário |

### 📝 Mensagens

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/messages` | Listar mensagens |
| `POST` | `/api/v1/messages` | Criar nova mensagem |
| `GET` | `/api/v1/messages/{id}` | Buscar mensagem por ID |
| `DELETE` | `/api/v1/messages/{id}` | Deletar mensagem |

## 🌐 WebSocket

### Conexão
```
WS ws://localhost:8080/ws?user_id={user_id}&room_id={room_id}
```

### Eventos Suportados

#### Cliente → Servidor
- `send_message`: Enviar mensagem
- `typing_start`: Iniciar digitação
- `typing_stop`: Parar digitação

#### Servidor → Cliente
- `new_message`: Nova mensagem recebida
- `typing_indicator`: Indicador de digitação
- `user_joined`: Usuário entrou na sala
- `user_left`: Usuário saiu da sala
- `message_history`: Histórico de mensagens
- `welcome`: Mensagem de boas-vindas

### Exemplo de Uso

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?user_id=123&room_id=456');

ws.onopen = () => {
  console.log('Conectado ao WebSocket');
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Mensagem recebida:', data);
};

// Enviar mensagem
ws.send(JSON.stringify({
  type: 'send_message',
  payload: { content: 'Olá, mundo!' }
}));
```

## 📚 Documentação

### Swagger UI
Acesse a documentação interativa em:
```
http://localhost:8080/swagger/index.html
```

### Gerar Documentação
```bash
# Gerar documentação Swagger
make docs

# Ou manualmente
swag init -g cmd/server/main.go
```

### Exemplos de Uso

#### Criar Usuário
```bash
curl -X POST "http://localhost:8080/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "joao123",
    "email": "joao@example.com",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=joao123"
  }'
```

#### Criar Sala
```bash
curl -X POST "http://localhost:8080/api/v1/rooms" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sala Geral",
    "description": "Sala para conversas gerais",
    "type": "public",
    "created_by": "user-id"
  }'
```

## 🔧 Desenvolvimento

### Estrutura de Código

O projeto segue a arquitetura em camadas:

1. **Controllers**: Recebem requisições HTTP
2. **Services**: Contêm a lógica de negócio
3. **Repository**: Acessam o banco de dados
4. **Models**: Definem as estruturas de dados

### Padrões Utilizados

- **Dependency Injection**: Injeção de dependências
- **Repository Pattern**: Abstração do banco de dados
- **Service Layer**: Separação de responsabilidades
- **Middleware Pattern**: Interceptadores de requisições
- **WebSocket Hub**: Gerenciamento de conexões

### Convenções

- **Nomenclatura**: camelCase para variáveis, PascalCase para tipos
- **Estrutura**: Um arquivo por funcionalidade
- **Documentação**: Comentários em português
- **Logs**: Estruturados e contextualizados

## 🧪 Testes

### Executar Testes
```bash
# Todos os testes
make test

# Testes específicos
go test ./internal/controllers/...
go test ./internal/services/...

# Com cobertura
go test -cover ./...

# Testes de integração
go test -tags=integration ./...
```

### Estrutura de Testes
```
internal/
├── controllers/
│   ├── user_controller.go
│   └── user_controller_test.go
├── services/
│   ├── user_service.go
│   └── user_service_test.go
└── repository/
    ├── user_repository.go
    └── user_repository_test.go
```

## 🚀 Deploy

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
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

### Build e Execução
```bash
# Build da imagem
docker build -t whatz-chat-backend .

# Executar container
docker run -p 8080:8080 whatz-chat-backend
```

### Variáveis de Produção
```env
PORT=8080
DB_PATH=/data/whatz.db
JWT_SECRET=your-production-secret
ALLOWED_ORIGINS=https://yourdomain.com
LOG_LEVEL=warn
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👨‍💻 Autor

**Rafael Bit**
- GitHub: [@rafael-bit](https://github.com/rafael-bit)
- Email: rafael@example.com

## 🙏 Agradecimentos

- [Gin Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Swaggo](https://github.com/swaggo/swag)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)

---

**Whatz Chat** - Sistema de Chat em Tempo Real 🚀
