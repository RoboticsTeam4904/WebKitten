package main

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
	content    []byte
}

var ()

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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
