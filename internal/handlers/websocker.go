package handlers

import "net/http"

type WSHandler struct {
}

func NewWsHandler() *WSHandler {
	return &WSHandler{}
}

func (ws WSHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}
