package main

type Message struct {
	SenderID string `json:"sender_id"`
	Name     string `json:"name"`
	Message  string `json:"message"`
}
