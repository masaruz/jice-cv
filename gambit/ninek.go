package gambit

import (
	"999k_engine/api"
	"999k_engine/constant"
	"999k_engine/engine"
	"999k_engine/handler"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"math"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

// NineK is 9K
type NineK struct {
	MaxPlayers      int
	DecisionTime    int64
	DelayNextRound  int64
	FinishGameDelay int64
	MaxAFKCount     int
	BlindsSmall     int
	BlindsBig       int
	BuyInMin        int
	BuyInMax        int
	Rake            float64 // percentage
	Cap             float64 // cap of rake
	GPSRestrcited   bool
}

// Init deck and environment variables
func (game NineK) Init() {
	// create counting afk
	state.GS.PlayerPots = make([]int, game.MaxPlayers)
	state.GS.AFKCounts = make([]int, game.MaxPlayers)
	// set the seats
	handler.CreateSeats(game.MaxPlayers)
	handler.SetMinimumBet(game.BlindsBig)
}

// Start game
func (game NineK) Start() bool {
	util.Print("Try to start")
	if handler.IsTableStart() &&
		!handler.IsGameStart() &&
		!handler.IsInExtendFinishRoundTime() &&
		!handler.IsTableExpired() {
		// everyone is assumed afk
		state.Snapshot.DoActions = make([]bool, game.MaxPlayers)
		state.Snapshot.Rakes = make(map[string]float64)
		state.Snapshot.PlayerPots = make([]int, game.MaxPlayers)
		handler.InitPots(&state.Snapshot)
		// filter players who are not ready to play
		for index := range state.Snapshot.Players {
			player := &state.Snapshot.Players[index]
			// If this position is empty seat continue
			if player.ID == "" {
				continue
			}
			// If someone has requested for topup then call buyin
			if player.TopUp.IsRequest {
				handler.TopUp(player.ID)
			}
			// If player has no chip enough
			if int(math.Floor(player.Chips)) < game.GetSettings().BlindsSmall {
				// Force to stand
				if !handler.Stand(player.ID, false) {
					return false
				}
				// If player has minimum chip for able to play
			} else if state.Snapshot.AFKCounts[index] >= game.MaxAFKCount {
				util.Print(player.ID, "Is AFK")
				// Force to stand
				if !handler.Stand(player.ID, false) {
					return false
				}
				// In case this room check GPS location
			} else if game.GPSRestrcited {
				// Check to others
				for _, competitor := range state.Snapshot.Players {
					if competitor.ID == "" || competitor.ID == player.ID {
						continue
					}
					if util.Distance(*player, competitor) <= 50 {
						util.Print(player.ID, "Is nearby someone")
						if !handler.Stand(player.ID, false) {
							return false
						}
					}
				}
			}
		}
		// After filtered with the critiria
		// if there are more than 2 players are sitting
		if util.CountSitting(state.Snapshot.Players) >= 2 {
			// Increase gameindex for backend process ex. realtime-data, analytic
			util.Print("Gameindex increased")
			util.Print("Prepare to start game")
			if state.Snapshot.Env != "dev" {
				// Request to start game
				body, err := api.StartGame()
				util.Print("Response from StartGame", string(body), err)
				resp := &api.Response{}
				json.Unmarshal(body, resp)
				// Is there any error when start game
				if resp.Error != (api.Error{}) {
					return false
				}
			}
			state.Snapshot.GameIndex++
			// set players to be ready
			handler.PreparePlayers(true)
			handler.StartGame()
			handler.SetMinimumBet(game.BlindsBig)
			// let all players bets to the pots
			handler.PlayersInvestToPots(game.BlindsBig)
			handler.MergePots(&state.Snapshot)
			// start turn
			handler.IncreaseTurn()
			// start new bets
			handler.PlayersInvestToPots(0)
			handler.MergePots(&state.Snapshot)
			handler.SetOtherActions("", constant.Check)
			handler.SetOtherDefaultAction("", constant.Check)
			handler.SetDealer()
			handler.BuildDeck()
			handler.Shuffle()
			handler.CreateTimeLine(game.DecisionTime)
			handler.Deal(2, game.MaxPlayers)
			util.Print("2 cards dealed")
			handler.SetPlayersRake(game.Rake, game.Cap*float64(game.BlindsBig))
			util.Print("Start Success")
			return true
		}
		for index := range state.Snapshot.Players {
			comp := &state.Snapshot.Players[index]
			if comp.ID == "" {
				continue
			}
			comp.CardAmount = 0
		}
		handler.PreparePlayers(false)
		// Need to update state because number of players might be changed
		state.GS = util.CloneState(state.Snapshot)
	} else if !handler.IsGameStart() &&
		!handler.IsInExtendFinishRoundTime() &&
		time.Now().Unix() >= state.Snapshot.FinishTableTime &&
		state.Snapshot.FinishTableTime != 0 {
		state.Snapshot.IsTableExpired = true
		handler.TryTerminate()
		state.GS = util.CloneState(state.Snapshot)
	}
	util.Print("Start Failed")
	return false
}

// NextRound game after round by round
func (game NineK) NextRound() bool {
	util.Print("Try to next round")
	// Assume every player must be default
	handler.OverwriteActionToBehindPlayers()
	// Game must be start
	// All players must not has 3 cards
	// All bets must be equal
	// Now must be more than finish round time
	if handler.IsGameStart() && !handler.IsFullHand(3) &&
		handler.BetsEqual() && handler.IsEndRound() &&
		util.CountPlayerNotFold(state.Snapshot.Players) > 1 {
		// Initialize values
		handler.Deal(1, game.MaxPlayers)
		handler.SetMinimumBet(game.BlindsBig)
		// If there are figthers
		if util.CountPlayerNotFoldAndNotAllIn(state.Snapshot.Players) > 1 {
			handler.SetOtherActions("", constant.Check)
			handler.SetOtherDefaultAction("", constant.Check)
		}
		handler.CreateTimeLine(game.DecisionTime)
		handler.PlayersInvestToPots(0)
		handler.MergePots(&state.Snapshot)
		util.Print("1 cards dealed")
		handler.IncreaseTurn()
		// no one is assumed afk
		state.Snapshot.DoActions = make([]bool, game.MaxPlayers)
		util.Print("Next round Success")
		return true
	}
	util.Print("Next round Failed")
	return false
}

// Finish game
func (game NineK) Finish() bool {
	util.Print("Try to finish")
	handler.OverwriteActionToBehindPlayers()
	// no others to play with or all players have 3 cards but bet is not equal
	if handler.IsGameStart() && ((util.CountPlayerNotFold(state.Snapshot.Players) <= 1) ||
		// if has 3 cards bet equal
		(handler.IsFullHand(3) && handler.BetsEqual() && handler.IsEndRound())) {
		util.Print("Prepare to finish game")
		handler.PlayersInvestToPots(0)
		handler.MergePots(&state.Snapshot)
		// calculate afk players
		for index, doaction := range state.Snapshot.DoActions {
			// Skip empty players
			if !util.IsPlayingAndNotFoldAndNotAllIn(state.Snapshot.Players[index]) {
				continue
			}
			// If this player done something in this round
			if doaction {
				state.Snapshot.AFKCounts[index] = 0
			} else {
				// If this player never done action in this game
				state.Snapshot.AFKCounts[index]++
			}
		}
		// Extend more time for client to play animation after game finished
		handler.ExtendFinishRoundTime()
		// Find winner and added their rewards
		hscore := -1
		hbonus1 := -1
		hbonus2 := -1
		hbonus3 := -1
		hbonus4 := -1
		pos := -1
		util.Print("Find the winner(s)")
		for i := 0; i < len(state.Snapshot.Players); i++ {
			// Evaluate score from everyone's hand
			// Find a winner
			for index, player := range state.Snapshot.Players {
				if !util.IsPlayingAndNotFold(player) ||
					len(player.Cards) == 0 ||
					player.IsEarned {
					continue
				}
				scores, _ := game.Evaluate(player.Cards)
				score := scores[0]
				bonus1 := scores[1]
				bonus2 := scores[2]
				bonus3 := scores[3]
				bonus4 := scores[4]
				highscore := (hscore < score) ||
					(hscore == score && hbonus1 < bonus1) ||
					(hscore == score && hbonus1 == bonus1 && hbonus2 < bonus2) ||
					(hscore == score && hbonus1 == bonus1 && hbonus2 == bonus2 && hbonus3 < bonus3) ||
					(hscore == score && hbonus1 == bonus1 && hbonus2 == bonus2 && hbonus3 == bonus3 && hbonus4 < bonus4)
				if highscore {
					hscore = score
					hbonus1 = bonus1
					hbonus2 = bonus2
					hbonus3 = bonus3
					hbonus4 = bonus4
					pos = index
				}
			}
			// This mean we found some winners
			if pos != -1 {
				winner := &state.Snapshot.Players[pos]
				handler.AssignWinnerToPots(&state.Snapshot, winner.ID)
				hscore = -1
				hbonus1 = -1
				hbonus2 = -1
				hbonus3 = -1
				pos = -1
			}
		}
		// Update scoreboard with winloss amount value
		for _, player := range state.Snapshot.Players {
			if !player.IsPlaying {
				continue
			}
			handler.UpdateWinningsAmount(player.ID, player.WinLossAmount)
		}
		if state.Snapshot.Env != "dev" {
			body, err := api.SaveSettlements()
			util.Print("Response from SaveSettlements", string(body), err)
			resp := &api.Response{}
			json.Unmarshal(body, resp)
			// Is there any error when start game
			if resp.Error != (api.Error{}) {
				return false
			}
			body, err = api.UpdateRealtimeData()
			util.Print("Response from UpdateRealtimeData", string(body), err)
			resp = &api.Response{}
			json.Unmarshal(body, resp)
			// Is there any error when start game
			if resp.Error != (api.Error{}) {
				return false
			}
		}
		// Revert minimum bet
		handler.SetMinimumBet(game.BlindsBig)
		handler.FlushPlayers()
		state.Snapshot.Turn = 0
		state.Snapshot.IsGameStart = false
		handler.SaveHistory()
		util.Print("Finish Success")
		return true
	}
	util.Print("Finish Failed")
	return false
}

// Check is doing nothing but only shift the timeline
func (game NineK) Check(id string) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.Snapshot.Players, id)
	player := &state.Snapshot.Players[index]
	if !util.IsPlayingAndNotFoldAndNotAllIn(*player) {
		return false
	}
	// Cannot check if player has less bet than highest
	if player.Bets[state.Snapshot.Turn] <
		util.GetHighestBetInTurn(state.Snapshot.Turn, state.Snapshot.Players) {
		return false
	}
	state.Snapshot.DoActions[index] = true
	player.Default = model.Action{Name: constant.Check}
	player.Action = model.Action{Name: constant.Check}
	diff := time.Now().Unix() - player.DeadLine
	handler.OverwriteActionToBehindPlayers()
	handler.ShortenTimeline(diff)
	return true
}

