package main

import (
	"999k_engine/api"
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
	"time"

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
	handler.RestoreStateData()
	handler.Initiate(gambit.Create(os.Getenv(constant.GambitType)))
	state.GS.Gambit.Init()
	// Create queue to receiving request
	queue := make(chan func())
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
			log.Println("++++ Enter ++++")
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Enter", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				channel = constant.Enter
				// Join the room
				so.Join(userid)
				handler.Enter(model.Player{
					ID:      userid,
					Name:    data.Header.DisplayName,
					Picture: "picture",
				})
				handler.BroadcastGameState(so, channel, userid)
				state.GS.IncreaseVersion()
				log.Println(userid, "Enter", "success")
				// If no seat then just return current state
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player need server to check something
		so.On(constant.Stimulate, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Stimulate", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				log.Println("Prepare to check Start(), NextRound(), Finish()")
				// If cannot start, next and finish then it is during gameplay
				if !state.GS.Gambit.Start() &&
					!state.GS.Gambit.NextRound() &&
					!state.GS.Gambit.Finish() {
					// log.Println(userid, "Stimulate", "nothing")
				} else {
					channel = constant.PushState
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					// log.Println(userid, "Stimulate", "success")
				}
				// If no seat then just result current state
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player need to get game state
		so.On(constant.GetState, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "GetState", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				result <- handler.CreateResponse(userid, constant.GetState)
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player call check
		so.On(constant.Check, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Check", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				if state.GS.Gambit.Check(userid) {
					channel = constant.Check
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					log.Println(userid, "Check", "success")
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player need to bet chips
		so.On(constant.Bet, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Bet", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				// client send amount of bet
				if state.GS.Gambit.Bet(userid, data.Payload.Parameters[0].IntegerValue) {
					channel = constant.Bet
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					log.Println(userid, "Bet", "success")
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player need to raise chips
		so.On(constant.Raise, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Raise", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				// client send amount of raise
				if state.GS.Gambit.Raise(userid, data.Payload.Parameters[0].IntegerValue) {
					channel = constant.Raise
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					log.Println(userid, "Raise", "success")
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player need to call chips
		so.On(constant.Call, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Call", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				if state.GS.Gambit.Call(userid) {
					channel = constant.Call
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					log.Println(userid, "Call", "success")
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player need to all in
		so.On(constant.AllIn, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "AllIn", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				if state.GS.Gambit.AllIn(userid) {
					channel = constant.Raise
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, constant.Raise, userid)
					log.Println(userid, "AllIn", "success")
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player fold their cards
		so.On(constant.Fold, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Fold", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				if state.GS.Gambit.Fold(userid) {
					channel = constant.Fold
					state.GS.Gambit.Finish()
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					log.Println(userid, "Fold", "success")
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// Start table and no ending until expire
		so.On(constant.StartTable, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "StartTable", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				if util.CountSitting(state.GS.Players) > 1 && !handler.IsTableStart() {
					channel = constant.StartTable
					handler.StartTable()
					state.GS.Gambit.Start()
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					log.Println(userid, "StartTable", "success")
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player sit down
		so.On(constant.Sit, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Sit", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				if handler.Sit(userid, data.Payload.Parameters[0].IntegerValue) {
					channel = constant.Sit
					state.GS.Gambit.Start()
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					log.Println(userid, "Sit", "success")
				} else {
					log.Println(userid, "Sit", "Fail")
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player stand up
		so.On(constant.Stand, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Stand", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				if handler.Stand(userid, false) {
					channel = constant.Stand
					state.GS.Gambit.Finish()
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					log.Println(userid, "Stand", "success")
				} else {
					log.Println(userid, "Stand", "Fail")
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When disconnected
		so.On(constant.Disconnection, func() {
			so.Disconnect()
			log.Println(so.Id(), "Disconnect")
		})
		// When exit
		so.On(constant.Leave, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Leave", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				channel = constant.Leave
				handler.Leave(userid)
				state.GS.Gambit.Finish()
				log.Printf("Sitting: %d", util.CountSitting(state.GS.Players))
				log.Printf("Visitors: %d", len(state.GS.Visitors))
				// If no one in the room terminate itself
				log.Println("Try terminate ...")
				if util.CountSitting(state.GS.Players) <= 0 &&
					len(state.GS.Visitors) <= 0 {
					log.Println("No players then terminate")
					if os.Getenv("env") != "dev" {
						// Delay 5 second before send signal to hawkeye that please kill this container
						go func() {
							time.Sleep(time.Second * 3)
							body, err := api.Terminate()
							log.Println("Response from Terminate", string(body), err)
						}()
					}
				}
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, userid)
				log.Println(userid, "Leave", "sucess")
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When send sticker
		so.On(constant.SendSticker, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Stand", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				// if cannot parse or client send nothing
				if len(data.Payload.Parameters) == 2 {
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
						handler.SendSticker(stickerid, userid, targetslot)
						state.GS.IncreaseVersion()
						// broadcast state to everyone
						handler.BroadcastGameState(so, channel, userid)
						log.Println(userid, "Send Sticker", "success")
					}
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// Extend a player action time with effect to everyone's timeline
		// And also finish round time of table
		so.On(constant.ExtendDecisionTime, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Extend Decision Time", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				if handler.ExtendPlayerTimeline(userid) {
					channel = constant.ExtendDecisionTime
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
				}
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When admin disband table it should be set finish table time
		so.On(constant.DisbandTable, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				if userid == "" {
					log.Println(userid, "Disband Table", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				// Set table expired equal 0 to make sure it actually expired
				handler.FinishTable()
				// Never let player force close this table when game is started
				if !handler.IsGameStart() {
					state.GS.IsTableExpired = true
					handler.TryTerminate()
				}
				channel = constant.DisbandTable
				handler.BroadcastGameState(so, channel, userid)
				log.Println(userid, "Disband", "success")
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer util.Log()
			defer close(result)
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
		if os.Getenv("env") == "dev" {
			for _, pair := range os.Environ() {
				str += fmt.Sprintf("%s ", pair)
			}
		}
		fmt.Fprintf(w, "All envs are here: %s", str)
	})
	http.Handle("/", router)
	log.Println(fmt.Sprintf("Serving at localhost%s", port))

	log.Fatal(http.ListenAndServe(port, nil))
}
