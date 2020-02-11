//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

//Suit for a standard Card
type Suit uint8

const (
	//Spade Suit
	Spade Suit = iota
	//Diamond Suit
	Diamond
	//Club Suit
	Club
	//Heart Suit
	Heart
	//Joker a special suit
	Joker
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

//Rank of standard Card
type Rank uint8

const (
	_ Rank = iota
	//Ace - 1st rank
	Ace
	//Two - 2nd rank
	Two
	//Three - 3rd rank
	Three
	//Four - 4th rank
	Four
	//Five - 5th rank
	Five
	//Six - 6th rank
	Six
	//Seven - 7th rank
	Seven
	//Eight - 8th rank
	Eight
	//Nine - 9th rank
	Nine
	//Ten - 10th rank
	Ten
	//Jack - 11th rank
	Jack
	//Queen - 12th rank
	Queen
	//King - 13th rank
	King
)

// Card defines a standard card
type Card struct {
	Suit
	Rank
}

const (
	minRank = Ace
	maxRank = King
)

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

// New retuns a new deck of cards.
// It can be customised by giving function parameters
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Rank: rank, Suit: suit})
		}
	}
	for _, opt := range opts {
		cards = opt(cards)
	}
	return cards
}

//DefaultSort sorts the card in default order, i.e., ascending order
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

//Sort sorts the cards based on function provided the user
func Sort(cmp func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, cmp(cards))
		return cards
	}
}

//Less returns a comparator function, i.e.,
//returns a function that checks whether index i less than j
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

//ShuffleDeterMinistic shuffles with the given random number generator the cards
func ShuffleDeterMinistic(r *rand.Rand) func([]Card) []Card {
	return func(cards []Card) []Card {
		ret := make([]Card, len(cards))
		perm := r.Perm(len(cards))
		for i, j := range perm {
			ret[i] = cards[j]
		}
		return ret
	}
}

//Shuffle the cards
func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	perm := r.Perm(len(cards))
	for i, j := range perm {
		ret[i] = cards[j]
	}
	return ret
}

//Jokers adds a number of jokers to the deck
func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{
				Rank: Rank(i),
				Suit: Joker,
			})
		}
		return cards
	}
}

//Filter outs the passed on card based on filter function provided by user
func Filter(f func(c Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for _, c := range cards {
			if !f(c) {
				ret = append(ret, c)
			}
		}
		return ret
	}
}

//Deck duplicates the current deck of cards n times
func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			ret = append(ret, cards...)
		}
		return ret
	}
}

//
