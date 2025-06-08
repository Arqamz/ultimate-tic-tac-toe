package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

type GameServer struct {
	clients map[*websocket.Conn]string // map client connection to player name
}

func NewGameServer() *GameServer {
	return &GameServer{
		clients: make(map[*websocket.Conn]string),
	}
}

func (gs *GameServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Assign player (X or O) based on number of connected clients
	var player string
	if len(gs.clients) == 0 {
		player = "X"
	} else {
		player = "O"
	}
	
	gs.clients[conn] = player
	
	log.Printf("Player %s connected from %s", player, r.RemoteAddr)
	
	// Welcome message with player assignment
	welcomeMsg := fmt.Sprintf("Welcome! You are player %s", player)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg)); err != nil {
		log.Printf("Error sending welcome message: %v", err)
		return
	}
	
	// Send the basic board
	board := gs.generateBasicBoard()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(board)); err != nil {
		log.Printf("Error sending board: %v", err)
		return
	}

	// Listen for messages from client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		
		if messageType == websocket.TextMessage {
			move := string(message)
			response := fmt.Sprintf("Move %s received from %s", move, player)
			log.Println(response)
		}
	}
	
	// Clean up when client disconnects
	delete(gs.clients, conn)
	log.Printf("Player %s disconnected", player)
}

func (gs *GameServer) generateBasicBoard() string {
	board := `
Ultimate Tic-Tac-Toe Board:

   | A |         | B |         | C | 
-----------   -----------   -----------
 1 | 2 | 3     1 | 2 | 3     1 | 2 | 3
-----------   -----------   -----------
 4 | 5 | 6     4 | 5 | 6     4 | 5 | 6
-----------   -----------   -----------
 7 | 8 | 9     7 | 8 | 9     7 | 8 | 9

-----------   -----------   -----------

   | D |         | E |         | F | 
-----------   -----------   -----------
 1 | 2 | 3     1 | 2 | 3     1 | 2 | 3
-----------   -----------   -----------
 4 | 5 | 6     4 | 5 | 6     4 | 5 | 6
-----------   -----------   -----------
 7 | 8 | 9     7 | 8 | 9     7 | 8 | 9

-----------   -----------   -----------

   | G |         | H |         | I | 
-----------   -----------   -----------
 1 | 2 | 3     1 | 2 | 3     1 | 2 | 3
-----------   -----------   -----------
 4 | 5 | 6     4 | 5 | 6     4 | 5 | 6
-----------   -----------   -----------
 7 | 8 | 9     7 | 8 | 9     7 | 8 | 9

Send your move in format like: A2, D5, G9, etc.
`
	return board
}

func main() {
	gameServer := NewGameServer()
	
	http.HandleFunc("/ws", gameServer.handleWebSocket)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ultimate Tic-Tac-Toe WebSocket Server is running!\nConnect to ws://localhost:8080/ws")
	})
	
	log.Println("Server starting on :8080")
	log.Println("WebSocket endpoint: ws://localhost:8080/ws")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
