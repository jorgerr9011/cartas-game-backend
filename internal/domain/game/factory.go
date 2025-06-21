package game

type GameFactory struct{}

func NewGameFactory() *GameFactory {
	return &GameFactory{}
}

func (gamefactory *GameFactory) NewGame(gameType string) Game {
	switch gameType {
	case "culo":
		return NewCuloCardGame()
	default:
		return NewCuloCardGame()
	}
}
