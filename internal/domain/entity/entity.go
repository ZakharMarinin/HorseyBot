package domain

import "time"

const (
	WaitingCommand    = "waiting_command"
	WaitingAction     = "waiting_action"
	WaitingImage      = "waiting_image"
	WaitingKeyword    = "waiting_keyword"
	WaitingUser       = "waiting_user"
	WaitingChance     = "waiting_chance"
	WaitingDelete     = "waiting_delete"
	WaitingFilter     = "waiting_for_filter"
	WaitingFilterData = "waiting_filter_data"
	WaitingChat       = "waiting_chat"
	WaitingOgo        = "waiting_ogo"
	WaitingOgoTimer   = "waiting_ogo_timer"
)

const (
	SendImage = iota + 1
	SendPing
	DeleteSub
	GetSubs
)

type User struct {
	UserID   int64
	UserName string
	ChatID   []int64
}

type Chat struct {
	ChatID   int64
	ChatName string
}

type TempUserState struct {
	UserID   int64  `json:"user_id"`
	ChatID   int64  `json:"chat_id"`
	ChatName string `json:"chat_name"`
	Filter   string `json:"filter"`
	Action   int    `json:"action"`
	Store    Store  `json:"store"`
	State    string `json:"state"`
}

type Subscription struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	ChatID   int64  `json:"chat_id"`
	ChatName string `json:"chat_name"`
	Feature  int    `json:"feature"`
	Store    Store  `json:"store"`
}

type Store struct {
	Threshold   int    `json:"threshold"`
	Chance      int    `json:"chance"`
	Image       string `json:"image"`
	ImageType   string `json:"image_type"`
	Keyword     string `json:"keyword"`
	TrackedUser string `json:"tracked_user"`
	StartTime   string `json:"start_time"`
	LastMessage string `json:"last_message"`
}

type OgoMeter struct {
	Count    int       `json:"count"`
	FirstOgo time.Time `json:"first_ogo"`
	LastOgo  time.Time `json:"last_ogo"`
	State    string    `json:"state"`
}
