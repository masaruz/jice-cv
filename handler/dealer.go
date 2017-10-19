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
	state.GS.Deck.Cards = cards
}

// Shuffle cards in deck
func Shuffle() {
	cards := state.GS.Deck.Cards
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := len(cards)
	for i := 0; i < n; i++ {
		j := r.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
	state.GS.Deck.Cards = cards
}

// Draw a card by poping the deck
func Draw() int {
	c := state.GS.Deck.Cards
	state.GS.Deck.Cards = c[:len(c)-1]
	return c[len(c)-1]
}

// SetDealer select who is dealer
func SetDealer() {
	next := -1
	amount := len(state.GS.Players)
	for index, player := range state.GS.Players {
		if player.Type == constant.Dealer && next == -1 {
			state.GS.Players[index].Type = constant.Normal
			// find next player to be a dealer
			for next <= amount*2 {
				next = (index + 1) % amount
				if state.GS.Players[next].IsPlaying {
					state.GS.Players[next].Type = constant.Dealer
					break
				}
				index++
			}
		} else if next != index && player.IsPlaying {
			state.GS.Players[index].Type = constant.Normal
		}
	}
	// set first player as a dealer
	if next == -1 {
		for index, player := range state.GS.Players {
			// skip empty seat
			if player.IsPlaying {
				state.GS.Players[index].Type = constant.Dealer
				break
			}
		}
	}
}
