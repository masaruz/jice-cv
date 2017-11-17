package main

import (
	"999k_engine/constant"
	"999k_engine/gambit"
	"999k_engine/handler"
	"999k_engine/state"
	"999k_engine/util"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/googollee/go-socket.io"
)

const port = ":3000"

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.SetAllowRequest(func(r *http.Request) error {
		log.Println(r.Header)
		return nil
	})
	handler.SetGambit(gambit.Create(os.Getenv(constant.GambitType)))
	state.GS.Gambit.Init()
	// when connection happend
	server.On(constant.Connection, func(so socketio.Socket) {
		// join the room
		so.Join(so.Id())
		handler.Connect(so.Id())
		handler.BroadcastGameState(so, constant.GetState, so.Id())
		state.GS.IncreaseVersion()
		log.Println(so.Id(), "Connect")
		// when player need server to check something
		so.On(constant.Stimulate, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			// if cannot start, next and finish then it is during gameplay
			if !state.GS.Gambit.Start() &&
				!state.GS.Gambit.NextRound() &&
				!state.GS.Gambit.Finish() {
				log.Println(so.Id(), "Stimulate", "Nothing")
			} else {
				channel = constant.PushState
				state.GS.IncreaseVersion()
				log.Println(so.Id(), "Stimulate", "Success")
			}
			handler.FinishProcess()
			// if no seat then just return current state
			return handler.CreateResponse(so.Id(), channel)
		})
		// when player need to get game state
		so.On(constant.GetState, func(msg string) string {
			return handler.CreateResponse(so.Id(), constant.GetState)
		})
		// when player call check
		so.On(constant.Check, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			if state.GS.Gambit.Check(so.Id()) {
				channel = constant.Check
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Check", "Success")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// when player need to bet chips
		so.On(constant.Bet, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			data, err := handler.ConvertStringToRequestStruct(msg)
			// if cannot parse or client send nothing
			if err != nil || len(data.Payload.Parameters) <= 0 {
				return handler.CreateResponse(so.Id(), channel)
			}
			// client send amount of bet
			if state.GS.Gambit.Bet(so.Id(), data.Payload.Parameters[0].IntegerValue) {
				channel = constant.Bet
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Bet", "Success")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// when player need to raise chips
		so.On(constant.Raise, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			data, err := handler.ConvertStringToRequestStruct(msg)
			// if cannot parse or client send nothing
			if err != nil || len(data.Payload.Parameters) <= 0 {
				return handler.CreateResponse(so.Id(), channel)
			}
			// client send amount of raise
			if state.GS.Gambit.Raise(so.Id(), data.Payload.Parameters[0].IntegerValue) {
				channel = constant.Raise
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Raise", "Success")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// when player need to call chips
		so.On(constant.Call, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			if state.GS.Gambit.Call(so.Id()) {
				channel = constant.Call
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Call", "Success")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// when player need to all in
		so.On(constant.AllIn, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			if state.GS.Gambit.AllIn(so.Id()) {
				channel = constant.Raise
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, constant.Raise, so.Id())
				log.Println(so.Id(), "AllIn", "Success")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// when player fold their cards
		so.On(constant.Fold, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			if state.GS.Gambit.Fold(so.Id()) {
				channel = constant.Fold
				state.GS.Gambit.Finish()
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Fold", "Success")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// start table and no ending until expire
		so.On(constant.StartTable, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			if util.CountSitting(state.GS.Players) > 1 && !handler.IsTableStart() {
				channel = constant.StartTable
				handler.StartTable()
				state.GS.Gambit.Start()
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "StartTable", "Success")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// when player sit down
		so.On(constant.Sit, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			data, err := handler.ConvertStringToRequestStruct(msg)
			if err == nil && handler.Sit(so.Id(), data.Payload.Parameters[0].IntegerValue) {
				channel = constant.Sit
				state.GS.Gambit.Start()
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Sit", "Success")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// when player stand up
		so.On(constant.Stand, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			if handler.Stand(so.Id()) {
				channel = constant.Stand
				state.GS.Gambit.Finish()
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Stand", "Success")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// when disconnected
		so.On(constant.Disconnection, func() {
			log.Println(so.Id(), "Disconnect")
		})
		// when exit
		so.On(constant.Leave, func(msg string) {
			handler.WaitQueue()
			handler.StartProcess()
			handler.Leave(so.Id())
			state.GS.Gambit.Finish()
			state.GS.IncreaseVersion()
			handler.BroadcastGameState(so, constant.Leave, so.Id())
			handler.FinishProcess()
			log.Println(so.Id(), "Leave")
		})
		// when send sticker
		so.On(constant.SendSticker, func(msg string) string {
			channel := ""
			data, err := handler.ConvertStringToRequestStruct(msg)
			// if cannot parse or client send nothing
			if err == nil && len(data.Payload.Parameters) == 2 {
				param1 := data.Payload.Parameters[0]
				param2 := data.Payload.Parameters[1]
				stickerid, targetslot := "", -1
				// handler prevent swap of parameters
				switch param1.Name {
				// in case of stickerid is param1
				case "stickerid":
					stickerid = param1.StringValue
					targetslot = param2.IntegerValue
					// in case of stickerid is param2
				case "targetslot":
					stickerid = param2.StringValue
					targetslot = param1.IntegerValue
				default:
				}
				if stickerid != "" && targetslot != -1 {
					channel = constant.SendSticker
					// set sticker state in player
					handler.SendSticker(stickerid, so.Id(), targetslot)
					state.GS.IncreaseVersion()
					// broadcast state to everyone
					handler.BroadcastGameState(so, channel, so.Id())
					log.Println(so.Id(), "Send Sticker", "Success", data.Payload)
				}
			}
			return handler.CreateResponse(so.Id(), channel)
		})
		so.On(constant.ExtendDecisionTime, func(msg string) string {
			channel := ""
			if handler.ExtendPlayerTimeline(so.Id()) {
				channel = constant.ExtendDecisionTime
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Extend", "Success")
				return handler.CreateResponse(so.Id(), channel)
			}
			return handler.CreateResponse(so.Id(), channel)
		})
	})
	// listening for error
	server.On(constant.Error, func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		str := ""
		for _, pair := range os.Environ() {
			str += fmt.Sprintf("%s ", pair)
		}
		fmt.Fprintf(w, "All envs are here: %s", str)
	})
	log.Println(fmt.Sprintf("Serving at localhost%s", port))

	log.Fatal(http.ListenAndServe(port, nil))
}
