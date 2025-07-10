package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

// NewHost creates a new libp2p host and initializes a Kademlia DHT.
func NewHost(ctx context.Context, bootstrapPeers []string) (host.Host, *dht.IpfsDHT, error) {
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.NATPortMap(), // Enable NAT port mapping
	)
	if err != nil {
		return nil, nil, err
	}

	// Create a new Kademlia DHT
	dht, err := dht.New(ctx, h, dht.Mode(dht.ModeServer))
	if err != nil {
		return nil, nil, err
	}

	// Bootstrap the DHT.
	if len(bootstrapPeers) > 0 {
		fmt.Println("Connecting to bootstrap peers...")
		var pis []peer.AddrInfo
		for _, addr := range bootstrapPeers {
			maddr, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				fmt.Printf("Invalid bootstrap peer address %s: %v\n", addr, err)
				continue
			}
			pi, err := peer.AddrInfoFromP2pAddr(maddr)
			if err != nil {
				fmt.Printf("Invalid bootstrap peer AddrInfo %s: %v\n", addr, err)
				continue
			}
			pis = append(pis, *pi)
		}
		if err := h.Connect(ctx, pis[0]); err != nil {
			fmt.Printf("Failed to connect to bootstrap peer: %v\n", err)
		} else {
			fmt.Println("Connected to bootstrap peer.")
		}
	} else {
		// Default to public IPFS bootstrap peers if none are specified
		if err = dht.Bootstrap(ctx); err != nil {
			return nil, nil, err
		}
	}

	// Periodically discover and connect to new peers
	go discoverPeers(ctx, h, dht)

	return h, dht, nil
}

// discoverPeers continuously finds and connects to new peers.
func discoverPeers(ctx context.Context, h host.Host, dht *dht.IpfsDHT) {
	routingDiscovery := routing.NewRoutingDiscovery(dht)
	// Advertise our presence to the DHT
	routingDiscovery.Advertise(ctx, "mesh-network")

	for {
		fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, "mesh-network")
		if err != nil {
			fmt.Printf("Error finding peers: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		for p := range peerChan {
			if p.ID == h.ID() || len(p.Addrs) == 0 {
				continue
			}
			fmt.Printf("Found peer: %s with addresses: %v\n", p.ID, p.Addrs)
			// Check if already connected using the ConnManager
			if h.Network().Connectedness(p.ID) == network.Connected {
				fmt.Printf("Already connected to peer %s\n", p.ID)
				continue
			}

			// Connect to the peer
			go func(p peer.AddrInfo) {
				fmt.Printf("Connecting to peer %s...\n", p.ID)
				ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
				defer cancel()
				if err := h.Connect(ctx, p); err != nil {
					fmt.Printf("Failed to connect to %s: %v\n", p.ID, err)
				} else {
					fmt.Printf("Connected to peer %s\n", p.ID)
				}
			}(p)
		}
		time.Sleep(30 * time.Second) // Wait before searching again
	}
}