// Bet is raising bet to the target
func (game NineK) Bet(id string, chips int) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	return game.pay(id, chips, constant.Bet)
}

// Raise is raising bet to the target
func (game NineK) Raise(id string, chips int) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.Snapshot.Players, id)
	// not less than minimum
	if state.Snapshot.Players[index].Bets[state.Snapshot.Turn]+chips <= state.Snapshot.MinimumBet {
		return false
	}
	return game.pay(id, chips, constant.Raise)
}

// Call is raising bet to the highest bet
func (game NineK) Call(id string) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.Snapshot.Players, id)
	player := &state.Snapshot.Players[index]
	if !util.IsPlayingAndNotFoldAndNotAllIn(*player) {
		return false
	}
	chips := util.GetHighestBetInTurn(state.Snapshot.Turn, state.Snapshot.Players) -
		player.Bets[state.Snapshot.Turn]
	chipsDecimal := decimal.NewFromFloat(float64(chips))
	// cannot call more than player's chips
	if int(math.Floor(player.Chips)) < chips || chips == 0 {
		return false
	}
	state.Snapshot.DoActions[index] = true
	player.Chips, _ = decimal.NewFromFloat(player.Chips).Sub(chipsDecimal).Float64()
	player.WinLossAmount, _ = decimal.NewFromFloat(player.WinLossAmount).Sub(chipsDecimal).Float64()
	player.Bets[state.Snapshot.Turn] += chips
	playerAction := model.Action{}
	if math.Floor(player.Chips) == 0 {
		playerAction = model.Action{Name: constant.AllIn}
	} else {
		playerAction = model.Action{Name: constant.Call}
	}
	player.Default = playerAction
	player.Action = playerAction
	player.Actions = game.Reducer(constant.Check, id)
	handler.IncreasePlayerPot(index, chips)
	// set action of everyone
	handler.OverwriteActionToBehindPlayers()
	handler.SetMaximumBet(util.SumPots(state.Snapshot.PlayerPots))
	// others need to know what to do next
	handler.SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - player.DeadLine
	handler.ShortenTimeline(diff)
	// set players rake
	handler.SetPlayersRake(game.Rake, game.Cap*float64(game.BlindsBig))
	return true
}

