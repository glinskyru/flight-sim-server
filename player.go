package main

import (
	"time"

	"github.com/pion/webrtc/v3"
)

// Position represents a 3D coordinate
type Position struct {
	X, Y, Z float64
}

// Velocity represents movement speed
type Velocity struct {
	X, Y, Z float64
}

// Player represents a connected player
type Player struct {
	ID          string
	Position    Position
	Velocity    Velocity
	PeerConn    *webrtc.PeerConnection
	GameUpdates DataChannel
	UserInput   DataChannel
	Violations  int // For anti-cheat tracking
	BanUntil    time.Time
	DataChannel DataChannel // Now accepts any type implementing DataChannel
}
