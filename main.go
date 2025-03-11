package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	players      = make(map[string]*Player)
	playersMutex sync.Mutex
)

func main() {
	// Start signaling server
	http.HandleFunc("/offer", handleOffer)
	http.HandleFunc("/disconnect", handleDisconnect)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Start game loop
	go gameLoop()

	// Keep server running
	select {}
}

// Game loop to update and broadcast positions every 100ms
func gameLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		updateGameState()
		broadcastPositions()
	}
}