// AllIn give all chips
func (game NineK) AllIn(id string) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.Snapshot.Players, id)
	player := &state.Snapshot.Players[index]
	if !util.IsPlayingAndNotFoldAndNotAllIn(*player) {
		return false
	}
	chips := int(math.Floor(player.Chips))
	// not more than maximum
	if player.Bets[state.Snapshot.Turn]+chips > state.Snapshot.MaximumBet {
		return false
	}
	highestbet := util.GetHighestBetInTurn(state.Snapshot.Turn, state.Snapshot.Players)
	state.Snapshot.DoActions[index] = true
	player.Bets[state.Snapshot.Turn] += chips
	chipsDecimal := decimal.NewFromFloat(float64(chips))
	player.Chips, _ = decimal.NewFromFloat(player.Chips).Sub(chipsDecimal).Float64()
	player.WinLossAmount, _ = decimal.NewFromFloat(player.WinLossAmount).Sub(chipsDecimal).Float64()
	player.Default = model.Action{Name: constant.AllIn}
	player.Action = model.Action{Name: constant.AllIn}
	player.Actions = game.Reducer(constant.Check, id)
	handler.IncreasePlayerPot(index, chips)
	handler.SetMaximumBet(util.SumPots(state.Snapshot.PlayerPots))
	// set action of everyone
	handler.OverwriteActionToBehindPlayers()
	// others automatic set to fold as default
	handler.SetOtherDefaultAction(id, constant.Fold)
	// others need to know what to do next
	handler.SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - player.DeadLine
	handler.ShortenTimeline(diff)
	// duration extend the timeline
	if player.Bets[state.Snapshot.Turn] > highestbet {
		handler.SetMinimumBet(player.Bets[state.Snapshot.Turn])
		handler.ShiftPlayersToEndOfTimeline(id, game.DecisionTime)
	}
	// set players rake
	handler.SetPlayersRake(game.Rake, game.Cap*float64(game.BlindsBig))
	return true
}

