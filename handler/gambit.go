package handler

import (
	"999k_engine/api"
	"999k_engine/constant"
	"999k_engine/engine"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"log"
	"os"
	"time"
)

// Initiate required variables
func Initiate(game engine.Gambit) {
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
	start := time.Now().Unix()
	state.GS.StartTableTime = start
	state.GS.FinishTableTime = start + (60 * 30)
	state.GS.IsTableStart = true
}

// FinishTable set table start
func FinishTable() {
	// Cannot set to 0 because 0 is a initialized value
	state.GS.FinishTableTime = time.Now().Unix() - 5
}

// StartGame set game start
func StartGame() {
	state.GS.IsGameStart = true
}

// IsTableStart true or false
func IsTableStart() bool {
	return state.GS.IsTableStart &&
		time.Now().Unix() < state.GS.FinishTableTime
}

// IsGameStart true or false
func IsGameStart() bool {
	return state.GS.IsGameStart
}

// IsEndRound if current time is more than finish time
func IsEndRound() bool {
	return state.GS.FinishRoundTime <= time.Now().Unix()
}

// IsInExtendFinishRoundTime when end round but still not actually ended
func IsInExtendFinishRoundTime() bool {
	return state.GS.FinishRoundTime > time.Now().Unix()
}

// IncreaseTurn to seperate player bets
func IncreaseTurn() {
	state.GS.Turn++
}

// IncreaseGameIndex every game start
func IncreaseGameIndex() {
	state.GS.GameIndex++
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
	if util.CountPlayerNotFoldAndNotAllIn(state.GS.Players) > 1 {
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
	}
	state.GS.FinishRoundTime = start
}

// MakePlayersReady make everyone
func MakePlayersReady() {
	for index := range state.GS.Players {
		player := &state.GS.Players[index]
		if player.ID == "" {
			continue
		}
		player.Cards = model.Cards{}
		player.Bets = []int{}
		player.IsPlaying = true
		player.IsEarned = false
		player.IsWinner = false
		player.Default = model.Action{Name: constant.Check}
		player.Action = model.Action{}
	}
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

// SetOtherActionsWhoAreNotPlaying set everyone action who are not playing
func SetOtherActionsWhoAreNotPlaying(action string) {
	// if others who are not playing then able to starttable or only stand
	for index, player := range state.GS.Players {
		// not a seat and not playing
		if player.ID != "" && !player.IsPlaying {
			state.GS.Players[index].Actions = Reducer(action, player.ID)
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
		if len(player.Cards) < maxcards {
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
			player := &state.GS.Players[index]
			if util.IsPlayingAndNotFold(*player) {
				player.Cards = append(player.Cards, Draw())
				player.CardAmount = len(player.Cards)
				if index == dealer {
					break
				}
			}
		}
	}
}

// FlushPlayers reset everything before new game and client no needs to see it
func FlushPlayers() {
	for index := range state.GS.Players {
		player := &state.GS.Players[index]
		player.IsPlaying = false
		player.Bets = []int{}
		player.Default = model.Action{}
		player.Action = model.Action{}
		player.Actions = model.Actions{}
		player.StartLine = 0
		player.DeadLine = 0
		player.IsEarned = false
	}
}

// ShortenTimeline shift timeline of everyone because someone take action
func ShortenTimeline(diff int64) {
	diff = util.Absolute(diff)
	for index, player := range state.GS.Players {
		if player.IsPlaying {
			state.GS.Players[index].StartLine -= diff
			state.GS.Players[index].DeadLine -= diff
		}
	}
	state.GS.FinishRoundTime = state.GS.FinishRoundTime - diff
}

// ShortenTimelineAfterTarget shift timeline of everyone behind target player
func ShortenTimelineAfterTarget(id string, second int64) {
	second = util.Absolute(second)
	_, caller := util.Get(state.GS.Players, id)
	for index, player := range state.GS.Players {
		// who start behind caller will be shifted
		if util.IsPlayingAndNotFoldAndNotAllIn(player) && player.StartLine >= caller.DeadLine {
			state.GS.Players[index].StartLine -= second
			state.GS.Players[index].DeadLine -= second
		}
	}
	state.GS.FinishRoundTime = state.GS.FinishRoundTime - second
}

// ExtendPlayerTimeline extend player timeline
// return boolean because it needs to validate wth another server
func ExtendPlayerTimeline(id string) bool {
	if !IsPlayerTurn(id) {
		return false
	}
	second := state.GS.Gambit.GetSettings().DecisionTime
	current, caller := util.Get(state.GS.Players, id)
	start := time.Now().Unix()
	diff := (start + second) - state.GS.Players[current].DeadLine
	state.GS.Players[current].StartLine = start
	state.GS.Players[current].DeadLine = start + second
	for index, player := range state.GS.Players {
		// who start behind caller will be shifted
		if util.IsPlayingAndNotFoldAndNotAllIn(player) && player.StartLine >= caller.DeadLine {
			state.GS.Players[index].StartLine += diff
			state.GS.Players[index].DeadLine += diff
		}
	}
	state.GS.FinishRoundTime += diff
	OverwriteActionToBehindPlayers()
	return true
}

// ExtendFinishRoundTime force this timeline to be ended
func ExtendFinishRoundTime() {
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
		if player := &state.GS.Players[index]; util.IsPlayingAndNotFoldAndNotAllIn(*player) {
			player.Chips -= chips
			player.WinLossAmount -= chips
			util.AddScoreboardWinAmount(player.ID, -chips)
			player.Bets = append(player.Bets, chips)
			IncreasePots(index, chips)
			// start with first element in pots
		} else if util.IsPlayingAndNotFold(*player) {
			player.Bets = append(player.Bets, 0)
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

// TryTerminate try to terminate the container
func TryTerminate() {
	// Check if current time is more than finish table time
	if time.Now().Unix() >= state.GS.FinishTableTime &&
		state.GS.FinishTableTime != 0 {
		log.Println("Table is timeout then terminate")
		// For force client to leave
		state.GS.IsTableExpired = true
		state.GS.IsTableStart = false
		// TODO call terminate api
		if os.Getenv("env") != "dev" {
			// Delay 5 second before send signal to hawkeye that please kill this container
			go func() {
				body, err := api.TableEnd()
				log.Println("Response from TableEnd", string(body), err)
				time.Sleep(time.Second * 3)
				body, err = api.Terminate()
				log.Println("Response from Terminate", string(body), err)
			}()
		}
	}
}
