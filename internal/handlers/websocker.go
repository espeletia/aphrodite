package handlers

import (
	"aphrodite/internal/config"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSHandler struct {
	Upgrader websocket.Upgrader
}

func NewWsHandler(configuration config.WebSocketConfig) *WSHandler {
	return &WSHandler{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  configuration.ReadBufferSize,
			WriteBufferSize: configuration.WriteBufferSize,
		},
	}
}

func (ws WSHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

func (ws WSHandler) Echo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := ws.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			log.Println("Message: ", string(message))
			err = conn.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	}
}