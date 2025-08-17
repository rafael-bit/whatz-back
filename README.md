# Whatz Chat - Sistema de Chat em Tempo Real

Sistema completo de chat em tempo real com backend em Go e frontend em Next.js 15, incluindo WebSocket, interface moderna e funcionalidades avançadas.

## 🚀 Funcionalidades

### Backend (Go)
- **API REST**: Endpoints para usuários, salas e mensagens
- **WebSocket**: Comunicação em tempo real
- **Banco SQLite**: Persistência de dados
- **Logging**: Sistema de logs informativos
- **Tracing**: Medição de performance
- **Swagger**: Documentação automática da API
- **CORS**: Configuração para desenvolvimento

### Frontend (Next.js 15)
- **Interface Moderna**: Design responsivo com shadcn/ui
- **Tema Claro/Escuro**: Suporte a múltiplos temas
- **Indicadores de Digitação**: Tempo real
- **Seletor de Emojis**: Biblioteca completa
- **Scroll Infinito**: Histórico de mensagens
- **Reconexão Automática**: WebSocket resiliente
- **Markdown Básico**: Formatação de texto
- **Agrupamento de Mensagens**: UX otimizada

## 🛠️ Tecnologias

### Backend
- **Go**: Linguagem principal
- **Fiber**: Framework web
- **SQLite**: Banco de dados
- **WebSocket**: Comunicação real-time
- **Swagger**: Documentação API

### Frontend
- **Next.js 15**: Framework React
- **TypeScript**: Tipagem estática
- **Tailwind CSS**: Estilização
- **shadcn/ui**: Componentes
- **Radix UI**: Primitivos acessíveis
- **Emoji Picker**: Seletor de emojis

## 📦 Instalação e Execução

### Pré-requisitos
- Go 1.21+
- Node.js 18+
- npm ou yarn

### 1. Backend

```bash
cd Back-end

# Instalar dependências
go mod download

# Configurar variáveis de ambiente (opcional)
cp env.example .env

# Executar o servidor
go run cmd/server/main.go
```

O backend estará disponível em:
- API: http://localhost:8080
- WebSocket: ws://localhost:8080/ws
- Swagger: http://localhost:8080/swagger/

### 2. Frontend

```bash
cd frontend

# Instalar dependências
npm install

# Configurar variáveis de ambiente
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" > .env.local
echo "NEXT_PUBLIC_WS_URL=ws://localhost:8080" >> .env.local

# Executar em desenvolvimento
npm run dev
```

O frontend estará disponível em: http://localhost:3000

## 🎯 Como Usar

1. **Acesse o frontend**: http://localhost:3000
2. **Crie um usuário** ou selecione um existente
3. **Crie uma sala** ou selecione uma existente
4. **Conecte-se ao chat** e comece a conversar!

### Funcionalidades do Chat

- **Envio de mensagens**: Digite e pressione Enter
- **Emojis**: Clique no ícone de emoji para selecionar
- **Markdown**: Use `**negrito**`, `*itálico*`, `` `código` ``
- **Indicadores**: Veja quando outros estão digitando
- **Tema**: Alterne entre claro/escuro no cabeçalho

## 📁 Estrutura do Projeto

```
Vaga/
├── Back-end/                 # Servidor Go
│   ├── cmd/server/          # Ponto de entrada
│   ├── internal/            # Código interno
│   │   ├── controllers/     # Controladores HTTP
│   │   ├── models/          # Modelos de dados
│   │   ├── services/        # Lógica de negócio
│   │   ├── repository/      # Acesso a dados
│   │   ├── websocket/       # WebSocket handlers
│   │   └── database/        # Configuração do banco
│   ├── docs/               # Documentação Swagger
│   └── scripts/            # Scripts utilitários
└── frontend/               # Cliente Next.js
    ├── src/
    │   ├── app/            # App Router
    │   ├── components/     # Componentes React
    │   ├── contexts/       # Contextos React
    │   ├── services/       # Serviços de API
    │   └── types/          # Tipos TypeScript
    └── public/             # Arquivos estáticos
```

## 🔌 API Endpoints

### Usuários
- `GET /api/v1/users` - Listar usuários
- `POST /api/v1/users` - Criar usuário
- `GET /api/v1/users/:id` - Buscar usuário
- `PUT /api/v1/users/:id` - Atualizar usuário
- `DELETE /api/v1/users/:id` - Deletar usuário

### Salas
- `GET /api/v1/rooms` - Listar salas
- `GET /api/v1/rooms/public` - Salas públicas
- `POST /api/v1/rooms` - Criar sala
- `GET /api/v1/rooms/:id` - Buscar sala
- `GET /api/v1/rooms/:id/messages` - Mensagens da sala

### WebSocket
- `ws://localhost:8080/ws?user_id=X&room_id=Y` - Conectar ao chat

## 🎨 Componentes Frontend

### Principais
- **ChatSetup**: Configuração inicial
- **ChatInterface**: Interface principal
- **ChatHeader**: Cabeçalho com informações
- **MessageList**: Lista de mensagens
- **MessageItem**: Item de mensagem
- **MessageInput**: Input de envio

### Funcionalidades
- Agrupamento automático de mensagens
- Scroll inteligente
- Indicadores de digitação
- Seletor de emojis
- Tema claro/escuro
- Responsividade completa

## 🔧 Configuração

### Variáveis de Ambiente

#### Backend (.env)
```env
PORT=8080
DB_PATH=./whatz.db
LOG_LEVEL=info
CORS_ORIGIN=*
```

#### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

## 🚀 Deploy

### Backend
- **Heroku**: Deploy direto
- **Railway**: Deploy automático
- **DigitalOcean**: Droplet personalizado
- **Vercel**: Serverless functions

### Frontend
- **Vercel**: Deploy automático (recomendado)
- **Netlify**: Deploy com drag & drop
- **Railway**: Deploy full-stack
- **AWS Amplify**: Deploy na AWS

## 📊 Logs e Monitoramento

### Backend
- Logs estruturados com timestamps
- Tracing de performance
- Métricas de WebSocket
- Health checks

### Frontend
- Logs de console para debugging
- Métricas de conexão WebSocket
- Tratamento de erros
- Estados de loading

## 🔒 Segurança

- Validação de entrada
- Sanitização de dados
- CORS configurado
- Tratamento de erros
- Logs de auditoria

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🆘 Suporte

### Problemas Comuns

1. **Backend não inicia**
   - Verifique se Go está instalado
   - Execute `go mod download`
   - Verifique a porta 8080

2. **Frontend não conecta**
   - Verifique se o backend está rodando
   - Confirme as variáveis de ambiente
   - Verifique o console do navegador

3. **WebSocket não funciona**
   - Verifique se o backend está rodando
   - Confirme a URL do WebSocket
   - Verifique os logs do backend

### Logs Úteis

- **Backend**: Logs no terminal
- **Frontend**: Console do navegador (F12)
- **WebSocket**: Network tab do DevTools

## 🔄 Roadmap

### Próximas Funcionalidades
- [ ] Autenticação JWT
- [ ] Upload de arquivos
- [ ] Notificações push
- [ ] Salas privadas
- [ ] Moderação de conteúdo
- [ ] Histórico de mensagens
- [ ] Exportação de dados
- [ ] Analytics e métricas

### Melhorias Técnicas
- [ ] Testes automatizados
- [ ] CI/CD pipeline
- [ ] Docker compose
- [ ] Monitoramento avançado
- [ ] Cache Redis
- [ ] Load balancing

---

**Desenvolvido com ❤️ usando Go e Next.js**

