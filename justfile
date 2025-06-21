# Muestra los comandos disponibles
default:
    @echo ""
    @echo "🃏 Comandos disponibles:"
    @echo ""
    @echo "  just test        Ejecuta los tests de WebSocket"

# Ejecuta los tests
test:
    docker exec -it cardgame-backend go test ./internal/ports/ws -v
