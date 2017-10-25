package main

import (
	"999k_engine/constant"
	"999k_engine/game"
	"999k_engine/handler"
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
)

const port = ":3000"

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	decisionTime := int64(20)
	ninek := game.NineK{
		MaxPlayers:   6,
		DecisionTime: decisionTime,
		MinimumBet:   10}
	handler.SetGambit(ninek)
	state.GS.Gambit.Init()
	// when connection happend
	server.On(constant.Connection, func(so socketio.Socket) {
		// join the room
		so.Join(so.Id())
		handler.Connect(so.Id())
		handler.BroadcastGameState(so, constant.GetState, so.Id())
		state.GS.IncreaseVersion()
		fmt.Println(so.Id(), "Connect")
		// when player need server to check something
		so.On(constant.Stimulate, func(msg string) string {
			// if cannot start, next and finish then it is during gameplay
			if !state.GS.Gambit.Start() &&
				!state.GS.Gambit.NextRound() &&
				!state.GS.Gambit.Finish() {
				fmt.Println(so.Id(), "Stimulate", "Nothing", msg)
				return handler.CreateResponse(so.Id(), "")
			}
			fmt.Println(so.Id(), "Stimulate", "Success", msg)
			state.GS.IncreaseVersion()
			// if no seat then just return current state
			return handler.CreateResponse(so.Id(), constant.PushState)
		})
		// when player need to get game state
		so.On(constant.GetState, func(msg string) string {
			return handler.CreateResponse(so.Id(), constant.GetState)
		})
		// when player call check
		so.On(constant.Check, func(msg string) string {
			if !state.GS.Gambit.Check(so.Id()) {
				fmt.Println(so.Id(), "Check", "Nothing", msg)
				// if no seat then just return current state
				return handler.CreateResponse(so.Id(), "")
			}
			fmt.Println(so.Id(), "Check", "Success", msg)
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.Check, so.Id())
			return handler.CreateResponse(so.Id(), constant.Check)
		})
		// when player need to bet chips
		so.On(constant.Bet, func(msg string) string {
			data := &state.Req{}
			err := json.Unmarshal([]byte(msg), data)
			// if cannot parse or client send nothing
			if err != nil || len(data.Payload.Parameters) <= 0 {
				return handler.CreateResponse(so.Id(), "")
			}
			// client send amount of bet
			if !state.GS.Gambit.Bet(so.Id(), data.Payload.Parameters[0].ValueInteger) {
				fmt.Println(so.Id(), "Bet", "Nothing", msg)
				// if no seat then just return current state
				return handler.CreateResponse(so.Id(), "")
			}
			fmt.Println(so.Id(), "Bet", "Success", msg)
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.Bet, so.Id())
			return handler.CreateResponse(so.Id(), constant.Bet)
		})
		// when player need to raise chips
		so.On(constant.Raise, func(msg string) string {
			data := &state.Req{}
			err := json.Unmarshal([]byte(msg), data)
			// if cannot parse or client send nothing
			if err != nil || len(data.Payload.Parameters) <= 0 {
				return handler.CreateResponse(so.Id(), "")
			}
			// client send amount of raise
			if !state.GS.Gambit.Raise(so.Id(), data.Payload.Parameters[0].ValueInteger) {
				fmt.Println(so.Id(), "Raise", "Nothing", msg)
				// if no seat then just return current state
				return handler.CreateResponse(so.Id(), "")
			}
			fmt.Println(so.Id(), "Raise", "Success", msg)
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.Raise, so.Id())
			return handler.CreateResponse(so.Id(), constant.Raise)
		})
		// when player need to call chips
		so.On(constant.Call, func(msg string) string {
			if !state.GS.Gambit.Call(so.Id()) {
				fmt.Println(so.Id(), "Call", "Nothing", msg)
				// if no seat then just return current state
				return handler.CreateResponse(so.Id(), "")
			}
			fmt.Println(so.Id(), "Call", "Success", msg)
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.Call, so.Id())
			return handler.CreateResponse(so.Id(), constant.Call)
		})
		// when player need to all in
		so.On(constant.AllIn, func(msg string) string {
			if !state.GS.Gambit.AllIn(so.Id()) {
				fmt.Println(so.Id(), "AllIn", "Nothing", msg)
				// if no seat then just return current state
				return handler.CreateResponse(so.Id(), "")
			}
			fmt.Println(so.Id(), "AllIn", "Success", msg)
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.Raise, so.Id())
			return handler.CreateResponse(so.Id(), constant.Raise)
		})
		// when player fold their cards
		so.On(constant.Fold, func(msg string) string {
			if !state.GS.Gambit.Fold(so.Id()) {
				fmt.Println(so.Id(), "Fold", "Nothing", msg)
				// if no seat then just return current state
				return handler.CreateResponse(so.Id(), "")
			}
			fmt.Println(so.Id(), "Fold", "Success", msg)
			state.GS.Gambit.Finish()
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.Fold, so.Id())
			return handler.CreateResponse(so.Id(), constant.Fold)
		})
		// start table and no ending until expire
		so.On(constant.StartTable, func(msg string) string {
			if util.CountSitting(state.GS.Players) <= 1 {
				fmt.Println(so.Id(), "StartTable", "Nothing", msg)
				return handler.CreateResponse(so.Id(), "")
			}
			fmt.Println(so.Id(), "StartTable", "Success", msg)
			handler.StartTable()
			state.GS.Gambit.Start()
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.StartTable, so.Id())
			return handler.CreateResponse(so.Id(), constant.StartTable)
		})
		// when player sit down
		so.On(constant.Sit, func(msg string) string {
			data := &state.Req{}
			err := json.Unmarshal([]byte(msg), data)
			// if cannot parse or client send nothing
			if err != nil || len(data.Payload.Parameters) <= 0 {
				if !handler.AutoSit(so.Id()) {
					fmt.Println(so.Id(), "Sit", "Nothing", msg)
					return handler.CreateResponse(so.Id(), "")
				}
			} else if !handler.Sit(so.Id(), data.Payload.Parameters[0].ValueInteger) {
				// client send seat position
				fmt.Println(so.Id(), "Sit", "Nothing", msg)
				// if no seat then just return current state
				return handler.CreateResponse(so.Id(), "")
			}
			fmt.Println(so.Id(), "Sit", "Success", msg)
			state.GS.Gambit.Start()
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.Sit, so.Id())
			return handler.CreateResponse(so.Id(), constant.Sit)
		})
		// when player stand up
		so.On(constant.Stand, func(msg string) string {
			if !handler.Stand(so.Id()) {
				fmt.Println(so.Id(), "Stand", "Nothing", msg)
				return handler.CreateResponse(so.Id(), "")
			}
			fmt.Println(so.Id(), "Stand", "Success", msg)
			state.GS.Gambit.Finish()
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.Stand, so.Id())
			return handler.CreateResponse(so.Id(), constant.Stand)
		})
		// when disconnected
		so.On(constant.Disconnection, func() {
			fmt.Println(so.Id(), "Disconnect")
			handler.Disconnect(so.Id())
			state.GS.Gambit.Finish()
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.Disconnection, so.Id())
		})
	})
	// listening for error
	server.On(constant.Error, func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	// http.Handle("/", http.FileServer(http.Dir("../asset")))
	http.Handle("/socket.io/", server)

	log.Println(fmt.Sprintf("Serving at localhost%s", port))

	log.Fatal(http.ListenAndServe(port, nil))
}
