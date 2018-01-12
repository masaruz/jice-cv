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
				util.Print("================ Start a task ================")
				// Execute task one by one
				function()
				util.Print("================ Finish a task ================")
			}
		}
	}()
	// When connection happend
	server.On(constant.Connection, func(so socketio.Socket) {
		util.Print(so.Id(), "Connect")
		// Create real enter work as connect
		// Because connect does not support message payload
		// Or retrieve player info
		so.On(constant.Enter, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Enter ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Enter", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				channel = constant.Enter
				// Join the room
				so.Join(userid)
				if handler.Enter(model.Player{
					ID:              userid,
					Name:            data.Header.DisplayName,
					AvatarSource:    data.Header.AvatarSource,
					AvatarBuiltinID: data.Header.AvatarBuiltinID,
					AvatarCustomID:  data.Header.AvatarCustomID,
					FacebookID:      data.Header.FacebookID,
				}) {
					handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					util.Print(userid, "Enter", "success")
					// If no seat then just return current state
					result <- handler.CreateResponse(userid, channel)
					return
				}
			}
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When player need server to check something
		so.On(constant.Stimulate, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Stimulate ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Stimulate", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				util.Print("Prepare to check Start(), NextRound(), Finish()")
				if state.GS.Gambit.Start() || state.GS.Gambit.NextRound() || state.GS.Gambit.Finish() {
					channel = constant.PushState
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
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
				util.Print(userid, "++++++++++++++ GetState ++++++++++++++")
				if userid == "" {
					util.Print(userid, "GetState", "Token is invalid")
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Check ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Check", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				if state.Snapshot.Gambit.Check(userid) {
					channel = constant.Check
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					util.Print(userid, "Check", "success")
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Bet ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Bet", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				// client send amount of bet
				if state.Snapshot.Gambit.Bet(userid, data.Payload.Parameters[0].IntegerValue) {
					channel = constant.Bet
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					util.Print(userid, "Bet", "success")
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Raise ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Raise", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				// client send amount of raise
				if state.Snapshot.Gambit.Raise(userid, data.Payload.Parameters[0].IntegerValue) {
					channel = constant.Raise
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					util.Print(userid, "Raise", "success")
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Call ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Call", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				if state.Snapshot.Gambit.Call(userid) {
					channel = constant.Call
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					util.Print(userid, "Call", "success")
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ AllIn ++++++++++++++")
				if userid == "" {
					util.Print(userid, "AllIn", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				if state.Snapshot.Gambit.AllIn(userid) {
					channel = constant.Raise
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, constant.Raise, userid)
					util.Print(userid, "AllIn", "success")
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Fold ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Fold", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				if state.Snapshot.Gambit.Fold(userid) {
					channel = constant.Fold
					state.GS.Gambit.Finish()
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					util.Print(userid, "Fold", "success")
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Start Table ++++++++++++++")
				if userid == "" {
					util.Print(userid, "StartTable", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				if state.Snapshot.PlayerTableKeys[userid].ClubMemberLevel != 1 {
					util.Print(userid, "Disband Table", "Not Allowed")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				if !handler.IsTableStart() {
					channel = constant.StartTable
					handler.StartTable(userid)
					if util.CountSitting(state.Snapshot.Players) > 1 {
						state.Snapshot.Gambit.Start()
					}
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					util.Print(userid, "StartTable", "success")
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Sit ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Sit", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				err := handler.Sit(userid, data.Payload.Parameters[0].IntegerValue)
				if err == nil {
					channel = constant.Sit
					state.GS.Gambit.Start()
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					util.Print(userid, "Sit", "success")
				} else {
					util.Print(userid, "Sit", "Fail")
				}
				result <- handler.CreateResponseWithCode(userid, channel, err)
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Stand ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Stand", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				if handler.Stand(userid, false) {
					channel = constant.Stand
					state.GS.Gambit.Finish()
					state.GS = util.CloneState(state.Snapshot)
					state.GS.IncreaseVersion()
					handler.BroadcastGameState(so, channel, userid)
					util.Print(userid, "Stand", "success")
				} else {
					util.Print(userid, "Stand", "Fail")
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
			util.Print(so.Id(), "Disconnect")
		})
		// When exit
		so.On(constant.Leave, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Leave ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Leave", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				channel = constant.Leave
				handler.Leave(userid)
				state.Snapshot.Gambit.Finish()
				log.Printf("Sitting: %d", util.CountSitting(state.Snapshot.Players))
				log.Printf("Visitors: %d", len(state.Snapshot.Visitors))
				state.GS = util.CloneState(state.Snapshot)
				state.GS.IncreaseVersion()
				handler.BroadcastGameState(so, channel, userid)
				util.Print(userid, "Leave", "sucess")
				result <- handler.CreateResponse(userid, channel)
				return
			}
			defer func() {
				// If no one in the room terminate itself
				util.Print("Try terminate ...")
				state.Snapshot = util.CloneState(state.GS)
				if util.CountSitting(state.Snapshot.Players) <= 0 &&
					len(state.Snapshot.Visitors) <= 0 {
					util.Print("No players then terminate")
					if state.Snapshot.Env != "dev" {
						// Delay before send signal to hawkeye that please kill this container
						go func() {
							time.Sleep(time.Millisecond * 100)
							body, err := api.Terminate()
							util.Print("Response from Terminate", string(body), err)
						}()
					}
				}
			}()
			defer util.Log()
			defer close(result)
			return <-result
		})
		// When send sticker
		so.On(constant.SendSticker, func(msg string) string {
			result := make(chan string)
			queue <- func() {
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Send Sticker ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Stand", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
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
						state.GS = util.CloneState(state.Snapshot)
						state.GS.IncreaseVersion()
						// broadcast state to everyone
						handler.BroadcastGameState(so, channel, userid)
						util.Print(userid, "Send Sticker", "success")
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Extend Decision Time ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Extend Decision Time", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				if handler.ExtendPlayerTimeline(userid) {
					channel = constant.ExtendDecisionTime
					state.GS = util.CloneState(state.Snapshot)
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
				state.Snapshot = util.CloneState(state.GS)
				channel := ""
				data, _ := handler.ConvertStringToRequestStruct(msg)
				userid := handler.GetUserIDFromToken(data.Header.Token)
				util.Print(userid, "++++++++++++++ Disable Table ++++++++++++++")
				if userid == "" {
					util.Print(userid, "Disband Table", "Token is invalid")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				handler.SetPlayerLocation(userid, data.Header.Lat, data.Header.Lon)
				if state.Snapshot.PlayerTableKeys[userid].ClubMemberLevel != 1 {
					util.Print(userid, "Disband Table", "Not Allowed")
					result <- handler.CreateResponse(userid, channel)
					return
				}
				// Set table expired less than finish table time to make sure it actually expired
				handler.FinishTable()
				// Never let player force close this table when game is started
				if !state.Snapshot.IsTableStart ||
					(!state.Snapshot.IsGameStart && !handler.IsInExtendFinishRoundTime()) {
					state.Snapshot.IsTableExpired = true
					handler.TryTerminate()
				}
				channel = constant.DisbandTable
				state.GS = util.CloneState(state.Snapshot)
				handler.BroadcastGameState(so, channel, userid)
				util.Print(userid, "Disband", "success")
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
		util.Print("error:", err)
	})
	// Create router to support wildcard
	router := mux.NewRouter()
	// Handler
	router.Handle("/socket.io/", server)
	router.Handle("/socket.io", server)
	router.HandleFunc("/updateAuth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		util.Print("++++++++++++++ Updating Auth ++++++++++++++")
		result := make(chan []byte)
		queue <- func() {
			state.Snapshot = util.CloneState(state.GS)
			handler.TryTerminate()
			playerTableKeys := []model.PlayerTableKey{}
			b, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(b, &playerTableKeys)
			for suid, sptk := range state.Snapshot.PlayerTableKeys {
				for index, ptk := range playerTableKeys {
					// Same table key but not same user
					if suid != ptk.UserID && sptk.TableKey == ptk.TableKey {
						playerTableKeys[index].TableKey = ""
					}
				}
			}
			// Forloop and save key into state
			for _, ptk := range playerTableKeys {
				// Hawkeye hint should delete the tablekey from this player
				if ptk.TableKey == "" {
					delete(state.Snapshot.PlayerTableKeys, ptk.UserID)
				} else {
					state.Snapshot.PlayerTableKeys[ptk.UserID] = ptk
				}
			}
			// Return success to hawkeye
			resp, _ := json.Marshal(struct {
				Code      int                             `json:"code"`
				Message   string                          `json:"message"`
				Resources map[string]model.PlayerTableKey `json:"resources"`
			}{
				Code:      200, // Success code
				Message:   "Update successfully",
				Resources: state.Snapshot.PlayerTableKeys,
			})
			state.GS = util.CloneState(state.Snapshot)
			result <- resp
		}
		w.Write(<-result)
	}).Methods("POST") // Receive only post
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		str := ""
		if state.Snapshot.Env == "dev" {
			for _, pair := range os.Environ() {
				str += fmt.Sprintf("%s ", pair)
			}
		}
		fmt.Fprintf(w, "All envs are here: %s", str)
	})
	http.Handle("/", router)
	util.Print(fmt.Sprintf("Serving at localhost%s", port))

	log.Fatal(http.ListenAndServe(port, nil))
}
