package main

import (
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
	content    []byte

	channelMerge *ChannelMerge
}

type ChannelMerge struct {
	sync.Mutex
	inputs []chan string
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),

		channelMerge: &ChannelMerge{
			inputs: []chan string{},
		},
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
			client.send <- []byte(hub.content)
			break
		case client := <-hub.unregister:
			_, ok := hub.clients[client]
			if ok {
				delete(hub.clients, client)
				close(client.send)
			}
			break
		case message := <-hub.broadcast:
			hub.broadcastMessage(message)
		}
	}
}

func (hub *Hub) broadcastMessage(message string) {
	hub.content = append([]byte(message), hub.content...)[:maxMessageSize]
	for client := range hub.clients {
		select {
		case client.send <- []byte(message):
			break
		default:
			close(client.send)
			delete(hub.clients, client)
		}
	}
}

func (hub *Hub) addInput(ch chan string) {
	hub.channelMerge.Lock()
	defer hub.channelMerge.Unlock()
	hub.channelMerge.inputs = append(hub.channelMerge.inputs, ch)
}

func (hub *Hub) merge() {
	hub.channelMerge.Lock()
	defer hub.channelMerge.Unlock()
	for _, ch := range hub.channelMerge.inputs {
		go func(c chan string) {
			for {
				message, ok := <-c
				if !ok {
					break
				}
				hub.broadcast <- message
			}
			return
		}(ch)
	}
}
