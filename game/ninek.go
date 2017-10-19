package game

import (
	"999k_engine/constant"
	"999k_engine/handler"
	"999k_engine/util"
	"sort"
)

// NineK is 9K
type NineK struct {
	MaxPlayers   int
	DecisionTime int64
	MinimumBet   int
}

// Payload data accessed by continue
type Payload struct {
	ID    string
	Chips int
}

// Init deck and environment variables
func (game NineK) Init() {
	// set the seats
	handler.CreateSeats(game.MaxPlayers)
}

// Start game
func (game NineK) Start() bool {
	if handler.IsTableStart() && !handler.IsGameStart() && handler.MakePlayersReady() {
		handler.StartGame()
		handler.InvestToPots(game.MinimumBet)
		handler.SetDealer()
		handler.BuildDeck()
		handler.Shuffle()
		handler.CreateTimeLine(game.DecisionTime)
		handler.Deal(2, game.MaxPlayers)
		return true
	}
	return false
}

// NextRound game after round by round
func (game NineK) NextRound() bool {
	if !handler.IsFullHand(3) && handler.BetsEqual() && handler.IsEndRound() &&
		util.CountPlaying(handler.GetPlayerState()) > 1 {
		handler.Deal(1, game.MaxPlayers)
		handler.CreateTimeLine(game.DecisionTime)
		handler.InvestToPots(0)
		handler.IncreaseTurn()
		return true
	}
	handler.OverwriteActionToBehindPlayers()
	return false
}

// Check is doing nothing only shift the timeline
func (game NineK) Check(id string) bool {
	return handler.Check(id)
}

// Bet is raising bet to the target
func (game NineK) Bet(id string, chips int) bool {
	return handler.Bet(id, chips, game.DecisionTime)
}

// Call is raising bet to the highest bet
func (game NineK) Call(id string) bool {
	return handler.Call(id, game.DecisionTime)
}

// AllIn give all chips
func (game NineK) AllIn(id string) bool {
	return false
}

// Fold quit the game but still lost bet
func (game NineK) Fold(id string) bool {
	return handler.Fold(id)
}

// Finish game
func (game NineK) Finish() bool {
	// no others to play with or all players have 3 cards but bet is not equal
	if util.CountPlaying(handler.GetPlayerState()) <= 1 ||
		(handler.IsFullHand(3) && handler.BetsEqual() && handler.IsEndRound()) {
		handler.ForceEndTimeline()
		handler.AssignWinner()
		handler.FlushGame()
		return true
	}
	return false
}

// End game
func (game NineK) End() {}

// Evaluate score of player cards
func (game NineK) Evaluate(values []int) (scores []int, kind string) {
	sorted := make([]int, len(values))
	copy(sorted, values)
	sort.Ints(sorted)
	// all cards are same number
	if game.threeOfAKind(sorted) {
		return game.summary(constant.ThreeOfAKind, sorted)
	}
	// all cards are in seqence and same suit
	if game.straightFlush(sorted) {
		return game.summary(constant.StraightFlush, sorted)
	}
	// all cards are J or Q or K
	if game.royal(sorted) {
		return game.summary(constant.Royal, sorted)
	}
	// all cards are in sequnce
	if game.straight(sorted) {
		return game.summary(constant.Straight, sorted)
	}
	return game.summary(constant.Nothing, sorted)
}

// Summary score of each kind
func (game NineK) summary(kind string, hands []int) ([]int, string) {
	bonus := hands[len(hands)-1]
	switch kind {
	case constant.ThreeOfAKind:
		// becase 333 is the highest
		if util.GetCardNumberFromValue(bonus) == 3 {
			bonus = 52 //max
		}
		return []int{10000000, bonus}, constant.ThreeOfAKind
	case constant.StraightFlush:
		return []int{1000000, bonus}, constant.StraightFlush
	case constant.Royal:
		return []int{100000, bonus}, constant.Royal
	case constant.Straight:
		return []int{10000, bonus}, constant.Straight
	case constant.Flush:
		return []int{1000, bonus}, constant.Straight
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
		return []int{score % 10, bonus}, constant.Nothing
	}
}

// ThreeOfAKind when three cards are same number
func (game NineK) threeOfAKind(values []int) bool {
	number := util.GetCardNumberFromValue(values[0])
	for _, value := range values {
		if number != util.GetCardNumberFromValue(value) {
			return false
		}
	}
	return true
}

// StraightFlush when three cards are same suit and order in sequence
func (game NineK) straightFlush(values []int) bool {
	return game.flush(values) && game.straight(values)
}

// Straight when three cards order in sequence
func (game NineK) straight(values []int) bool {
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
func (game NineK) royal(values []int) bool {
	for _, value := range values {
		number := util.GetCardNumberFromValue(value)
		if number <= 10 || number == 14 {
			return false
		}
	}
	return true
}

// Flush when 3 cards have same suit
func (game NineK) flush(values []int) bool {
	suit := util.GetCardSuitFromValue(values[0])
	for _, value := range values {
		// check suit
		if suit != util.GetCardSuitFromValue(value) {
			return false
		}
	}
	return true
}
