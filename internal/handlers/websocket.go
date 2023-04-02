package handlers

import (
	"aphrodite/internal/config"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	Upgrader websocket.Upgrader
	Dio         string
	Yev         string
}

var (
	Connections map[string]*websocket.Conn
	Mutex       sync.Mutex

)

func NewWsHandler(configuration config.WebSocketConfig) *WSHandler {
	Connections = map[string]*websocket.Conn{}
	return &WSHandler{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  configuration.ReadBufferSize,
			WriteBufferSize: configuration.WriteBufferSize,
		},
		Dio: configuration.Dio,
		Yev: configuration.Yev,
	}
}

func (ws WSHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, client := range Connections {
			client.WriteMessage(websocket.TextMessage, []byte("ping"))
		}
		w.Write([]byte("pong"))
		log.Println("Pong")
	}
}

func (ws WSHandler) HandleConnections() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := ws.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		if Connections[mux.Vars(r)["id"]] != nil {
			Connections[mux.Vars(r)["id"]].Close()
		}

		Mutex.Lock()
		Connections[mux.Vars(r)["id"]] = conn
		log.Printf("New client connected. Total clients: %d", len(Connections))
		Mutex.Unlock()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				Mutex.Lock()
				delete(Connections, mux.Vars(r)["id"])
				log.Printf("Client disconnected. Total clients: %d", len(Connections))
				Mutex.Unlock()

				log.Printf("Client disconnected: %v", err)
				break
			}
		}
	}
}

func (ws WSHandler) GetConnections() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Total clients: %d", len(Connections))
		Mutex.Lock()
		w.Write([]byte("Number of connections: " + fmt.Sprint(len(Connections))))
		Mutex.Unlock()
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

func (ws WSHandler) DIOEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if Connections[ws.Yev] == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		Connections[ws.Yev].WriteMessage(websocket.TextMessage, []byte("ping"))

		w.WriteHeader(http.StatusOK)
	}
}

func (ws WSHandler) YEVEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if Connections[ws.Dio] == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		Connections[ws.Dio].WriteMessage(websocket.TextMessage, []byte("ping"))

		w.WriteHeader(http.StatusOK)
	}
}
