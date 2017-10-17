package util

import (
	"999k_websocket/constant"
	"999k_websocket/model"
	"fmt"
	"math"
)

// GetCardNumberFromValue convert value to card readable number
func GetCardNumberFromValue(value int) int {
	return int(math.Floor(float64(value)/4)) + 2
}

// GetCardSuitFromValue convert value to card suit
func GetCardSuitFromValue(value int) int {
	return value % 4
}

// PrintCard convert card into human readable
func PrintCard(value int) {
	fmt.Printf("%s:%s ",
		model.Number[GetCardNumberFromValue(value)-2],
		model.Suit[GetCardSuitFromValue(value)])
}

// FindDealer return index of dealer
func FindDealer(players model.Players) (int, model.Player) {
	// find dealer
	for index, player := range players {
		if player.Type == constant.Dealer {
			return index, player
		}
	}
	return 0, model.Player{Type: constant.Dealer}
}

// SumBet to sum each player chips bet
func SumBet(player model.Player) int {
	sum := 0
	for _, bet := range player.Bets {
		sum += bet
	}
	return sum
}

// SumBets to sum every player chips bet
func SumBets(turn int, players model.Players) int {
	sum := 0
	for _, player := range players {
		if !player.IsPlaying {
			continue
		}
		sum += player.Bets[turn]
	}
	return sum
}

// SumPots to sum every pot
func SumPots(pots []int) int {
	sum := 0
	for _, pot := range pots {
		sum += pot
	}
	return sum
}
