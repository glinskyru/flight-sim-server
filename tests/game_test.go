package flight_sim_server

import (
	"encoding/json"
	"testing"
	"time"
)

type mockDataChannel struct {
	sentData [][]byte
}

func (m *mockDataChannel) Send(data []byte) error {
	m.sentData = append(m.sentData, data)
	return nil
}

func TestBroadcastPositions(t *testing.T) {
	playersMutex.Lock()
	players = make(map[string]*Player)
	mockChannel := &mockDataChannel{}
	players["p1"] = &Player{
		ID:          "p1",
		Position:    Position{1, 2, 3},
		GameUpdates: mockChannel,
	}
	playersMutex.Unlock()

	broadcastPositions()

	if len(mockChannel.sentData) != 1 {
		t.Errorf("Expected 1 broadcast, got %d", len(mockChannel.sentData))
	}
	var pos Position
	json.Unmarshal(mockChannel.sentData[0], &pos)
	if pos != (Position{1, 2, 3}) {
		t.Errorf("Expected position {1,2,3}, got %+v", pos)
	}
}

func TestGameLoopTiming(t *testing.T) {
	playersMutex.Lock()
	players = make(map[string]*Player)
	mockChannel := &mockDataChannel{}
	players["p1"] = &Player{
		ID:          "p1",
		Position:    Position{0, 0, 0},
		Velocity:    Velocity{1, 0, 0},
		GameUpdates: mockChannel,
	}
	playersMutex.Unlock()

	start := time.Now()
	go gameLoop()
	time.Sleep(250 * time.Millisecond) // Allow 2+ ticks

	playersMutex.Lock()
	p := players["p1"]
	playersMutex.Unlock()

	elapsed := time.Since(start).Seconds()
	expectedTicks := int(elapsed / 0.1)
	expectedX := float64(expectedTicks) * 0.1 // Velocity * tick duration
	if p.Position.X < expectedX-0.1 || p.Position.X > expectedX+0.1 {
		t.Errorf("Expected X around %.1f, got %.1f", expectedX, p.Position.X)
	}
}
