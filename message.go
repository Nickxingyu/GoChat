package main

type Message struct {
	SenderID uint   `json:"sender_id"`
	Message  string `json:"message"`
}
