package model

// Cards array of card
type Cards []int

// Deck of cards
type Deck struct {
	Cards Cards `json:"cards"`
}

// Suit is for human readable
var Suit = [4]string{"clubs", "diamonds", "hearts", "spades"}

// Number is for human readable
var Number = [13]string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
