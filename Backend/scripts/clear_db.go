package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	log.Printf("🧹 Limpando banco de dados...")

	dbPath := "./whatz.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	// Conectar ao banco
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("❌ Erro ao abrir banco: %v", err)
	}
	defer db.Close()

	// Limpar todas as tabelas
	tables := []string{"messages", "rooms", "users"}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			log.Printf("⚠️ Erro ao limpar tabela %s: %v", table, err)
		} else {
			log.Printf("✅ Tabela %s limpa", table)
		}
	}

	log.Printf("🎉 Banco de dados limpo com sucesso!")
	log.Printf("💡 Agora você pode criar usuários e salas via API")
}
