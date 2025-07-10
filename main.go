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

	// Initialize GeoIP Service
	geoIPService, err := pkg.NewGeoIPService(cfg.GeoIPDatabasePath)
	if err != nil {
		log.Printf("Warning: Could not load GeoIP database. Country codes will not be available: %v", err)
		// Continue without GeoIP if it fails to load
	}
	defer func() {
		if geoIPService != nil {
			if err := geoIPService.Close(); err != nil {
				log.Printf("Error closing GeoIP database: %v", err)
			}
		}
	}()

	// Create local peer info (placeholders for now)
	localPeerInfo := pkg.PeerInfo{
		PeerID:        "", // Will be filled by NewHost
		IP:            "", // Will be filled by NewHost
		CountryCode:   "", // Will be filled by NewHost using GeoIP
		BandwidthScore: 100.0,               // Placeholder
		IsExitNode:    true,                  // Placeholder
	}

	h, dht, peerInfoMgr, err := pkg.NewHost(ctx, cfg.BootstrapPeers, &localPeerInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Populate localPeerInfo with actual host details and GeoIP
	// The PeerID and IP are already set inside NewHost
	if geoIPService != nil && localPeerInfo.IP != "" {
		// Extract just the IP address from the multiaddr string
		parsedAddr, err := pkg.ParseMultiaddrForIP(localPeerInfo.IP)
		if err == nil {
			country, err := geoIPService.GetCountryCode(parsedAddr)
			if err == nil {
				localPeerInfo.CountryCode = country
			} else {
				log.Printf("Warning: Could not get country code for %s: %v", parsedAddr, err)
			}
		} else {
			log.Printf("Warning: Could not parse IP from multiaddr %s: %v", localPeerInfo.IP, err)
		}
	}

	fmt.Printf("Host created with ID: %s\n", h.ID())
	fmt.Printf("Listening on addresses: %v\n", h.Addrs())

	// Initialize and start CircuitManager
	cm := pkg.NewCircuitManager(h, peerInfoMgr, cfg)
	cm.BuildNewCircuit(ctx) // Build initial circuit
	cm.StartCircuitRotation(ctx)

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
			// Also print current circuit info
			if cm.GetCurrentCircuit() != nil {
				fmt.Printf("Current Circuit: Hops: %v, Exit: %s\n", cm.GetCurrentCircuit().Hops, cm.GetCurrentCircuit().ExitNode.PeerID)
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
