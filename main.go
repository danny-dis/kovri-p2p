package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mesh-network/pkg"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := pkg.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	h, dht, err := pkg.NewHost(ctx, cfg.BootstrapPeers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Host created with ID: %s\n", h.ID())
	fmt.Printf("Listening on addresses: %v\n", h.Addrs())

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