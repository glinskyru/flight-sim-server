package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/pion/webrtc/v3"
)

var playerIDCounter uint64

func generateUniqueID() string {
	id := atomic.AddUint64(&playerIDCounter, 1)
	return "player_" + string(id)
}

func handleOffer(w http.ResponseWriter, r *http.Request) {
	var offer webrtc.SessionDescription
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		http.Error(w, "Invalid offer", http.StatusBadRequest)
		return
	}

	// Create peer connection
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}
	peerConn, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Printf("Failed to create peer connection: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Set remote description
	if err := peerConn.SetRemoteDescription(offer); err != nil {
		log.Printf("Failed to set remote description: %v", err)
		peerConn.Close()
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Create data channels
	gameUpdates, err := peerConn.CreateDataChannel("game_updates", nil)
	if err != nil {
		log.Printf("Failed to create game_updates channel: %v", err)
		peerConn.Close()
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	userInput, err := peerConn.CreateDataChannel("user_input", nil)
	if err != nil {
		log.Printf("Failed to create user_input channel: %v", err)
		peerConn.Close()
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Create player
	playerID := generateUniqueID()
	player := &Player{
		ID:          playerID,
		Position:    Position{0, 0, 0},
		Velocity:    Velocity{0, 0, 0},
		PeerConn:    peerConn,
		GameUpdates: gameUpdates,
		UserInput:   userInput,
	}

	// Handle user input
	userInput.OnMessage(func(msg webrtc.DataChannelMessage) {
		velocity := parseInput(msg.Data)
		playersMutex.Lock()
		if p, ok := players[playerID]; ok && p.BanUntil.IsZero() {
			p.Velocity = velocity
		}
		playersMutex.Unlock()
	})

	// Create answer
	answer, err := peerConn.CreateAnswer(nil)
	if err != nil {
		log.Printf("Failed to create answer: %v", err)
		peerConn.Close()
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if err := peerConn.SetLocalDescription(answer); err != nil {
		log.Printf("Failed to set local description: %v", err)
		peerConn.Close()
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Store player
	playersMutex.Lock()
	players[playerID] = player
	playersMutex.Unlock()

	// Send answer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}

func handleDisconnect(w http.ResponseWriter, r *http.Request) {
	var data struct {
		PlayerID string `json:"player_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	playersMutex.Lock()
	if player, ok := players[data.PlayerID]; ok {
		player.PeerConn.Close()
		delete(players, data.PlayerID)
		log.Printf("Player %s disconnected", data.PlayerID)
	}
	playersMutex.Unlock()
	w.WriteHeader(http.StatusOK)
}

// Mock input parsing (in production, decode actual client input)
func parseInput(data []byte) Velocity {
	return Velocity{X: 1.0, Y: 0, Z: 0} // Simplified for testing
}
