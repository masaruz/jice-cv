package handler

import (
	"999k_engine/constant"
	"999k_engine/engine"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"time"
)

// WaitQueue check if server is processing result
func WaitQueue() {
	for state.GS.IsProcessing {
	}
}

// StartProcess set IsProcessing to true to blocking
func StartProcess() {
	state.GS.IsProcessing = true
}

// FinishProcess set IsProcessing to false to unblocking
func FinishProcess() {
	state.GS.IsProcessing = false
}

// SetGambit current game to the gamestate
func SetGambit(game engine.Gambit) {
	state.GS.Gambit = game
}

// CreateSeats prepare empty seat for players
func CreateSeats(seats int) {
	// TODO
	names := []string{"Eleven", "Dustin", "Mike", "Lucas", "Nancy", "Will"}
	for i := 0; i < seats; i++ {
		state.GS.Players = append(state.GS.Players,
			model.Player{Slot: i, Name: names[i]})
	}
}

// CreatePots per seat
func CreatePots(length int) {
	for i := 0; i < length; i++ {
		state.GS.Pots = append(state.GS.Pots, 0)
	}
}

// StartTable set table start
func StartTable() {
	state.GS.IsTableStart = true
}

// StartGame set game start
func StartGame() {
	state.GS.IsGameStart = true
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
	return state.GS.FinishRoundTime <= time.Now().Unix()
}

// IsInExtendTime when end round but still not actually ended
func IsInExtendTime() bool {
	return state.GS.FinishRoundTime > time.Now().Unix()
}

// IncreaseTurn to seperate player bets
func IncreaseTurn() {
	state.GS.Turn++
}

// IncreasePots when players pay chips
func IncreasePots(index int, chips int) {
	state.GS.Pots[index] += chips
}

// SetMinimumBet set minimum players can bet
func SetMinimumBet(chips int) {
	state.GS.MinimumBet = chips
}

// SetMaximumBet set that maximum players can bet
func SetMaximumBet(chips int) {
	state.GS.MaximumBet = chips
}

// IsPlayerTurn if player do something before deadline
func IsPlayerTurn(id string) bool {
	index, _ := util.Get(state.GS.Players, id)
	if index == -1 {
		return false
	}
	nowline := time.Now().Unix()
	startline := state.GS.Players[index].StartLine
	deadline := state.GS.Players[index].DeadLine
	return nowline >= startline && nowline < deadline
}

// CreateTimeLine set timeline for game and any players
func CreateTimeLine(decisionTime int64) {
	loop := 0
	delay := int64(0)
	start, amount := time.Now().Unix(), len(state.GS.Players)
	state.GS.StartRoundTime = start
	// need at least one competitors
	if util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) <= 1 {
		state.GS.FinishRoundTime = start
	} else {
		dealer, _ := util.FindDealer(state.GS.Players)
		for loop < amount {
			next := (dealer + 1) % amount
			player := state.GS.Players[next]
			if util.IsPlayingAndNotFoldAndNotAllIn(player) {
				startline := start + delay
				state.GS.Players[next].StartLine = startline
				state.GS.Players[next].DeadLine = startline + decisionTime
				start = state.GS.Players[next].DeadLine
			}
			dealer++
			loop++
		}
		state.GS.FinishRoundTime = start
	}
}

// MakePlayersReady make everyone isPlayer = true
func MakePlayersReady() bool {
	for index, player := range state.GS.Players {
		if player.ID == "" {
			continue
		}
		// force to stand when player has no chips enough
		if player.Chips < state.GS.MinimumBet {
			Stand(player.ID)
			continue
		}
		state.GS.Players[index].Cards = model.Cards{}
		state.GS.Players[index].Bets = []int{}
		state.GS.Players[index].IsPlaying = true
		state.GS.Players[index].IsEarned = false
		state.GS.Players[index].IsWinner = false
		state.GS.Players[index].Default = model.Action{Name: constant.Check}
		state.GS.Players[index].Action = model.Action{}
	}
	return util.CountSitting(state.GS.Players) >= 2
}

// SetOtherDefaultAction make every has default action
func SetOtherDefaultAction(id string, action string) {
	daction := model.Action{Name: action}
	for index := range state.GS.Players {
		if !util.IsPlayingAndNotFoldAndNotAllIn(state.GS.Players[index]) {
			continue
		}
		// if chips equal 0 then must be allin
		if state.GS.Players[index].Chips == 0 {
			state.GS.Players[index].Default = model.Action{Name: constant.AllIn}
			state.GS.Players[index].Action = model.Action{Name: constant.AllIn}
			continue
		}
		if id != "" && id != state.GS.Players[index].ID {
			_, caller := util.Get(state.GS.Players, id)
			// if caller's bet more than others then overwrite their action
			if caller.Bets[state.GS.Turn] > state.GS.Players[index].Bets[state.GS.Turn] {
				state.GS.Players[index].Default = daction
				state.GS.Players[index].Action = model.Action{}
			}
		} else if id == "" {
			state.GS.Players[index].Default = daction
			state.GS.Players[index].Action = model.Action{}
		}
	}
}

// SetOtherActions make every has default action
func SetOtherActions(id string, action string) {
	for index, player := range state.GS.Players {
		if !util.IsPlayingAndNotFoldAndNotAllIn(player) {
			continue
		}
		if id != "" && id != player.ID {
			state.GS.Players[index].Actions = state.GS.Gambit.Reducer(action, player.ID)
		} else if id == "" {
			state.GS.Players[index].Actions = state.GS.Gambit.Reducer(action, player.ID)
		}
	}
}

