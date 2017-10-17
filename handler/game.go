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
	return state.GS.FinishGameTime.Unix() <= time.Now().Unix()
}

// IsBeforeDeadLine if player do something before deadline
func IsBeforeDeadLine(id int) bool {
	// deadline := state.GS.Players[id].DeadLine.Unix()
	// startline := deadline
	return true
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
		if !player.IsPlaying {
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
func CreateTimeLine(limit int) {
	pos, loop := 0, 0
	start, amount := time.Now(), len(state.GS.Players)
	index, _ := util.FindDealer(state.GS.Players)
	for loop < amount {
		next := (index + 1) % amount
		player := state.GS.Players[next]
		if player.IsPlaying {
			pos++
			state.GS.Players[next].Action = model.Action{Name: constant.Check}
			state.GS.Players[next].DeadLine = start.Add(time.Second * time.Duration(pos*limit))
		}
		index++
		loop++
	}
	state.GS.StartGameTime = start
	state.GS.FinishGameTime = start.Add(time.Second * time.Duration(pos*limit))
}

// IsFullHand check if hold max cards
func IsFullHand(maxcards int) bool {
	for _, player := range state.GS.Players {
		if !player.IsPlaying {
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
		if !player.IsPlaying {
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
		for round <= playerAmount {
			start++
			round++
			index = start % playerAmount
			// skip empty seat
			if state.GS.Players[index].IsPlaying {
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
	state.GS.IsGameStart = false
}
