package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	serverURL := "ws://localhost:8080/ws"
	
	fmt.Println("Connecting to Ultimate Tic-Tac-Toe server...")
	
	// Connect to WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()
	
	fmt.Println("Connected to server!")
	
	// Handle interrupt signal to gracefully close connection
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	
	// Channel to receive messages from server
	done := make(chan struct{})
	
	// Goroutine to read messages from server
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			fmt.Printf("\nServer: %s\n", string(message))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your move (e.g., A2, D5): ")

	// Main loop for sending moves
	go func() {
		for scanner.Scan() {
			move := strings.TrimSpace(scanner.Text())

			if move == "quit" || move == "exit" {
				fmt.Println("Disconnecting...")
				close(done)
				return
			}

			if move != "" {
				err := conn.WriteMessage(websocket.TextMessage, []byte(move))
				if err != nil {
					log.Println("Error sending move:", err)
					return
				}
			}

			fmt.Print("Enter your move (e.g., A2, D5): ")
		}
	}()
	
	// Wait for interrupt signal or done channel
	select {
	case <-done:
		return
	case <-interrupt:
		fmt.Println("\nReceived interrupt signal. Closing connection...")
		
		// Send close message to server
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("Error sending close message:", err)
			return
		}
		
		return
	}
}
