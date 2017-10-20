package util

import (
	"999k_engine/constant"
	"999k_engine/model"
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

// SumBetsInTurn to sum every player chips bet
func SumBetsInTurn(turn int, players model.Players) int {
	sum := 0
	for _, player := range players {
		if !player.IsPlaying {
			continue
		}
		sum += player.Bets[turn]
	}
	return sum
}

// GetHighestBet for get the highest bet and make a call action
func GetHighestBet(players model.Players) int {
	highest := 0
	for _, player := range players {
		if !InGame(player) {
			continue
		}
		sum := SumBet(player)
		if sum > highest {
			highest = sum
		}
	}
	return highest
}

// GetHighestBetInTurn for get the highest bet in specific turn
func GetHighestBetInTurn(turn int, players model.Players) int {
	highest := 0
	for _, player := range players {
		if !InGame(player) || len(player.Bets) <= turn {
			continue
		}
		bet := player.Bets[turn]
		if bet > highest {
			highest = bet
		}
	}
	return highest
}

// SumPots to sum every pot
func SumPots(pots []int) int {
	sum := 0
	for _, pot := range pots {
		sum += pot
	}
	return sum
}
