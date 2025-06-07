// internal/domain/game/factory.go
package game

func NewGame(gameType, id GameID) Game {
	switch gameType {
	case "culo":
		return NewCuloCardGame(id)
	// case "truco": return NewTrucoGame(id)
	default:
		return nil
	}
}
