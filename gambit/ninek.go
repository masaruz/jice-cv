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
	"log"
	"os"
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

// Payload data accessed by continue
type Payload struct {
	ID    string
	Chips int
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
	if handler.IsTableStart() &&
		!handler.IsGameStart() &&
		!handler.IsInExtendFinishRoundTime() {
		// filter players who are not ready to play
		for index := range state.GS.Players {
			player := &state.GS.Players[index]
			// If this position is empty seat continue
			if player.ID == "" {
				continue
			}
			// If player has no chip
			// TODO if player has not enough chip then all-in their chips
			if player.Chips < game.GetSettings().BlindsSmall {
				// Validate with other server when is not in dev
				if os.Getenv("env") != "test" {
					// Need request to server for buyin
					body, _ := api.BuyIn(player.ID, game.GetSettings().BuyInMin)
					resp := &api.PlayerResponse{}
					json.Unmarshal(body, resp)
					// If this player request buy in success
					if resp.Error.StatusCode == 0 {
						log.Println("Buy-in success")
						// Assign how much they buy-in
						player.Chips = game.GetSettings().BuyInMin
					} else {
						log.Fatal(resp)
					}
				} else {
					player.Chips = game.GetSettings().BuyInMin
				}
			}
			// If player has minimum chip for able to play
			if player.Chips >= game.GetSettings().BlindsSmall &&
				state.GS.AFKCounts[index] < game.MaxAFKCount {
				continue
			}
			handler.Stand(player.ID)
		}
		// After pass through the critiria
		// if there are more than 2 players are sitting
		if util.CountSitting(state.GS.Players) >= 2 {
			// everyone is assumed afk
			state.GS.DoActions = make([]bool, game.MaxPlayers)
			state.GS.Rakes = make(map[string]float64)
			state.GS.Pots = make([]int, game.MaxPlayers)
			// request to start game
			_, err := api.StartGame()
			if err != nil {
				panic(err)
			}
			// set players to be ready
			handler.MakePlayersReady()
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
			handler.SetPlayersRake(game.Rake, game.Cap*float64(game.BlindsBig))
			return true
		}
	}
	return false
}

// NextRound game after round by round
func (game NineK) NextRound() bool {
	// Assume every player must be default
	handler.OverwriteActionToBehindPlayers()
	// Game must be start
	// All players must not has 3 cards
	// All bets must be equal
	// Now must be more than finish round time
	if handler.IsGameStart() && !handler.IsFullHand(3) &&
		handler.BetsEqual() && handler.IsEndRound() &&
		util.CountPlayerNotFold(state.GS.Players) > 1 {
		// Initialize values
		handler.SetMinimumBet(game.BlindsBig)
		handler.SetOtherActions("", constant.Check)
		handler.SetOtherDefaultAction("", constant.Check)
		handler.CreateTimeLine(game.DecisionTime)
		handler.PlayersInvestToPots(0)
		handler.Deal(1, game.MaxPlayers)
		handler.IncreaseTurn()
		// no one is assumed afk
		state.GS.DoActions = make([]bool, game.MaxPlayers)
		return true
	}
	return false
}

// Finish game
func (game NineK) Finish() bool {
	handler.OverwriteActionToBehindPlayers()
	// no others to play with or all players have 3 cards but bet is not equal
	if handler.IsGameStart() && ((util.CountPlayerNotFold(state.GS.Players) <= 1) ||
		// if has 3 cards bet equal
		(handler.IsFullHand(3) && handler.BetsEqual() && handler.IsEndRound())) {
		// calculate afk players
		for index, doaction := range state.GS.DoActions {
			// Skip empty players
			if !util.IsPlayingAndNotFoldAndNotAllIn(state.GS.Players[index]) {
				continue
			}
			// If this player done something in this round
			if doaction {
				state.GS.AFKCounts[index] = 0
			} else {
				// If this player never done action in this game
				state.GS.AFKCounts[index]++
			}
		}
		// Extend more time for client to play animation after game finished
		handler.ExtendFinishRoundTime()
		// Find winner and added their rewards
		hscore := -1
		hbonus := -1
		pos := -1
		// Evaluate score from everyone's hand
		for i := 0; i < len(state.GS.Players); i++ {
			for index, player := range state.GS.Players {
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
				for poti, pot := range state.GS.Pots {
					if pot == 0 {
						continue
					}
					playerbet := pot
					winnerbet := state.GS.Pots[pos]
					earnedbet := 0
					if winnerbet > playerbet {
						// If winner has higher bet
						earnedbet = playerbet
					} else {
						// If winner has lower bet
						earnedbet = winnerbet
					}
					if earnedbet != 0 {
						state.GS.Players[pos].Chips += earnedbet
						state.GS.Players[pos].WinLossAmount += earnedbet
						earnedplayers := util.CountPlayerAlreadyEarned(state.GS.Players)
						if util.CountPlayerNotFold(state.GS.Players)-earnedplayers > 1 || earnedplayers == 0 {
							state.GS.Players[pos].IsWinner = true
						}
					}
					// If not caller
					if poti != pos {
						handler.BurnBet(poti, earnedbet)
					}
				}
				state.GS.Players[pos].IsEarned = true
				handler.BurnBet(pos, util.SumBet(state.GS.Players[pos]))
				hscore = -1
				hbonus = -1
				pos = -1
			}
		}
		// Update all players' buy-in amount
		// api.SaveSettlements()
		// Increase gameindex for backend process ex. realtime-data, analytic
		handler.IncreaseGameIndex()
		// Revert minimum bet
		handler.SetMinimumBet(game.BlindsBig)
		handler.FlushPlayers()
		state.GS.Turn = 0
		state.GS.IsGameStart = false
		// Check if current time is more than finish table time
		if time.Now().Unix() >= state.GS.FinishTableTime {
			// For force client to leave
			state.GS.IsTableExpired = true
			state.GS.IsTableStart = false
			defer handler.PrepareDestroyed()
		}
		return true
	}
	return false
}

