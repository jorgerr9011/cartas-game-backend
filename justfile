# Muestra los comandos disponibles
default:
    @echo ""
    @echo "🃏 Comandos disponibles:"
    @echo ""
    @echo "  just install     Obtiene imágenes docker y despliega el proyecto"
    @echo "  just init        Despliega el proyecto"
    @echo "  just test        Ejecuta los tests de WebSocket"

# Instalar proyecto
install: 
    docker compose -f docker-compose.yml up --build -d

# Levantar contenedores
init:
    docker compose -f docker-compose.yml up -d

# Ejecuta los tests
test:
    docker exec -it cardgame-backend go test ./internal/ports/ws -v
