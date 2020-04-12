package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestEchoHandler(t *testing.T) {
	t.Run("Given a connected client, when a message is sent,  the client received the message", func(t *testing.T) {
		connections := &[]*websocket.Conn{}
		s := httptest.NewServer(echoHandler(connections))
		defer s.Close()

		// Convert http://127.0.0.1 to ws://127.0.0.
		url := "ws" + strings.TrimPrefix(s.URL, "http")
		message := "hello"

		ws1 := openWsConnection(t, url)
		defer ws1.Close()

		writeMessage(t, ws1, message)
		require.Equal(t, message, readMessage(t, ws1))
	})

	t.Run("Given two clients connected, when a message is sent, it is read by all the open connections", func(t *testing.T) {
		connections := &[]*websocket.Conn{}
		s := httptest.NewServer(echoHandler(connections))
		defer s.Close()

		// Convert http://127.0.0.1 to ws://127.0.0.
		url := "ws" + strings.TrimPrefix(s.URL, "http")
		message := "hello"

		ws1 := openWsConnection(t, url)
		defer ws1.Close()

		ws2 := openWsConnection(t, url)
		defer ws2.Close()

		writeMessage(t, ws1, message)

		require.Equal(t, message, readMessage(t, ws1))
		require.Equal(t, message, readMessage(t, ws2))
	})
}

func openWsConnection(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	return ws
}

func readMessage(t *testing.T, ws *websocket.Conn) string {
	_, p, err := ws.ReadMessage()
	require.NoError(t, err)
	return string(p)
}

func writeMessage(t *testing.T, ws *websocket.Conn, message string) {
	err := ws.WriteMessage(websocket.TextMessage, []byte(message))
	require.NoError(t, err)
}
