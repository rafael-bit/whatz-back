package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	start := time.Now()
	log.Printf("🔧 Iniciando conexão com banco de dados: %s", dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco de dados: %v", err)
	}

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar com banco de dados: %v", err)
	}

	// Configurar conexão
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	database := &Database{DB: db}

	// Executar migrações
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("erro ao executar migrações: %v", err)
	}

	log.Printf("✅ Banco de dados conectado com sucesso em %v", time.Since(start))
	return database, nil
}

func (d *Database) migrate() error {
	start := time.Now()
	log.Printf("🔄 Executando migrações do banco de dados...")

	// Criar tabela de usuários
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		avatar TEXT,
		status TEXT DEFAULT 'online',
		role TEXT DEFAULT 'user',
		tags TEXT DEFAULT '[]',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	// Criar tabela de salas
	createRoomsTable := `
	CREATE TABLE IF NOT EXISTS rooms (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		type TEXT DEFAULT 'public',
		access_tags TEXT DEFAULT '[]',
		created_by TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (created_by) REFERENCES users (id)
	);`

	// Criar tabela de mensagens
	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		content TEXT NOT NULL,
		user_id TEXT NOT NULL,
		username TEXT NOT NULL,
		avatar TEXT,
		type TEXT DEFAULT 'text',
		room_id TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (id),
		FOREIGN KEY (room_id) REFERENCES rooms (id)
	);`

	// Criar tabela de tags
	createTagsTable := `
	CREATE TABLE IF NOT EXISTS tags (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	// Criar índices para performance
	createIndexes := `
	CREATE INDEX IF NOT EXISTS idx_messages_room_id ON messages (room_id);
	CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages (created_at);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
	CREATE INDEX IF NOT EXISTS idx_rooms_type ON rooms (type);
	`

	queries := []string{
		createUsersTable,
		createRoomsTable,
		createMessagesTable,
		createTagsTable,
		createIndexes,
	}

	for i, query := range queries {
		if _, err := d.DB.Exec(query); err != nil {
			return fmt.Errorf("erro na migração %d: %v", i+1, err)
		}
	}

	log.Printf("✅ Migrações executadas com sucesso em %v", time.Since(start))
	return nil
}

func (d *Database) Close() error {
	log.Printf("🔌 Fechando conexão com banco de dados...")
	return d.DB.Close()
}
