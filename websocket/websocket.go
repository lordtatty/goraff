package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type WebSocketServer struct {
	Addr      string
	Upgrader  websocket.Upgrader
	clients   map[*websocket.Conn]bool
	broadcast chan string
}

func NewWebSocketServer(addr string) *WebSocketServer {
	return &WebSocketServer{
		Addr: addr,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan string),
	}
}

func (s *WebSocketServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	s.clients[ws] = true
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
	}
	delete(s.clients, ws)
}

func (server *WebSocketServer) handleMessages() {
	for {
		msg := <-server.broadcast
		for client := range server.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(server.clients, client)
			}
		}
	}
}

func (server *WebSocketServer) Serve() {
	http.HandleFunc("/ws", server.handleConnections)
	go server.handleMessages()
	log.Println("WebSocket server started on", server.Addr)
	log.Fatal(http.ListenAndServe(server.Addr, nil))
}

func (server *WebSocketServer) Send(msg string) {
	// msg to a json map
	server.broadcast <- msg
}

func (s *WebSocketServer) WaitForConnection() {
	// Create a channel to signal when a client connects
	connected := make(chan struct{})

	// Goroutine to wait for a client to connect
	go func() {
		for {
			// Check if there are any clients connected
			if len(s.clients) > 0 {
				// Signal that a client is connected
				close(connected)
				return
			}
			// Sleep for a short duration before checking again
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Wait until a client connects
	<-connected
}
