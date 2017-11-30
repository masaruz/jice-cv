package main

import (
	"999k_engine/constant"
	"999k_engine/gambit"
	"999k_engine/handler"
	"999k_engine/model"
	"999k_engine/state"
	"999k_engine/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
)

const port = ":3000"

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds) //Log in microsecond
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	handler.ResumeState()
	handler.Initiate(gambit.Create(os.Getenv(constant.GambitType)))
	state.GS.Gambit.Init()
	// Create queue to receiving request
	queue := make(chan func(), 1)
	// Create a worker to standby
	go func() {
		for {
			// When queue arrived
			select {
			case function := <-queue:
				log.Println("================ Start a task ================")
				// Execute task one by one
				function()
				log.Println("================ Finish a task ================")
			}
		}
	}()
	// When connection happend
	server.On(constant.Connection, func(so socketio.Socket) {
		log.Println(so.Id(), "Connect")
		// Create real enter work as connect
		// Because connect does not support message payload
		// Or retrieve player info
		so.On(constant.Enter, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "Enter", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Enter", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				channel = constant.Enter
				// Join the room
				so.Join(so.Id())
				handler.Connect(so.Id())
				handler.Enter(model.Player{
					ID:      so.Id(),
					Name:    "name",
					Picture: "picture",
				})
				handler.BroadcastGameState(so, constant.GetState, so.Id())
				state.GS.IncreaseVersion()
				log.Println(so.Id(), "Enter", "success")
				// If no seat then just return current state
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When player need server to check something
		so.On(constant.Stimulate, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "Stimulate", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Stimulate", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				log.Println("Prepare to check Start(), NextRound(), Finish()")
				// If cannot start, next and finish then it is during gameplay
				if !state.GS.Gambit.Start() &&
					!state.GS.Gambit.NextRound() &&
					!state.GS.Gambit.Finish() {
					// log.Println(so.Id(), "Stimulate", "nothing")
				} else {
					channel = constant.PushState
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, so.Id())
					// log.Println(so.Id(), "Stimulate", "success")
				}
				// If no seat then just result current state
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When player need to get game state
		so.On(constant.GetState, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "GetState", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "GetState", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				result <- handler.CreateResponse(so.Id(), constant.GetState)
			}
			return <-result
		})
		// When player call check
		so.On(constant.Check, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "Check", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Check", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if state.GS.Gambit.Check(so.Id()) {
					channel = constant.Check
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, so.Id())
					log.Println(so.Id(), "Check", "success")
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When player need to bet chips
		so.On(constant.Bet, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil || len(data.Payload.Parameters) <= 0 {
					log.Println(so.Id(), "Bet", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Bet", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				// client send amount of bet
				if state.GS.Gambit.Bet(so.Id(), data.Payload.Parameters[0].IntegerValue) {
					channel = constant.Bet
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, so.Id())
					log.Println(so.Id(), "Bet", "success")
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When player need to raise chips
		so.On(constant.Raise, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "Raise", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Raise", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				// client send amount of raise
				if state.GS.Gambit.Raise(so.Id(), data.Payload.Parameters[0].IntegerValue) {
					channel = constant.Raise
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, so.Id())
					log.Println(so.Id(), "Raise", "success")
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When player need to call chips
		so.On(constant.Call, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "Call", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Call", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if state.GS.Gambit.Call(so.Id()) {
					channel = constant.Call
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, so.Id())
					log.Println(so.Id(), "Call", "success")
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When player need to all in
		so.On(constant.AllIn, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "AllIn", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "AllIn", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if state.GS.Gambit.AllIn(so.Id()) {
					channel = constant.Raise
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, constant.Raise, so.Id())
					log.Println(so.Id(), "AllIn", "success")
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When player fold their cards
		so.On(constant.Fold, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "Fold", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Fold", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if state.GS.Gambit.Fold(so.Id()) {
					channel = constant.Fold
					state.GS.Gambit.Finish()
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, so.Id())
					log.Println(so.Id(), "Fold", "success")
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// Start table and no ending until expire
		so.On(constant.StartTable, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "StartTable", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "StartTable", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if util.CountSitting(state.GS.Players) > 1 && !handler.IsTableStart() {
					channel = constant.StartTable
					handler.StartTable()
					state.GS.Gambit.Start()
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, so.Id())
					log.Println(so.Id(), "StartTable", "success")
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When player sit down
		so.On(constant.Sit, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				if err != nil || len(data.Payload.Parameters) <= 0 {
					log.Println(so.Id(), "Sit", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Sit", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if err == nil && handler.Sit(so.Id(), data.Payload.Parameters[0].IntegerValue) {
					channel = constant.Sit
					state.GS.Gambit.Start()
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, so.Id())
					log.Println(so.Id(), "Sit", "success")
				} else {
					log.Println(so.Id(), "Sit", "Fail")
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When player stand up
		so.On(constant.Stand, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "Stand", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Stand", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if handler.Stand(so.Id()) {
					channel = constant.Stand
					state.GS.Gambit.Finish()
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, so.Id())
					log.Println(so.Id(), "Stand", "success")
				} else {
					log.Println(so.Id(), "Stand", "Fail")
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When disconnected
		so.On(constant.Disconnection, func() {
			so.Disconnect()
			log.Println(so.Id(), "Disconnect")
		})
		// When exit
		so.On(constant.Leave, func(msg string) {
			queue <- func() {
				handler.Leave(so.Id())
				state.GS.Gambit.Finish()
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, constant.Leave, so.Id())
				log.Println(so.Id(), "Leave")
			}
		})
		// When send sticker
		so.On(constant.SendSticker, func(msg string) string {
			result := make(chan string)
			queue <- func() {
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
						log.Println(so.Id(), "Send Sticker", "success")
					}
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// Extend a player action time with effect to everyone's timeline
		// And also finish round time of table
		so.On(constant.ExtendDecisionTime, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "Extend Decision Time", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Extend Decision Time", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if handler.ExtendPlayerTimeline(so.Id()) {
					channel = constant.ExtendDecisionTime
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, so.Id())
				}
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
		// When admin disband table it should be set finish table time
		so.On(constant.DisbandTable, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, err := handler.ConvertStringToRequestStruct(msg)
				// if cannot parse or client send nothing
				if err != nil {
					log.Println(so.Id(), "Disband Table", "Payload is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				if !handler.IsTableKeyValid(data.Header.Token) {
					log.Println(so.Id(), "Disband Table", "Token is invalid")
					result <- handler.CreateResponse(so.Id(), channel)
				}
				channel = constant.DisbandTable
				handler.FinishTable()
				if !handler.IsGameStart() {
					state.GS.IsTableExpired = true
				}
				handler.BroadcastGameState(so, channel, so.Id())
				log.Println(so.Id(), "Disband", "success")
				defer handler.PrepareDestroyed()
				result <- handler.CreateResponse(so.Id(), channel)
			}
			return <-result
		})
	})
	// listening for error
	server.On(constant.Error, func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	// Create router to support wildcard
	router := mux.NewRouter()
	// Handler
	router.Handle("/socket.io/", server)
	router.Handle("/socket.io", server)
	router.HandleFunc("/updateAuth", func(w http.ResponseWriter, r *http.Request) {
		// Set header to response as json format
		w.Header().Set("Content-Type", "application/json")
		var playerTableKeys []struct {
			TableKey string `json:"tablekey"`
			UserID   string `json:"userid"`
		}
		b, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(b, &playerTableKeys)
		for suid, sptk := range state.GS.PlayerTableKeys {
			for index, ptk := range playerTableKeys {
				// Same table key but not same user
				if suid != ptk.UserID && sptk == ptk.TableKey {
					playerTableKeys[index].TableKey = ""
				}
			}
		}
		// Forloop and save key into state
		for _, ptk := range playerTableKeys {
			// Hawkeye hint should delete the tablekey from this player
			if ptk.TableKey == "" {
				delete(state.GS.PlayerTableKeys, ptk.UserID)
			} else {
				state.GS.PlayerTableKeys[ptk.UserID] = ptk.TableKey
			}
		}
		// Return success to hawkeye
		resp, _ := json.Marshal(struct {
			Code      int               `json:"code"`
			Message   string            `json:"message"`
			Resources map[string]string `json:"resources"`
		}{
			Code:      200, // Success code
			Message:   "Update successfully",
			Resources: state.GS.PlayerTableKeys,
		})
		w.Write(resp)
	}).Methods("POST") // Receive only post
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
