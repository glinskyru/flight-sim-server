package main

import (
	"log"
	"math"
	"time"
)

const (
	maxSpeed    = 10.0 // Units per second
	tolerance   = 0.1
	warning     = 1
	tempBan     = 3
	extendedBan = 5
)

func isValidMovement(oldPos, newPos Position, velocity Velocity) bool {
	distance := math.Sqrt(
		math.Pow(newPos.X-oldPos.X, 2) +
			math.Pow(newPos.Y-oldPos.Y, 2) +
			math.Pow(newPos.Z-oldPos.Z, 2),
	)
	expectedDistance := math.Sqrt(
		math.Pow(velocity.X*0.1, 2) +
			math.Pow(velocity.Y*0.1, 2) +
			math.Pow(velocity.Z*0.1, 2),
	)
	return distance <= expectedDistance+tolerance && distance <= maxSpeed*0.1
}

func handleViolation(player *Player) {
	player.Violations++
	switch {
	case player.Violations >= extendedBan:
		player.BanUntil = time.Now().Add(180 * 24 * time.Hour) // 6 months
		log.Printf("Player %s banned for 6 months", player.ID)
	case player.Violations >= tempBan:
		player.BanUntil = time.Now().Add(1 * time.Hour) // 1 hour
		log.Printf("Player %s temporarily banned for 1 hour", player.ID)
	case player.Violations >= warning:
		log.Printf("Warning issued to player %s for suspicious movement", player.ID)
	}
}
