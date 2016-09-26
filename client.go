package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
		Subprotocols:    []string{"webkitten-v1"},
	}
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (client *Client) readPump() {
	defer func() {
		client.hub.unregister <- client
		closeErr := client.conn.Close()
		if closeErr != nil {
			Error.Println(closeErr.Error())
		}
	}()
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, messageErr := client.conn.ReadMessage()
		if messageErr != nil {
			Error.Println(messageErr.Error())
			break
		}
		Info.Println(string(message))
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		closeErr := client.conn.Close()
		if closeErr != nil {
			Error.Println(closeErr.Error())
		}
	}()
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				writeErr := client.write(websocket.CloseMessage, []byte{})
				if writeErr != nil {
					Error.Println(writeErr.Error())
				}
				return
			}
			if writeErr := client.write(websocket.TextMessage, []byte(message)); writeErr != nil {
				return
			}
		case <-ticker.C:
			if writeErr := client.write(websocket.PingMessage, []byte{}); writeErr != nil {
				return
			}
		}
	}
}

func (client *Client) write(mType int, payload []byte) error {
	client.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return client.conn.WriteMessage(mType, payload)
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, connErr := upgrader.Upgrade(w, r, nil)
	if connErr != nil {
		Error.Println(connErr)
		return
	}
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, maxMessageSize),
	}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}
