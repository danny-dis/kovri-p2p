package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	
	"syscall"
	"time"

	"mesh-network/pkg"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := pkg.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create local peer info (placeholders for now)
	localPeerInfo := pkg.PeerInfo{
		PeerID:        "", // Will be filled by NewHost
		IP:            "", // Will be filled by NewHost
		CountryCode:   "US",                  // Placeholder, will be dynamic later
		BandwidthScore: 100.0,               // Placeholder
		IsExitNode:    true,                  // Placeholder
	}

	h, dht, peerInfoMgr, err := pkg.NewHost(ctx, cfg.BootstrapPeers, &localPeerInfo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Host created with ID: %s\n", h.ID())
	fmt.Printf("Listening on addresses: %v\n", h.Addrs())

	// Periodically print discovered peer info for debugging
	go func() {
		for {
			fmt.Println("\n--- Discovered Peers ---")
			peers := peerInfoMgr.GetAllPeerInfo()
			if len(peers) == 0 {
				fmt.Println("No peers discovered yet.")
			}
			for _, p := range peers {
				fmt.Printf("ID: %s, IP: %s, Country: %s, Bandwidth: %.2f, Exit: %t\n",
					p.PeerID, p.IP, p.CountryCode, p.BandwidthScore, p.IsExitNode)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	// Wait for a termination signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Println("Shutting down host...")
	if err := h.Close(); err != nil {
		log.Println("Error closing host:", err)
	}
	if dht != nil {
		if err := dht.Close(); err != nil {
			log.Println("Error closing DHT:", err)
		}
	}
	fmt.Println("Host shut down.")
}