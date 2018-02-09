package handler

import (
	"999k_engine/api"
	"999k_engine/state"
	"999k_engine/util"
	"time"
)

// TryTerminateWhenNoPlayers if no player in table then terminate
func TryTerminateWhenNoPlayers() {
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
}
