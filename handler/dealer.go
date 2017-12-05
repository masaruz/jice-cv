package handler

import (
	"999k_engine/constant"
	"999k_engine/model"
	"999k_engine/state"
	"math/rand"
	"time"
)

// BuildDeck a deck of cards
func BuildDeck() {
	cards := model.Cards{}
	for v := 0; v < 51; v++ {
		cards = append(cards, v)
	}
	state.Snapshot.Deck.Cards = cards
}

// Shuffle cards in deck
func Shuffle() {
	cards := state.Snapshot.Deck.Cards
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := len(cards)
	for i := 0; i < n; i++ {
		j := r.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
	state.Snapshot.Deck.Cards = cards
}

// Draw a card by poping the deck
func Draw() int {
	c := state.Snapshot.Deck.Cards
	state.Snapshot.Deck.Cards = c[:len(c)-1]
	return c[len(c)-1]
}

// SetDealer select who is dealer
func SetDealer() {
	next := -1
	amount := len(state.Snapshot.Players)
	for index, player := range state.Snapshot.Players {
		if player.Type == constant.Dealer && next == -1 {
			state.Snapshot.Players[index].Type = constant.Normal
			// find next player to be a dealer
			for next <= amount*2 {
				next = (index + 1) % amount
				if state.Snapshot.Players[next].IsPlaying {
					state.Snapshot.Players[next].Type = constant.Dealer
					break
				}
				index++
			}
		} else if next != index && player.IsPlaying {
			state.Snapshot.Players[index].Type = constant.Normal
		}
	}
	// set first player as a dealer
	if next == -1 {
		for index, player := range state.Snapshot.Players {
			// skip empty seat
			if player.IsPlaying {
				state.Snapshot.Players[index].Type = constant.Dealer
				break
			}
		}
	}
}
