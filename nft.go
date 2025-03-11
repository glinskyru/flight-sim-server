package main

import (
	"errors"
	"log"
	"time"
)

// Mock blockchain client (replace with real Bitcoin Ordinals API in production)
type BlockchainClient struct{}

func (c *BlockchainClient) TransferNFT(nftID, fromPlayerID, toPlayerID string) error {
	// Simulate blockchain call
	time.Sleep(100 * time.Millisecond) // Mock latency
	return nil
}

var (
	nftOwnership = make(map[string]string) // NFT ID to Player ID
	bcClient     = &BlockchainClient{}
)

func transferNFT(nftID, fromPlayerID, toPlayerID, jwtToken string) error {
	// Validate JWT for purchase (mocked for simplicity)
	if jwtToken == "" {
		return errors.New("authentication required for NFT transfer")
	}
	// Check ownership
	playersMutex.Lock()
	defer playersMutex.Unlock()
	if owner, ok := nftOwnership[nftID]; !ok || owner != fromPlayerID {
		return errors.New("invalid ownership")
	}
	// Perform blockchain transfer
	if err := bcClient.TransferNFT(nftID, fromPlayerID, toPlayerID); err != nil {
		log.Printf("Blockchain transfer failed: %v", err)
		return err
	}
	// Update off-chain tracking
	nftOwnership[nftID] = toPlayerID
	log.Printf("NFT %s transferred from %s to %s", nftID, fromPlayerID, toPlayerID)
	return nil
}
