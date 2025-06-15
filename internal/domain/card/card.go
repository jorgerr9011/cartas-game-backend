package card

type Suit string

const (
	Oros    Suit = "oros"
	Copas   Suit = "copas"
	Espadas Suit = "espadas"
	Bastos  Suit = "bastos"
)

type Rank int

const (
	As      Rank = 1
	Dos     Rank = 2
	Tres    Rank = 3
	Cuatro  Rank = 4
	Cinco   Rank = 5
	Seis    Rank = 6
	Siete   Rank = 7
	Sota    Rank = 10
	Caballo Rank = 11
	Rey     Rank = 12
)

type Card struct {
	Suit Suit
	Rank Rank
}

func NewSpanishDeck() []Card {
	suits := []Suit{Oros, Copas, Espadas, Bastos}
	ranks := []Rank{1, 2, 3, 4, 5, 6, 7, 10, 11, 12}

	var deck []Card
	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{Suit: suit, Rank: rank})
		}
	}
	return deck
}
