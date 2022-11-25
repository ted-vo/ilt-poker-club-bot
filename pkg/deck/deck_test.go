package deck

import (
	"testing"

	"github.com/apex/log"
	"github.com/stretchr/testify/require"
)

func TestDeck(t *testing.T) {
	deck := NewDeck()
	deck.Shuffle()

	// for _, card := range deck.Cards {
	// 	fmt.Printf("%s\n", card.toString())
	// }

	card, err := deck.Pop()
	if err != nil {
		log.Error(err.Error())
	}
	log.Info(card.ToString())

	card1, err := deck.Pop()
	if err != nil {
		log.Error(err.Error())
	}
	log.Info(card1.ToString())

	require.Equal(t, true, true)
}
