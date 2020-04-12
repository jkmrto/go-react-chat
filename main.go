package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func home(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./assets/home.html")
}

func main() {
	connections := &[]*websocket.Conn{}
	http.HandleFunc("/ws", echoHandler(connections))
	http.HandleFunc("/home", home)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("%+v", err)
	}
}

func echoHandler(connections *[]*websocket.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("%+v", err)
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
