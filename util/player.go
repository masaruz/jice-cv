package util

import (
	"999k_engine/constant"
	"999k_engine/model"
	"time"
)

// Remove by remove element from array
func Remove(slice model.Players, id string) model.Players {
	target := -1
	for index, player := range slice {
		if id == player.ID {
			target = index
		}
	}
	if target == -1 {
		return slice
	}
	return append(slice[:target], slice[target+1:]...)
}

// Kick by marking element to empty in slice
func Kick(slice model.Players, id string) model.Players {
	for index, player := range slice {
		if id == player.ID {
			if player.Type == constant.Normal {
				slice[index] = model.Player{Slot: player.Slot}
			} else {
				slice[index] = model.Player{Slot: player.Slot, Type: player.Type}
			}
		}
	}
	return slice
}

// Add element from array
func Add(slice model.Players, player model.Player) model.Players {
	for _, other := range slice {
		if other.ID == player.ID {
			slice = Remove(slice, player.ID)
			break
		}
	}
	return append(slice, player)
}

// Get element from array by key
func Get(slice model.Players, id string) (int, model.Player) {
	for index, player := range slice {
		if id == player.ID {
			return index, player
		}
	}
	return -1, model.Player{}
}

// CountPlayerAlreadyEarned who has eared
func CountPlayerAlreadyEarned(players model.Players) int {
	count := 0
	for _, player := range players {
		if player.IsEarned {
			count++
		}
	}
	return count
}

// CountPlayerNotFoldAndNotAllIn who has right to play
func CountPlayerNotFoldAndNotAllIn(players model.Players) int {
	count := 0
	for _, player := range players {
		if IsPlayingAndNotFoldAndNotAllIn(player) {
			count++
		}
	}
	return count
}

// CountPlayerNotFold count who is not fold
func CountPlayerNotFold(players model.Players) int {
	count := 0
	for _, player := range players {
		if IsPlayingAndNotFold(player) {
			count++
		}
	}
	return count
}

// CountSitting who is actually sit
func CountSitting(players model.Players) int {
	count := 0
	for _, player := range players {
		if player.ID != "" {
			count++
		}
	}
	return count
}

// GetLastPlayerInTimeline who is the last man
func GetLastPlayerInTimeline(players model.Players) (int, model.Player) {
	last := model.Player{}
	index := -1
	for i, player := range players {
		if player.StartLine > last.StartLine {
			last = player
			index = i
		}
	}
	return index, last
}

// IsPlayingAndNotFold if player is not fold and playing
func IsPlayingAndNotFold(player model.Player) bool {
	return player.Action.Name != constant.Fold && player.IsPlaying
}

// IsPlayingAndNotFoldAndNotAllIn if player is not fold and playing and not allin
func IsPlayingAndNotFoldAndNotAllIn(player model.Player) bool {
	return IsPlayingAndNotFold(player) && player.Action.Name != constant.AllIn
}

// IsPlayerBehindTheTimeline check if player behind the timeline
func IsPlayerBehindTheTimeline(player model.Player) bool {
	return time.Now().Unix() > player.DeadLine
}
