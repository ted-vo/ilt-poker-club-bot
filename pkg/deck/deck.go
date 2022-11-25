package deck

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SPADES   = "♠️ "
	HEARTS   = "♥️ "
	DIAMONDS = "♦️ "
	CLUBS    = "♣️ "
)

type Card struct {
	Suite string
	Value string
}

func (card *Card) ToString() string {
	return fmt.Sprintf("%s%s", card.Value, card.Suite)
}

type Deck struct {
	Cards []*Card
	lock  sync.RWMutex
}

func (deck *Deck) Shuffle() {
	t := time.Now()
	rand.Seed(int64(t.Nanosecond()))
	for i := range deck.Cards {
		j := rand.Intn(i + 1)
		if i != j {
			deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
		}
	}
}

func (deck *Deck) Push(card *Card) {
	deck.lock.Lock()
	defer deck.lock.Unlock()

	deck.Cards = append(deck.Cards, card)
}

func (deck *Deck) Pop() (*Card, error) {
	len := len(deck.Cards)
	var card *Card
	if deck.isNotEmpty() {
		deck.lock.Lock()
		defer deck.lock.Unlock()

		top := len - 1
		card = deck.Cards[top]
		deck.Cards = deck.Cards[:top]

		return card, nil
	}

	return card, fmt.Errorf("deck is empty")
}

func (deck *Deck) Size() int {
	return len(deck.Cards)
}

func (deck *Deck) isEmpty() bool {
	return len(deck.Cards) == 0
}

func (deck *Deck) isNotEmpty() bool {
	return len(deck.Cards) > 0
}

func NewDeck() *Deck {
	suites := []string{SPADES, HEARTS, DIAMONDS, CLUBS}
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	deck := &Deck{}
	for _, suite := range suites {
		for _, value := range values {
			deck.Cards = append(deck.Cards, &Card{Suite: suite, Value: value})
		}
	}

	return deck
}
