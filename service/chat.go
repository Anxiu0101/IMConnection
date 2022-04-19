package service

import "github.com/gorilla/websocket"

type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}
