package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

const PeerInfoProtocol = "/mesh-network/peer-info/1.0.0"

// PeerInfo holds metadata about a peer.
type PeerInfo struct {
	PeerID        string         `json:"peer_id"` // Changed to string
	IP            string         `json:"ip"`
	CountryCode   string         `json:"country_code"`
	BandwidthScore float64      `json:"bandwidth_score"`
	IsExitNode    bool           `json:"is_exit_node"`
}

// PeerInfoManager manages discovered peer information.
type PeerInfoManager struct {
	mu    sync.RWMutex
	peers map[peer.ID]PeerInfo
}

// NewPeerInfoManager creates a new PeerInfoManager.
func NewPeerInfoManager() *PeerInfoManager {
	return &PeerInfoManager{
		peers: make(map[peer.ID]PeerInfo),
	}
}

// AddPeerInfo adds or updates peer information.
func (pm *PeerInfoManager) AddPeerInfo(info PeerInfo) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	// Convert string PeerID back to peer.ID for map key
	id, err := peer.Decode(info.PeerID)
	if err != nil {
		log.Printf("Error decoding peer ID %s: %v\n", info.PeerID, err)
		return
	}
	pm.peers[id] = info
}

// GetPeerInfo retrieves peer information by ID.
func (pm *PeerInfoManager) GetPeerInfo(id peer.ID) (PeerInfo, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	info, ok := pm.peers[id]
	return info, ok
}

// GetAllPeerInfo returns all stored peer information.
func (pm *PeerInfoManager) GetAllPeerInfo() []PeerInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	infos := make([]PeerInfo, 0, len(pm.peers))
	for _, info := range pm.peers {
		infos = append(infos, info)
	}
	return infos
}

// NewHost creates a new libp2p host and initializes a Kademlia DHT.
// It also sets up the PeerInfo exchange handler.
func NewHost(ctx context.Context, bootstrapPeers []string, localPeerInfo *PeerInfo) (host.Host, *dht.IpfsDHT, *PeerInfoManager, error) {
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.NATPortMap(), // Enable NAT port mapping
	)
	if err != nil {
		return nil, nil, nil, err
	}

	// Populate localPeerInfo with actual host details immediately
	localPeerInfo.PeerID = h.ID().String() // Convert to string
	// Use the first non-loopback address if available, otherwise fallback
	for _, addr := range h.Addrs() {
		if strings.Contains(addr.String(), "/ip4/") && !strings.Contains(addr.String(), "127.0.0.1") && !strings.Contains(addr.String(), "0.0.0.0") {
			localPeerInfo.IP = addr.String()
			break
		}
	}
	if localPeerInfo.IP == "" && len(h.Addrs()) > 0 {
		localPeerInfo.IP = h.Addrs()[0].String()
	}

	// Create a new Kademlia DHT
	dht, err := dht.New(ctx, h, dht.Mode(dht.ModeServer))
	if err != nil {
		return nil, nil, nil, err
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
			return nil, nil, nil, err
		}
	}

	// Setup PeerInfo exchange handler
	peerInfoMgr := NewPeerInfoManager()
	h.SetStreamHandler(PeerInfoProtocol, func(s network.Stream) {
		defer s.Close()
		remotePeer := s.Conn().RemotePeer()
		log.Printf("Handling PeerInfo stream from %s\n", remotePeer)

		// Send our PeerInfo
		encoder := json.NewEncoder(s)
		if err := encoder.Encode(localPeerInfo); err != nil {
			log.Printf("Error sending PeerInfo to %s: %v\n", remotePeer, err)
			return
		}

		// Receive remote PeerInfo
		var remotePeerInfo PeerInfo
		decoder := json.NewDecoder(s)
		if err := decoder.Decode(&remotePeerInfo); err != nil {
			log.Printf("Error receiving PeerInfo from %s: %v\n", remotePeer, err)
			return
		}

		peerInfoMgr.AddPeerInfo(remotePeerInfo)
		log.Printf("Received PeerInfo from %s: %+v\n", remotePeer, remotePeerInfo)
	})

	// Periodically discover and connect to new peers
	go discoverPeers(ctx, h, dht, peerInfoMgr, *localPeerInfo)

	return h, dht, peerInfoMgr, nil
}

// ParseMultiaddrForIP extracts the IP address from a multiaddr string.
func ParseMultiaddrForIP(multiaddrStr string) (string, error) {
	maddr, err := multiaddr.NewMultiaddr(multiaddrStr)
	if err != nil {
		return "", fmt.Errorf("invalid multiaddress: %w", err)
	}

	// Extract the IP component
	for _, p := range maddr.Protocols() {
		if p.Code == multiaddr.P_IP4 || p.Code == multiaddr.P_IP6 {
			return maddr.ValueForProtocol(p.Code)
		}
	}
	return "", fmt.Errorf("no IP address found in multiaddress: %s", multiaddrStr)
}

// discoverPeers continuously finds and connects to new peers.
func discoverPeers(ctx context.Context, h host.Host, dht *dht.IpfsDHT, peerInfoMgr *PeerInfoManager, localPeerInfo PeerInfo) {
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
				// If already connected, try to exchange PeerInfo if not already done
				if _, ok := peerInfoMgr.GetPeerInfo(p.ID); !ok {
					go exchangePeerInfo(ctx, h, p.ID, peerInfoMgr, localPeerInfo)
				}
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
					// Exchange PeerInfo immediately after connecting
					go exchangePeerInfo(ctx, h, p.ID, peerInfoMgr, localPeerInfo)
				}
			}(p)
		}
		time.Sleep(30 * time.Second) // Wait before searching again
	}
}

// exchangePeerInfo opens a new stream to a peer and exchanges PeerInfo.
func exchangePeerInfo(ctx context.Context, h host.Host, peerID peer.ID, peerInfoMgr *PeerInfoManager, localPeerInfo PeerInfo) {
	log.Printf("Exchanging PeerInfo with %s\n", peerID)
	stream, err := h.NewStream(ctx, peerID, PeerInfoProtocol)
	if err != nil {
		log.Printf("Error opening PeerInfo stream to %s: %v\n", peerID, err)
		return
	}
	defer stream.Close()

	// Send our PeerInfo
	encoder := json.NewEncoder(stream)
	if err := encoder.Encode(localPeerInfo); err != nil {
		log.Printf("Error sending PeerInfo to %s: %v\n", peerID, err)
		return
	}

	// Receive remote PeerInfo
	var remotePeerInfo PeerInfo
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&remotePeerInfo); err != nil {
		log.Printf("Error receiving PeerInfo from %s: %v\n", peerID, err)
		return
	}

	peerInfoMgr.AddPeerInfo(remotePeerInfo)
	log.Printf("Received PeerInfo from %s: %+v\n", peerID, remotePeerInfo)
}