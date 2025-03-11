package flight_sim_server

import (
	"testing"
)

func TestTransferNFT(t *testing.T) {
	playersMutex.Lock()
	players = make(map[string]*Player)
	players["p1"] = &Player{ID: "p1"}
	players["p2"] = &Player{ID: "p2"}
	nftOwnership = make(map[string]string)
	nftOwnership["nft1"] = "p1"
	playersMutex.Unlock()

	// Success case
	err := transferNFT("nft1", "p1", "p2", "mock_jwt")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if nftOwnership["nft1"] != "p2" {
		t.Errorf("Expected owner p2, got %s", nftOwnership["nft1"])
	}

	// No auth
	err = transferNFT("nft1", "p2", "p1", "")
	if err == nil || err.Error() != "authentication required for NFT transfer" {
		t.Errorf("Expected auth error, got %v", err)
	}

	// Invalid ownership
	err = transferNFT("nft1", "p1", "p2", "mock_jwt")
	if err == nil || err.Error() != "invalid ownership" {
		t.Errorf("Expected ownership error, got %v", err)
	}
}
