package main

import (
	"encoding/json"
	"log"
	"time"
)

type DataChannel interface {
	Send(data []byte) error
	Close() error
	// Include any other methods your code calls on DataChannel
}

type Game struct {
	DataChannel DatraChannel
}

func updateGameState() {
	playersMutex.Lock()
	defer playersMutex.Unlock()
	for _, player := range players {
		if player.BanUntil.After(time.Now()) {
			continue
		}
		// Basic velocity-based movement
		newPos := Position{
			X: player.Position.X + player.Velocity.X*0.1, // 100ms tick
			Y: player.Position.Y + player.Velocity.Y*0.1,
			Z: player.Position.Z + player.Velocity.Z*0.1,
		}
		// Anti-cheat validation
		if isValidMovement(player.Position, newPos, player.Velocity) {
			player.Position = newPos
		} else {
			handleViolation(player)
		}
	}
}

func broadcastPositions() {
	playersMutex.Lock()
	defer playersMutex.Unlock()
	for _, player := range players {
		if player.BanUntil.After(time.Now()) {
			continue
		}
		data, err := json.Marshal(player.Position)
		if err != nil {
			log.Printf("Failed to serialize position for %s: %v", player.ID, err)
			continue
		}
		if err := player.GameUpdates.Send(data); err != nil {
			log.Printf("Failed to send position to %s: %v", player.ID, err)
		}
	}
}
