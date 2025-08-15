#!/bin/bash

# Script para gerar documentação Swagger
echo "🔧 Gerando documentação Swagger..."

# Verificar se o swag está instalado
if ! command -v swag &> /dev/null; then
    echo "📦 Instalando swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Gerar documentação
echo "📝 Gerando docs..."
swag init -g cmd/server/main.go -o docs

# Verificar se a geração foi bem-sucedida
if [ $? -eq 0 ]; then
    echo "✅ Documentação gerada com sucesso!"
    echo "📚 Acesse: http://localhost:8080/swagger/"
    echo "📁 Arquivos gerados em: docs/"
else
    echo "❌ Erro ao gerar documentação"
    exit 1
fi