// Check is doing nothing but only shift the timeline
func (game NineK) Check(id string) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.GS.Players, id)
	player := &state.GS.Players[index]
	// Cannot check if player has less bet than highest
	if player.Bets[state.GS.Turn] <
		util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players) {
		return false
	}
	state.GS.DoActions[index] = true
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
	index, _ := util.Get(state.GS.Players, id)
	player := &state.GS.Players[index]
	// not less than minimum
	if player.Bets[state.GS.Turn]+chips < state.GS.MinimumBet {
		return false
	}
	// not more than maximum
	if player.Bets[state.GS.Turn]+chips > state.GS.MaximumBet {
		return false
	}
	// cannot bet more than player's chips
	if player.Chips < chips {
		return false
	}
	state.GS.DoActions[index] = true
	// added value to the bet in this turn
	player.Chips -= chips
	player.WinLossAmount -= chips
	player.Bets[state.GS.Turn] += chips
	// broadcast to everyone that I bet
	player.Default = model.Action{Name: constant.Bet}
	player.Action = model.Action{Name: constant.Bet}
	player.Actions = game.Reducer(constant.Check, id)
	handler.IncreasePots(index, chips)
	// assign minimum bet
	handler.SetMinimumBet(player.Bets[state.GS.Turn])
	// assign maximum bet
	handler.SetMaximumBet(util.SumPots(state.GS.Pots))
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

// Raise is raising bet to the target
func (game NineK) Raise(id string, chips int) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.GS.Players, id)
	// not less than minimum
	if state.GS.Players[index].Bets[state.GS.Turn]+chips <= state.GS.MinimumBet {
		return false
	}
	return game.Bet(id, chips)
}

// Call is raising bet to the highest bet
func (game NineK) Call(id string) bool {
	if !handler.IsPlayerTurn(id) {
		return false
	}
	index, _ := util.Get(state.GS.Players, id)
	player := &state.GS.Players[index]
	chips := util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players) -
		player.Bets[state.GS.Turn]
	// cannot call more than player's chips
	if player.Chips < chips || chips == 0 {
		return false
	}
	state.GS.DoActions[index] = true
	player.Chips -= chips
	player.WinLossAmount -= chips
	player.Bets[state.GS.Turn] += chips
	player.Default = model.Action{Name: constant.Call}
	player.Action = model.Action{Name: constant.Call}
	player.Actions = game.Reducer(constant.Check, id)
	handler.IncreasePots(index, chips)
	// set action of everyone
	handler.OverwriteActionToBehindPlayers()
	handler.SetMaximumBet(util.SumPots(state.GS.Pots))
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
	index, _ := util.Get(state.GS.Players, id)
	player := &state.GS.Players[index]
	chips := player.Chips
	// not more than maximum
	if player.Bets[state.GS.Turn]+chips > state.GS.MaximumBet {
		return false
	}
	state.GS.DoActions[index] = true
	player.Bets[state.GS.Turn] += chips
	player.WinLossAmount -= chips
	player.Chips = 0
	player.Default = model.Action{Name: constant.AllIn}
	player.Action = model.Action{Name: constant.AllIn}
	player.Actions = game.Reducer(constant.Check, id)
	handler.IncreasePots(index, chips)
	handler.SetMaximumBet(util.SumPots(state.GS.Pots))
	// set action of everyone
	handler.OverwriteActionToBehindPlayers()
	// others automatic set to fold as default
	handler.SetOtherDefaultAction(id, constant.Fold)
	// others need to know what to do next
	handler.SetOtherActions(id, constant.Bet)
	diff := time.Now().Unix() - player.DeadLine
	handler.ShortenTimeline(diff)
	// duration extend the timeline
	if player.Bets[state.GS.Turn] >= util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players) {
		handler.SetMinimumBet(player.Bets[state.GS.Turn])
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
	index, _ := util.Get(state.GS.Players, id)
	player := &state.GS.Players[index]
	state.GS.DoActions[index] = true
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
		_, player := util.Get(state.GS.Players, id)
		if player.Chips == 0 {
			return model.Actions{
				model.Action{Name: constant.Fold},
				model.Action{Name: constant.Check},
				extendAction}
		}
		// maximum will be player's chips if not enough
		maximum := 0
		if state.GS.MaximumBet > player.Chips {
			maximum = player.Chips
		} else {
			maximum = state.GS.MaximumBet
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
						Name: "amount", Type: "integer", Value: state.GS.MinimumBet},
					model.Hint{
						Name: "amount_max", Type: "integer", Value: maximum}}},
			extendAction}
	case constant.Bet:
		_, player := util.Get(state.GS.Players, id)
		playerchips := player.Chips + player.Bets[state.GS.Turn]
		// highest bet in that turn
		highestbet := util.GetHighestBetInTurn(state.GS.Turn, state.GS.Players)
		playerbet := player.Bets[state.GS.Turn]
		// raise must be highest * 2
		raise := highestbet * 2
		// all sum bets
		pots := util.SumPots(state.GS.Pots)
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
		if state.GS.MaximumBet > playerchips {
			maximum = playerchips
		} else {
			maximum = state.GS.MaximumBet
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