// Fold quit the game but still lost bet
func (game NineK) Fold(id string) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.Snapshot.Players, id)
	player := &state.Snapshot.Players[index]
	if !util.IsPlayingAndNotFoldAndNotAllIn(*player) {
		return false
	}
	state.Snapshot.DoActions[index] = true
	player.Default = model.Action{Name: constant.Fold}
	player.Action = model.Action{Name: constant.Fold}
	player.Actions = game.Reducer(constant.Fold, id)
	diff := time.Now().Unix() - player.DeadLine
	handler.OverwriteActionToBehindPlayers()
	handler.ShortenTimeline(diff)
	return true
}

// End game
func (game NineK) End() {}

// Reducer reduce the action and when receive the event
func (game NineK) Reducer(event string, id string) model.Actions {
	_, player := util.Get(state.Snapshot.Players, id)
	chip := int(math.Floor(player.Chips))
	topupAction := handler.GetTopUpHint(id)
	extendAction := model.Action{
		Name: constant.ExtendDecisionTime,
		Hints: model.Hints{
			model.Hint{Name: "amount", Type: "integer", Value: 15},
		}}
	standAction := model.Action{Name: constant.Stand}
	switch event {
	case constant.Bet:
		playerchips := chip + player.Bets[state.Snapshot.Turn]
		// highest bet in that turn
		highestbet := util.GetHighestBetInTurn(state.Snapshot.Turn, state.Snapshot.Players)
		playerbet := player.Bets[state.Snapshot.Turn]
		// raise must be highest * 2
		raise := highestbet * 2
		// all sum bets
		pots := util.SumPots(state.Snapshot.PlayerPots)
		if highestbet <= playerbet {
			return game.Reducer(constant.Check, id)
		}
		// no more than pots
		if raise > pots {
			raise = pots
		}
		if playerchips <= highestbet {
			return model.Actions{
				model.Action{Name: constant.Fold},
				model.Action{Name: constant.AllIn,
					Hints: model.Hints{
						model.Hint{
							Name: "amount", Type: "integer", Value: chip}}},
				extendAction, standAction, topupAction}
		}
		diff := highestbet - playerbet
		if playerchips < raise {
			return model.Actions{
				model.Action{Name: constant.Fold},
				model.Action{Name: constant.Call,
					Hints: model.Hints{
						model.Hint{
							Name: "amount", Type: "integer", Value: diff}}},
				model.Action{Name: constant.AllIn,
					Hints: model.Hints{
						model.Hint{
							Name: "amount", Type: "integer", Value: chip}}},
				extendAction, standAction, topupAction}
		}
		// maximum will be player's chips if not enough
		maximum := 0
		if state.Snapshot.MaximumBet > playerchips {
			maximum = playerchips
		} else {
			maximum = state.Snapshot.MaximumBet
		}
		return model.Actions{
			model.Action{Name: constant.Fold},
			model.Action{Name: constant.Call,
				Hints: model.Hints{
					model.Hint{
						Name: "amount", Type: "integer", Value: diff}}},
			model.Action{Name: constant.Raise,
				Parameters: model.Parameters{
					model.Parameter{
						Name: "amount", Type: "integer"}},
				Hints: model.Hints{
					model.Hint{
						Name: "amount", Type: "integer", Value: raise - playerbet},
					model.Hint{
						Name: "amount_max", Type: "integer", Value: maximum - playerbet}}},
			extendAction, standAction, topupAction}
	case constant.Fold:
		return model.Actions{standAction, topupAction}
	default:
		if math.Floor(player.Chips) == 0 {
			return model.Actions{
				model.Action{Name: constant.Fold},
				model.Action{Name: constant.Check},
				extendAction, standAction, topupAction}
		}
		// maximum will be player's chips if not enough
		maximum := 0
		if state.Snapshot.MaximumBet > chip {
			maximum = chip
		} else {
			maximum = state.Snapshot.MaximumBet
		}
		return model.Actions{
			model.Action{Name: constant.Fold},
			model.Action{Name: constant.Check},
			model.Action{Name: constant.Bet,
				Parameters: model.Parameters{
					model.Parameter{
						Name: "amount", Type: "integer"}},
				Hints: model.Hints{
					model.Hint{
						Name: "amount", Type: "integer", Value: state.Snapshot.MinimumBet},
					model.Hint{
						Name: "amount_max", Type: "integer", Value: maximum}}},
			extendAction, standAction, topupAction}
	}
}

