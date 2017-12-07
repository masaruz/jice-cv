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
	"sort"
	"time"
)

// NineK is 9K
type NineK struct {
	MaxPlayers   int
	DecisionTime int64
	MaxAFKCount  int
	BlindsSmall  int
	BlindsBig    int
	BuyInMin     int
	BuyInMax     int
	Rake         float64 // percentage
	Cap          float64 // cap of rake
}

// Init deck and environment variables
func (game NineK) Init() {
	// create counting afk
	state.GS.Pots = make([]int, game.MaxPlayers)
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
		!handler.IsInExtendFinishRoundTime() {
		// filter players who are not ready to play
		for index := range state.Snapshot.Players {
			player := &state.Snapshot.Players[index]
			// If this position is empty seat continue
			if player.ID == "" {
				continue
			}
			// If player has no chip enough
			if player.Chips < game.GetSettings().BlindsSmall {
				// If in development assign chips and continue
				if state.Snapshot.Env == "dev" {
					player.Chips = game.GetSettings().BuyInMin
					continue
				}
				// Make sure this player ready to buyin
				// Cashback can be fail if they not buyin yet
				body, err := api.CashBack(player.ID)
				util.Print("Response from CashBack", string(body), err)
				resp := &api.Response{}
				json.Unmarshal(body, resp)
				// If cashback error
				if resp.Error != (api.Error{}) && resp.Error.StatusCode != 404 {
					// Force to stand
					if !handler.Stand(player.ID) {
						return false
					}
					continue
				}
				// After cashback success set chips to be 0
				player.Chips = 0
				// Need request to server for buyin
				body, err = api.BuyIn(player.ID, game.GetSettings().BuyInMin)
				util.Print("Response from BuyIn", string(body), err)
				resp = &api.Response{}
				json.Unmarshal(body, resp)
				// BuyIn must be successful
				if resp.Error != (api.Error{}) {
					// Force to stand
					if !handler.Stand(player.ID) {
						return false
					}
					continue
				}
				util.Print("Buy-in success")
				// Assign how much they buy-in
				player.Chips = game.GetSettings().BuyInMin
				// Update scoreboard
				// If actually buyin success
				scoreboard, sbindex := util.GetScoreboard(player.ID)
				// If not found player in scoreboard then add them
				if sbindex == -1 {
					state.Snapshot.Scoreboard = append(state.Snapshot.Scoreboard, model.Scoreboard{
						UserID:      player.ID,
						DisplayName: player.Name,
						BuyInAmount: player.Chips,
					})
				} else {
					scoreboard.BuyInAmount += player.Chips
				}
			}
			// If player has minimum chip for able to play
			if state.Snapshot.AFKCounts[index] >= game.MaxAFKCount {
				util.Print(player.ID, "Is AFK")
				// Force to stand
				if !handler.Stand(player.ID) {
					return false
				}
				continue
			}
		}
		// After filtered with the critiria
		// if there are more than 2 players are sitting
		if util.CountSitting(state.Snapshot.Players) >= 2 {
			// Increase gameindex for backend process ex. realtime-data, analytic
			state.Snapshot.GameIndex++
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
					state.Snapshot.GameIndex--
					return false
				}
			}
			// everyone is assumed afk
			state.Snapshot.DoActions = make([]bool, game.MaxPlayers)
			state.Snapshot.Rakes = make(map[string]float64)
			state.Snapshot.Pots = make([]int, game.MaxPlayers)
			// set players to be ready
			handler.PreparePlayers(true)
			handler.StartGame()
			handler.SetMinimumBet(game.BlindsBig)
			// let all players bets to the pots
			handler.PlayersInvestToPots(game.BlindsBig)
			// start turn
			handler.IncreaseTurn()
			// start new bets
			handler.PlayersInvestToPots(0)
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
		handler.PreparePlayers(false)
		// Need to update state because number of players might be changed
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
		handler.SetOtherActions("", constant.Check)
		handler.SetOtherDefaultAction("", constant.Check)
		handler.CreateTimeLine(game.DecisionTime)
		handler.PlayersInvestToPots(0)
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
		hbonus := -1
		pos := -1
		util.Print("Find the winner(s)")
		// Evaluate score from everyone's hand
		for i := 0; i < len(state.Snapshot.Players); i++ {
			for index, player := range state.Snapshot.Players {
				if !util.IsPlayingAndNotFold(player) ||
					len(player.Cards) == 0 ||
					player.IsEarned {
					continue
				}
				scores, _ := game.Evaluate(player.Cards)
				score := scores[0]
				bonus := scores[1]
				if hscore < score {
					hscore = score
					hbonus = bonus
					pos = index
				} else if hscore == score && hbonus < bonus {
					hscore = score
					hbonus = bonus
					pos = index
				}
			}
			// This mean we found some winners
			if pos != -1 {
				winner := &state.Snapshot.Players[pos]
				for poti, pot := range state.Snapshot.Pots {
					if pot == 0 {
						continue
					}
					playerbet := pot
					winnerbet := state.Snapshot.Pots[pos]
					earnedbet := 0
					if winnerbet > playerbet {
						// If winner has higher bet
						earnedbet = playerbet
					} else {
						// If winner has lower bet
						earnedbet = winnerbet
					}
					if earnedbet != 0 {
						winner.Chips += earnedbet
						winner.WinLossAmount += earnedbet
						util.AddScoreboardWinAmount(winner.ID, earnedbet)
						earnedplayers := util.CountPlayerAlreadyEarned(state.Snapshot.Players)
						if util.CountPlayerNotFold(state.Snapshot.Players)-earnedplayers > 1 || earnedplayers == 0 {
							winner.IsWinner = true
						}
					}
					// If not caller
					if poti != pos {
						handler.BurnBet(poti, earnedbet)
					}
				}
				winner.IsEarned = true
				handler.BurnBet(pos, util.SumBet(*winner))
				hscore = -1
				hbonus = -1
				pos = -1
			}
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
		// Check if table expired then terminate
		handler.TryTerminate()
		state.Snapshot.Turn = 0
		state.Snapshot.IsGameStart = false
		util.Print("Finish Success")
		return true
	}
	// Check if table expired then terminate
	handler.TryTerminate()
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
	chips := util.GetHighestBetInTurn(state.Snapshot.Turn, state.Snapshot.Players) -
		player.Bets[state.Snapshot.Turn]
	// cannot call more than player's chips
	if player.Chips < chips || chips == 0 {
		return false
	}
	util.AddScoreboardWinAmount(player.ID, -chips)
	state.Snapshot.DoActions[index] = true
	player.Chips -= chips
	player.WinLossAmount -= chips
	player.Bets[state.Snapshot.Turn] += chips
	player.Default = model.Action{Name: constant.Call}
	player.Action = model.Action{Name: constant.Call}
	player.Actions = game.Reducer(constant.Check, id)
	handler.IncreasePots(index, chips)
	// set action of everyone
	handler.OverwriteActionToBehindPlayers()
	handler.SetMaximumBet(util.SumPots(state.Snapshot.Pots))
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
	chips := player.Chips
	// not more than maximum
	if player.Bets[state.Snapshot.Turn]+chips > state.Snapshot.MaximumBet {
		return false
	}
	util.AddScoreboardWinAmount(player.ID, -chips)
	state.Snapshot.DoActions[index] = true
	player.Bets[state.Snapshot.Turn] += chips
	player.WinLossAmount -= chips
	player.Chips = 0
	player.Default = model.Action{Name: constant.AllIn}
	player.Action = model.Action{Name: constant.AllIn}
	player.Actions = game.Reducer(constant.Check, id)
	handler.IncreasePots(index, chips)
	handler.SetMaximumBet(util.SumPots(state.Snapshot.Pots))
	// set action of everyone
	handler.OverwriteActionToBehindPlayers()
	// others automatic set to fold as default
	handler.SetOtherDefaultAction(id, constant.Fold)
	// others need to know what to do next
	handler.SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - player.DeadLine
	handler.ShortenTimeline(diff)
	// duration extend the timeline
	if player.Bets[state.Snapshot.Turn] >= util.GetHighestBetInTurn(state.Snapshot.Turn, state.Snapshot.Players) {
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
	extendAction := model.Action{
		Name: constant.ExtendDecisionTime,
		Hints: model.Hints{
			model.Hint{Name: "amount", Type: "integer", Value: 15},
		}}
	switch event {
	case constant.Check:
		_, player := util.Get(state.Snapshot.Players, id)
		if player.Chips == 0 {
			return model.Actions{
				model.Action{Name: constant.Fold},
				model.Action{Name: constant.Check},
				extendAction}
		}
		// maximum will be player's chips if not enough
		maximum := 0
		if state.Snapshot.MaximumBet > player.Chips {
			maximum = player.Chips
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
			extendAction}
	case constant.Bet:
		_, player := util.Get(state.Snapshot.Players, id)
		playerchips := player.Chips + player.Bets[state.Snapshot.Turn]
		// highest bet in that turn
		highestbet := util.GetHighestBetInTurn(state.Snapshot.Turn, state.Snapshot.Players)
		playerbet := player.Bets[state.Snapshot.Turn]
		// raise must be highest * 2
		raise := highestbet * 2
		// all sum bets
		pots := util.SumPots(state.Snapshot.Pots)
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
							Name: "amount", Type: "integer", Value: player.Chips}}},
				extendAction}
		}
		diff := highestbet - playerbet
		if playerchips < raise {
			return model.Actions{
				model.Action{Name: constant.Fold},
				model.Action{Name: constant.Call,
					Hints: model.Hints{
						model.Hint{
							Name: "amount", Type: "integer", Value: diff}}},
				extendAction}
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
			extendAction}
	case constant.Fold:
		return model.Actions{
			model.Action{Name: constant.Stand}}
	default:
		return model.Actions{
			model.Action{Name: constant.Sit}}
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
		MaxPlayers:   game.MaxPlayers,
		DecisionTime: game.DecisionTime,
		MaxAFKCount:  game.MaxAFKCount,
		BlindsSmall:  game.BlindsSmall,
		BlindsBig:    game.BlindsSmall,
		BuyInMin:     game.BuyInMin,
		BuyInMax:     game.BuyInMax,
		Rake:         game.Rake,
		Cap:          game.Cap,
	}
}

func (game NineK) pay(id string, chips int, action string) bool {
	index, _ := util.Get(state.Snapshot.Players, id)
	player := &state.Snapshot.Players[index]
	// not less than minimum
	if player.Bets[state.Snapshot.Turn]+chips < state.Snapshot.MinimumBet {
		return false
	}
	// not more than maximum
	if player.Bets[state.Snapshot.Turn]+chips > state.Snapshot.MaximumBet {
		return false
	}
	// cannot bet more than player's chips
	if player.Chips < chips {
		return false
	}
	util.AddScoreboardWinAmount(player.ID, -chips)
	state.Snapshot.DoActions[index] = true
	// added value to the bet in this turn
	player.Chips -= chips
	player.WinLossAmount -= chips
	player.Bets[state.Snapshot.Turn] += chips
	// broadcast to everyone that I bet
	player.Default = model.Action{Name: action}
	player.Action = model.Action{Name: action}
	player.Actions = game.Reducer(constant.Check, id)
	handler.IncreasePots(index, chips)
	// assign minimum bet
	handler.SetMinimumBet(player.Bets[state.Snapshot.Turn])
	// assign maximum bet
	handler.SetMaximumBet(util.SumPots(state.Snapshot.Pots))
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
