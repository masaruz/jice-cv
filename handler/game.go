package handler

import (
	"999k_engine/constant"
	"999k_engine/engine"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"time"
)

// SetGambit current game to the gamestate
func SetGambit(game engine.Gambit) {
	state.GS.Gambit = game
}

// CreateSeats prepare empty seat for players
func CreateSeats(seats int) {
	for i := 0; i < seats; i++ {
		state.GS.Players = append(state.GS.Players, model.Player{Slot: i})
	}
}

// StartTable make table start
func StartTable() {
	state.GS.IsTableStart = true
}

// IsTableStart true or false
func IsTableStart() bool {
	return state.GS.IsTableStart
}

// IsGameStart true or false
func IsGameStart() bool {
	return state.GS.IsGameStart
}

// IsEndRound if current time is more than finish time
func IsEndRound() bool {
	return state.GS.FinishRoundTime.Unix() <= time.Now().Unix()
}

// IsPlayerTurn if player do something before deadline
func IsPlayerTurn(id string) bool {
	index, _ := util.Get(state.GS.Players, id)
	nowline := time.Now().Unix()
	startline := state.GS.Players[index].StartLine.Unix()
	deadline := state.GS.Players[index].DeadLine.Unix()
	return nowline >= startline && nowline < deadline
}

// FindWinner find a winner by evaluate his cards
func FindWinner() {
	hscore := -1
	hbonus := -1
	pos := -1
	// hkind := ""
	// winner := model.Player{}
	players := state.GS.Players
	for index, player := range players {
		if !util.InGame(player) {
			continue
		}
		scores, _ := state.GS.Gambit.Evaluate(player.Cards)
		score := scores[0]
		bonus := scores[1]
		if hscore < score {
			hscore = score
			hbonus = bonus
			// winner = player
			// hkind = kind
			pos = index
		} else if hscore == score && hbonus < bonus {
			hscore = score
			hbonus = bonus
			// winner = player
			// hkind = kind
			pos = index
		}
	}
	state.GS.Players[pos].IsWinner = true
}

// CreateTimeLine set timeline for game and any players
func CreateTimeLine(decisionTime int) {
	loop, delay := 0, 0
	start, amount := time.Now(), len(state.GS.Players)
	state.GS.StartRoundTime = start
	dealer, _ := util.FindDealer(state.GS.Players)
	for loop < amount {
		next := (dealer + 1) % amount
		player := state.GS.Players[next]
		if util.InGame(player) {
			startline := start.Add(time.Second * time.Duration(delay))
			state.GS.Players[next].Action = model.Action{}
			state.GS.Players[next].StartLine = startline
			state.GS.Players[next].DeadLine = startline.Add(time.Second * time.Duration(decisionTime))
			start = state.GS.Players[next].DeadLine
		}
		dealer++
		loop++
	}
	SetOtherDefaultAction("", model.Action{Name: constant.Check})
	state.GS.FinishRoundTime = start
}

// IncreaseTurn to seperate player bets
func IncreaseTurn() {
	state.GS.Turn++
}

// IsFullHand check if hold max cards
func IsFullHand(maxcards int) bool {
	for _, player := range state.GS.Players {
		if !util.InGame(player) {
			continue
		}
		// amount of cards are not equal
		if len(player.Cards) != 3 {
			return false
		}
	}
	return true
}

// BetsEqual check if same bet
func BetsEqual() bool {
	baseBet := -1
	for _, player := range state.GS.Players {
		if !util.InGame(player) {
			continue
		}
		// check if everyone has the same bet
		bet := util.SumBet(player)
		if baseBet != bet && baseBet == -1 {
			baseBet = bet
		} else if baseBet != bet {
			return false
		}
	}
	return true
}

// Deal card to the players
func Deal(cardAmount int, playerAmount int) {
	dealer, _ := util.FindDealer(state.GS.Players)
	index := -1
	// deal card start from next to dealer
	for i := 0; i < cardAmount; i++ {
		start := dealer
		round := 0
		for round < playerAmount {
			start++
			round++
			index = start % playerAmount
			// skip empty seat
			if util.InGame(state.GS.Players[index]) {
				state.GS.Players[index].Cards = append(state.GS.Players[index].Cards, Draw())
				if index == dealer {
					break
				}
			}
		}
	}
}

// FinishGame prepare to finish game
func FinishGame() {
	state.GS.Turn = 0
	state.GS.IsGameStart = false
}

// ShiftTimeline shift timeline of everyone because someone take action
func ShiftTimeline(diff time.Duration) {
	for index := range state.GS.Players {
		state.GS.Players[index].StartLine = state.GS.Players[index].StartLine.Add(diff)
		state.GS.Players[index].DeadLine = state.GS.Players[index].DeadLine.Add(diff)
	}
	state.GS.FinishRoundTime = state.GS.FinishRoundTime.Add(diff)
}

// ShiftPlayerToEndOfTimeline shift player to the end of timeline
func ShiftPlayerToEndOfTimeline(id string, second int) {
	duration := time.Duration(time.Second * time.Duration(second))
	index, _ := util.Get(state.GS.Players, id)
	finishRoundTime := state.GS.FinishRoundTime
	state.GS.Players[index].StartLine = finishRoundTime
	state.GS.Players[index].DeadLine = finishRoundTime.Add(duration)
	state.GS.FinishRoundTime = finishRoundTime.Add(duration)
}

// ShiftPlayersToEndOfTimeline shift current and prev player to the end of timeline
func ShiftPlayersToEndOfTimeline(id string, second int) {
	start, _ := util.Get(state.GS.Players, id)
	round, amount := 0, len(state.GS.Players)
	for round < amount {
		start++
		round++
		index := start % amount
		// force shift only 2 lastest players
		if util.InGame(state.GS.Players[index]) &&
			util.IsPlayerBehindTheTimeline(state.GS.Players[index]) {
			ShiftPlayerToEndOfTimeline(state.GS.Players[index].ID, second)
		}
	}
}

// GetCurrentTurn get current turn number
func GetCurrentTurn() int {
	return state.GS.Turn
}

// IncreasePots when increase pots values
func IncreasePots(chips int, index int) {
	if len(state.GS.Pots) <= 0 {
		state.GS.Pots = []int{0}
	}
	// increase pot values
	state.GS.Pots[0] += chips
}

// InvestToPots added bet to everyone base on turn
func InvestToPots(chips int) {
	// initiate bet value to players
	for index := range state.GS.Players {
		if util.InGame(state.GS.Players[index]) {
			state.GS.Players[index].Bets = append(state.GS.Players[index].Bets, chips)
			IncreasePots(chips, GetCurrentTurn()) // start with first element in pots
		}
	}
}
