package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type socketStruct struct {
	master *websocket.Conn
}

var SocketHandler = &socketStruct{}

func (apiconfig apiConfig) HandleWebSocketConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error while connecting to socket server", err)
		return
	}
	SocketHandler.master = conn

}
