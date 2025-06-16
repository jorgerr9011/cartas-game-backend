# Muestra los comandos disponibles
default:
    @echo ""
    @echo "🃏 Comandos disponibles:"
    @echo ""
    @echo "  just test        Ejecuta los tests de WebSocket"

# Ejecuta los tests
test:
    go test ./internal/ports/ws -v
