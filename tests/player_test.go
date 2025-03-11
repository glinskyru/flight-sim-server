package flight_sim_server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pion/webrtc/v3"
)

func TestHandleOffer(t *testing.T) {
	// Mock offer
	offer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: "mock_sdp"}
	offerData, _ := json.Marshal(offer)
	req := httptest.NewRequest("POST", "/offer", bytes.NewReader(offerData))
	w := httptest.NewRecorder()

	// Clear players
	playersMutex.Lock()
	players = make(map[string]*Player)
	playersMutex.Unlock()

	handleOffer(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	playersMutex.Lock()
	defer playersMutex.Unlock()
	if len(players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(players))
	}
	for _, p := range players {
		if p.GameUpdates == nil || p.UserInput == nil {
			t.Error("Data channels not initialized")
		}
	}
}

func TestHandleDisconnect(t *testing.T) {
	playersMutex.Lock()
	players["test_player"] = &Player{ID: "test_player", PeerConn: &webrtc.PeerConnection{}}
	playersMutex.Unlock()

	data, _ := json.Marshal(struct{ PlayerID string }{PlayerID: "test_player"})
	req := httptest.NewRequest("POST", "/disconnect", bytes.NewReader(data))
	w := httptest.NewRecorder()

	handleDisconnect(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	playersMutex.Lock()
	defer playersMutex.Unlock()
	if len(players) != 0 {
		t.Errorf("Expected 0 players, got %d", len(players))
	}
}