// Evaluate score of player cards
func (game NineK) Evaluate(values []int) (scores []int, kind string) {
	sorted := make([]int, len(values))
	copy(sorted, values)
	sort.Ints(sorted)
	// all cards are same number
	if threeOfAKind(sorted) {
		return summary(constant.ThreeOfAKind, sorted)
	}
	// all cards are in seqence and same suit
	if straightFlush(sorted) {
		return summary(constant.StraightFlush, sorted)
	}
	// all cards are J or Q or K
	if royal(sorted) {
		return summary(constant.Royal, sorted)
	}
	// all cards are in sequnce
	if straight(sorted) {
		return summary(constant.Straight, sorted)
	}
	// all cards are same color and kind
	if flush(sorted) {
		return summary(constant.Flush, sorted)
	}
	return summary(constant.Nothing, sorted)
}

// GetSettings return settings variables
func (game NineK) GetSettings() engine.Settings {
	return engine.Settings{
		MaxPlayers:      game.MaxPlayers,
		DecisionTime:    game.DecisionTime,
		MaxAFKCount:     game.MaxAFKCount,
		BlindsSmall:     game.BlindsSmall,
		BlindsBig:       game.BlindsSmall,
		BuyInMin:        game.BuyInMin,
		BuyInMax:        game.BuyInMax,
		Rake:            game.Rake,
		Cap:             game.Cap,
		FinishGameDelay: game.FinishGameDelay,
		DelayNextRound:  game.DelayNextRound,
		GPSRestrcited:   game.GPSRestrcited,
	}
}

