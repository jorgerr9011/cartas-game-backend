package game

type GameFactory struct{}

func NewGameFactory() *GameFactory {
	return &GameFactory{}
}

func (gamefactory *GameFactory) NewGame(gameType string) Game {
	switch gameType {
	case "culo":
		return NewCuloCardGame()
	// case "truco": return NewTrucoGame()
	default:
		return NewCuloCardGame()
	}
}
