package gambit

import (
	"999k_engine/constant"
	"999k_engine/util"
)

// Summary score of each kind
func summary(kind string, hands []int) ([]int, string) {
	bonus := hands[len(hands)-1]
	switch kind {
	case constant.ThreeOfAKind:
		return []int{10000000, bonus, 0, 0, 0}, constant.ThreeOfAKind
	case constant.StraightFlush:
		return []int{1000000, bonus, 0, 0, 0}, constant.StraightFlush
	case constant.Royal:
		return []int{100000, bonus, 0, 0, 0}, constant.Royal
	case constant.Straight:
		return []int{10000, bonus, 0, 0, 0}, constant.Straight
	case constant.Flush:
		scores := []int{1000}
		for index := len(hands) - 1; index >= 0; index-- {
			if index < len(hands) {
				scores = append(scores, util.GetCardNumberFromValue(hands[index]))
			}
		}
		for index := len(scores); index < 4; index++ {
			scores = append(scores, 0)
		}
		scores = append(scores, bonus)
		return scores, constant.Flush
	default:
		score := 0
		for _, value := range hands {
			number := util.GetCardNumberFromValue(value)
			// if value is A
			if number == 14 {
				number = 1
			}
			if number == 11 || number == 12 || number == 13 {
				number = 10
			}
			score += number
		}
		return []int{score % 10, bonus, 0, 0, 0, 0}, constant.Nothing
	}
}

// ThreeOfAKind when three cards are same number
func threeOfAKind(values []int) bool {
	number := util.GetCardNumberFromValue(values[0])
	for _, value := range values {
		if number != util.GetCardNumberFromValue(value) {
			return false
		}
	}
	return true
}

// StraightFlush when three cards are same suit and order in sequence
func straightFlush(values []int) bool {
	return flush(values) && straight(values)
}

// Straight when three cards order in sequence
func straight(values []int) bool {
	number := util.GetCardNumberFromValue(values[0])
	for i := 1; i < len(values); i++ {
		current := util.GetCardNumberFromValue(values[i])
		if current-number != 1 {
			return false
		}
		number = current
	}
	// because 48 - 51 are A
	return number < 13 && util.GetCardNumberFromValue(values[1]) < 13
}

// Royal when 3 cards have no any number (only J,Q,K)
func royal(values []int) bool {
	for _, value := range values {
		number := util.GetCardNumberFromValue(value)
		if number <= 10 || number == 14 {
			return false
		}
	}
	return true
}

// Flush when 3 cards have same suit
func flush(values []int) bool {
	suit := util.GetCardSuitFromValue(values[0])
	for _, value := range values {
		// check suit
		if suit != util.GetCardSuitFromValue(value) {
			return false
		}
	}
	return true
}
