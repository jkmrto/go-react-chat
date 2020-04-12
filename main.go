package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echoHandler(connections *[]*websocket.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		(*connections) = append((*connections), c)

		defer c.Close()
		for {
			// Read simultaneously
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}

			// Send message to all sockets
			for _, c := range *connections {
				if c.WriteMessage(mt, message); err != nil {
					break
				}
			}
		}
	}
}
