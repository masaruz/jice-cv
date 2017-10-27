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

// IncreaseTurn to seperate player bets
func IncreaseTurn() {
	state.GS.Turn++
}

// GetTurn get current turn number
func GetTurn() int {
	return state.GS.Turn
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

// AssignWinners find a winner by evaluate his cards
func AssignWinners() {
	hscore := -1
	hbonus := -1
	pos := -1
	// hkind := ""
	// winner := model.Player{}
	for i := 0; i < util.CountPlayerNotFold(state.GS.Players); i++ {
		for index, player := range state.GS.Players {
			if !util.IsPlayingAndNotFold(player) ||
				len(player.Cards) == 0 ||
				player.IsEarned {
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
		if pos != -1 {
			for index, player := range state.GS.Players {
				if !player.IsPlaying {
					continue
				}
				playerbet := util.SumBet(player)
				winnerbet := util.SumBet(state.GS.Players[pos])
				// if winner has higher bet
				if winnerbet > playerbet {
					state.GS.Players[pos].Chips += playerbet
					// player bet will be 0
					BurnBet(player.ID, playerbet)
				} else {
					// if winner has lower bet
					state.GS.Players[pos].Chips += winnerbet
					// if not caller
					if index != pos {
						BurnBet(player.ID, winnerbet)
					}
				}
			}
			state.GS.Players[pos].IsEarned = true
			BurnBet(state.GS.Players[pos].ID, util.SumBet(state.GS.Players[pos]))
			hscore = -1
			hbonus = -1
			pos = -1
		}
	}
}

// CreateTimeLine set timeline for game and any players
func CreateTimeLine(decisionTime int64) {
	// need at least one competitors
	if util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) <= 1 {
		return
	}
	loop := 0
	delay := int64(0)
	start, amount := time.Now().Unix(), len(state.GS.Players)
	state.GS.StartRoundTime = start
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

// AssignPlayersCheckOrAllIn if player has no chips then let them allin otherwise check
func AssignPlayersCheckOrAllIn() {
	for index, player := range state.GS.Players {
		if !util.IsPlayingAndNotFold(player) {
			continue
		}
		action := model.Action{}
		if state.GS.Players[index].Chips == 0 {
			action = model.Action{Name: constant.AllIn}
		} else {
			action = model.Action{Name: constant.Check}
		}
		state.GS.Players[index].Action = action
		state.GS.Players[index].Default = action
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
		if !util.IsPlayingAndNotFold(player) {
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
		// state.GS.Players[index].Cards = model.Cards{}
		state.GS.Players[index].Bets = []int{}
		state.GS.Players[index].Default = model.Action{}
		state.GS.Players[index].Action = model.Action{}
		state.GS.Players[index].Actions = model.Actions{}
		state.GS.Players[index].StartLine = 0
		state.GS.Players[index].DeadLine = 0
		state.GS.Players[index].IsEarned = false
	}
	state.GS.Pots = []int{}
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

// ForceEndTimeline force this timeline to be ended
func ForceEndTimeline() {
	state.GS.FinishRoundTime = time.Now().Unix()
}

// ShiftPlayerToEndOfTimeline shift player to the end of timeline
func ShiftPlayerToEndOfTimeline(id string, second int64) {
	index, _ := util.Get(state.GS.Players, id)
	finishRoundTime := state.GS.FinishRoundTime
	state.GS.Players[index].StartLine = finishRoundTime
	state.GS.Players[index].DeadLine = finishRoundTime + second
	state.GS.FinishRoundTime = finishRoundTime + second
}

// ShiftPlayersToEndOfTimeline shift current and prev player to the end of timeline
func ShiftPlayersToEndOfTimeline(id string, second int64) {
	start, _ := util.Get(state.GS.Players, id)
	round, amount := 0, len(state.GS.Players)
	for round < amount {
		start++
		round++
		index := start % amount
		// force shift to players who is in game not allin and behine the timeline
		if util.IsPlayingAndNotFoldAndNotAllIn(state.GS.Players[index]) &&
			util.IsPlayerBehindTheTimeline(state.GS.Players[index]) {
			ShiftPlayerToEndOfTimeline(state.GS.Players[index].ID, second)
		}
	}
}

// InvestToPots added bet to everyone base on turn
func InvestToPots(chips int) {
	// initiate bet value to players
	for index := range state.GS.Players {
		if util.IsPlayingAndNotFoldAndNotAllIn(state.GS.Players[index]) {
			state.GS.Players[index].Chips -= chips
			state.GS.Players[index].Bets = append(state.GS.Players[index].Bets, chips)
			SetMaximumBet(util.SumBets(state.GS.Players)) // start with first element in pots
		} else if util.IsPlayingAndNotFold(state.GS.Players[index]) {
			state.GS.Players[index].Bets = append(state.GS.Players[index].Bets, 0)
			SetMaximumBet(util.SumBets(state.GS.Players))
		}
	}
}