func (game NineK) pay(id string, chips int, action string) bool {
	index, _ := util.Get(state.Snapshot.Players, id)
	player := &state.Snapshot.Players[index]
	if !util.IsPlayingAndNotFoldAndNotAllIn(*player) {
		return false
	}
	// not less than minimum
	if player.Bets[state.Snapshot.Turn]+chips < state.Snapshot.MinimumBet {
		return false
	}
	// not more than maximum
	if player.Bets[state.Snapshot.Turn]+chips > state.Snapshot.MaximumBet {
		return false
	}
	// cannot bet more than player's chips
	if int(math.Floor(player.Chips)) < chips {
		return false
	}
	state.Snapshot.DoActions[index] = true
	// added value to the bet in this turn
	chipsDecimal := decimal.NewFromFloat(float64(chips))
	player.Chips, _ = decimal.NewFromFloat(player.Chips).Sub(chipsDecimal).Float64()
	player.WinLossAmount, _ = decimal.NewFromFloat(player.WinLossAmount).Sub(chipsDecimal).Float64()
	player.Bets[state.Snapshot.Turn] += chips
	// broadcast to everyone that I bet
	playerAction := model.Action{}
	if math.Floor(player.Chips) == 0 {
		playerAction = model.Action{Name: constant.AllIn}
	} else {
		playerAction = model.Action{Name: action}
	}
	player.Default = playerAction
	player.Action = playerAction
	player.Actions = game.Reducer(constant.Check, id)
	handler.IncreasePlayerPot(index, chips)
	// assign minimum bet
	handler.SetMinimumBet(player.Bets[state.Snapshot.Turn])
	// assign maximum bet
	handler.SetMaximumBet(util.SumPots(state.Snapshot.PlayerPots))
	// set action of everyone
	handler.OverwriteActionToBehindPlayers()
	// others automatic set to fold as default
	handler.SetOtherDefaultAction(id, constant.Fold)
	// others need to know what to do next
	handler.SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - player.DeadLine
	handler.ShortenTimeline(diff)
	// duration extend the timeline
	handler.ShiftPlayersToEndOfTimeline(id, game.DecisionTime)
	// set players rake
	handler.SetPlayersRake(game.Rake, game.Cap*float64(game.BlindsBig))
	return true
}