// IsFullHand check if hold max cards
func IsFullHand(maxcards int) bool {
	for _, player := range state.GS.Players {
		if !util.IsPlayingAndNotFold(player) {
			continue
		}
		// amount of cards are not equal
		if len(player.Cards) != 3 {
			return false
		}
	}
	return util.CountPlayerNotFold(state.GS.Players) > 1
}

// BetsEqual check if same bet
func BetsEqual() bool {
	baseBet := util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players)
	for _, player := range state.GS.Players {
		if !util.IsPlayingAndNotFoldAndNotAllIn(player) {
			continue
		}
		if len(player.Bets)-1 < state.GS.Turn {
			return false
		}
		// check if everyone has the same bet
		if baseBet != player.Bets[state.GS.Turn] && player.Chips != 0 && baseBet != 0 {
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
			if util.IsPlayingAndNotFold(state.GS.Players[index]) {
				state.GS.Players[index].Cards = append(state.GS.Players[index].Cards, Draw())
				if index == dealer {
					break
				}
			}
		}
	}
}

// FlushGame reset everything before new game and client no needs to see it
func FlushGame() {
	for index := range state.GS.Players {
		state.GS.Players[index].IsPlaying = false
		state.GS.Players[index].Bets = []int{}
		state.GS.Players[index].Default = model.Action{}
		state.GS.Players[index].Action = model.Action{}
		state.GS.Players[index].Actions = model.Actions{}
		state.GS.Players[index].StartLine = 0
		state.GS.Players[index].DeadLine = 0
		state.GS.Players[index].IsEarned = false
	}
	for index := range state.GS.Pots {
		state.GS.Pots[index] = 0
	}
	state.GS.Turn = 0
	state.GS.IsGameStart = false
}

// ShortenTimeline shift timeline of everyone because someone take action
func ShortenTimeline(diff int64) {
	diff = util.Absolute(diff)
	for index, player := range state.GS.Players {
		if player.IsPlaying {
			state.GS.Players[index].StartLine = state.GS.Players[index].StartLine - diff
			state.GS.Players[index].DeadLine = state.GS.Players[index].DeadLine - diff
		}
	}
	state.GS.FinishRoundTime = state.GS.FinishRoundTime - diff
}

// ShortenTimelineAfterTarget shift timeline of everyone behind target player
func ShortenTimelineAfterTarget(id string, diff int64) {
	diff = util.Absolute(diff)
	_, caller := util.Get(state.GS.Players, id)
	for index, player := range state.GS.Players {
		// who start behind caller will be shifted
		if util.IsPlayingAndNotFoldAndNotAllIn(player) && player.StartLine >= caller.DeadLine {
			state.GS.Players[index].StartLine = state.GS.Players[index].StartLine - diff
			state.GS.Players[index].DeadLine = state.GS.Players[index].DeadLine - diff
		}
	}
	state.GS.FinishRoundTime = state.GS.FinishRoundTime - diff
}

// ExtendTime force this timeline to be ended
func ExtendTime() {
	// added delay for display the winner
	state.GS.FinishRoundTime = time.Now().Unix() + 5
}

// ShiftPlayersToEndOfTimeline shift current and prev player to the end of timeline
func ShiftPlayersToEndOfTimeline(id string, second int64) {
	start, _ := util.Get(state.GS.Players, id)
	round, amount := 0, len(state.GS.Players)
	for round < amount {
		start++
		round++
		next := start % amount
		// force shift to players who is in game not allin and behine the timeline
		if util.IsPlayingAndNotFoldAndNotAllIn(state.GS.Players[next]) &&
			util.IsPlayerBehindTheTimeline(state.GS.Players[next]) {
			finishRoundTime := state.GS.FinishRoundTime
			state.GS.Players[next].StartLine = finishRoundTime
			state.GS.Players[next].DeadLine = finishRoundTime + second
			state.GS.FinishRoundTime = finishRoundTime + second
		}
	}
}

// PlayersInvestToPots added bet to everyone base on turn
func PlayersInvestToPots(chips int) {
	// initiate bet value to players
	for index := range state.GS.Players {
		if util.IsPlayingAndNotFoldAndNotAllIn(state.GS.Players[index]) {
			state.GS.Players[index].Chips -= chips
			state.GS.Players[index].Bets = append(state.GS.Players[index].Bets, chips)
			IncreasePots(index, chips)
			// start with first element in pots
		} else if util.IsPlayingAndNotFold(state.GS.Players[index]) {
			state.GS.Players[index].Bets = append(state.GS.Players[index].Bets, 0)
			IncreasePots(index, chips)
		}
	}
	SetMaximumBet(util.SumPots(state.GS.Pots))
}

// OverwriteActionToBehindPlayers overwritten action with default
func OverwriteActionToBehindPlayers() {
	for index := range state.GS.Players {
		if util.IsPlayingAndNotFoldAndNotAllIn(state.GS.Players[index]) &&
			util.IsPlayerBehindTheTimeline(state.GS.Players[index]) {
			state.GS.Players[index].Action = state.GS.Players[index].Default
		}
	}
}

// BurnBet burn bet from player
func BurnBet(index int, burn int) {
	// if this player cannot pay all of it
	if burn > state.GS.Pots[index] {
		state.GS.Pots[index] = 0
	} else {
		state.GS.Pots[index] -= burn
	}
}
