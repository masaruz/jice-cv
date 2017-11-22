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
	"time"

	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
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
	handler.Initiate(gambit.Create(os.Getenv(constant.GambitType)))
	state.GS.Gambit.Init()
	// When connection happend
	server.On(constant.Connection, func(so socketio.Socket) {
		// Join the room
		so.Join(so.Id())
		handler.Connect(so.Id())
		handler.BroadcastGameState(so, constant.GetState, so.Id())
		state.GS.IncreaseVersion()
		log.Println(so.Id(), "Connect")
		// When player need server to check something
		so.On(constant.Stimulate, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			// If cannot start, next and finish then it is during gameplay
			if !state.GS.Gambit.Start() &&
				!state.GS.Gambit.NextRound() &&
				!state.GS.Gambit.Finish() {
				log.Println(so.Id(), "Stimulate", "Nothing")
			} else {
				channel = constant.PushState
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Stimulate", "Success")
			}
			handler.FinishProcess()
			// If no seat then just return current state
			return handler.CreateResponse(so.Id(), channel)
		})
		// When player need to get game state
		so.On(constant.GetState, func(msg string) string {
			return handler.CreateResponse(so.Id(), constant.GetState)
		})
		// When player call check
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
		// When player need to bet chips
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
		// When player need to raise chips
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
		// When player need to call chips
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
		// When player need to all in
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
		// When player fold their cards
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
		// Start table and no ending until expire
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
		// When player sit down
		so.On(constant.Sit, func(msg string) string {
			log.Println(so.Id(), "Sit", "Trying")
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
			} else {
				log.Println(so.Id(), "Sit", "Fail")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// When player stand up
		so.On(constant.Stand, func(msg string) string {
			log.Println(so.Id(), "Stand", "Trying")
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			if handler.Stand(so.Id()) {
				channel = constant.Stand
				state.GS.Gambit.Finish()
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Stand", "Success")
			} else {
				log.Println(so.Id(), "Stand", "Fail")
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// When disconnected
		so.On(constant.Disconnection, func() {
			log.Println(so.Id(), "Disconnect")
		})
		// When exit
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
		// When send sticker
		so.On(constant.SendSticker, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
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
					log.Println(so.Id(), "Send Sticker", "Success")
				}
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// Extend a player action time with effect to everyone's timeline
		// And also finish round time of table
		so.On(constant.ExtendDecisionTime, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			if handler.ExtendPlayerTimeline(so.Id()) {
				channel = constant.ExtendDecisionTime
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, so.Id())
			}
			handler.FinishProcess()
			return handler.CreateResponse(so.Id(), channel)
		})
		// When admin disband table it should be set finish table time
		so.On(constant.DisbandTable, func(msg string) string {
			handler.WaitQueue()
			handler.StartProcess()
			channel := ""
			channel = constant.DisbandTable
			handler.FinishTable()
			if !handler.IsGameStart() {
				state.GS.IsTableExpired = true
			}
			handler.FinishProcess()
			handler.BroadcastGameState(so, channel, so.Id())
			go func() {
				// fmt.Printf("caught sig: %+v", sig)
				fmt.Println("Wait for 2 second to finish processing")
				time.Sleep(2 * time.Second)
				os.Exit(0)
			}()
			log.Println(so.Id(), "Disband", "Success")
			return handler.CreateResponse(so.Id(), channel)
		})
	})
	// listening for error
	server.On(constant.Error, func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	// Create router to support wildcard
	router := mux.NewRouter()
	// Handler
	router.Handle("/{tableid}/socket.io/", server)
	router.Handle("/{tableid}/socket.io", server)
	router.Handle("/socket.io/", server)
	router.Handle("/socket.io", server)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		str := ""
		for _, pair := range os.Environ() {
			str += fmt.Sprintf("%s ", pair)
		}
		fmt.Fprintf(w, "All envs are here: %s", str)
	})
	http.Handle("/", router)
	log.Println(fmt.Sprintf("Serving at localhost%s", port))

	log.Fatal(http.ListenAndServe(port, nil))
}
