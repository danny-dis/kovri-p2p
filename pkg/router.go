package pkg

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// Circuit represents a path through the mesh network.
type Circuit struct {
	Hops []PeerInfo
	ExitNode PeerInfo
}

// CircuitManager manages the creation and rotation of circuits.
type CircuitManager struct {
	host host.Host
	peerInfoMgr *PeerInfoManager
	config *Config
	currentCircuit *Circuit
	circuitRotationTicker *time.Ticker
}

// NewCircuitManager creates a new CircuitManager.
func NewCircuitManager(h host.Host, pim *PeerInfoManager, cfg *Config) *CircuitManager {
	cm := &CircuitManager{
		host: h,
		peerInfoMgr: pim,
		config: cfg,
	}
	return cm
}

// StartCircuitRotation starts the periodic circuit rotation.
func (cm *CircuitManager) StartCircuitRotation(ctx context.Context) {
	interval := cm.config.CircuitRotationInterval
	if interval == 0 {
		interval = 15 * time.Minute // Default if not set in config
	}
	cm.circuitRotationTicker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-cm.circuitRotationTicker.C:
				log.Println("Rotating circuit...")
				cm.BuildNewCircuit(ctx)
			case <-ctx.Done():
				log.Println("Circuit rotation stopped.")
				cm.circuitRotationTicker.Stop()
				return
			}
		}
	}()
}

// BuildNewCircuit attempts to build a new circuit.
func (cm *CircuitManager) BuildNewCircuit(ctx context.Context) {
	peers := cm.peerInfoMgr.GetAllPeerInfo()
	if len(peers) < 3 {
		log.Println("Not enough peers to build a circuit (need at least 3).")
		return
	}

	// Filter for exit nodes if a country is specified
	var exitNodes []PeerInfo
	if cm.config.ExitCountry != "" {
		for _, p := range peers {
			if p.IsExitNode && p.CountryCode == cm.config.ExitCountry {
				exitNodes = append(exitNodes, p)
			}
		}
	} else {
		// If no exit country, consider all exit nodes (or all peers for now)
		for _, p := range peers {
			if p.IsExitNode {
				exitNodes = append(exitNodes, p)
			}
		}
		// For simplicity, if no explicit exit nodes, any peer can be an exit node
		if len(exitNodes) == 0 {
			for _, p := range peers {
				exitNodes = append(exitNodes, p)
			}
		}
	}

	if len(exitNodes) == 0 {
		log.Println("No suitable exit nodes found.")
		return
	}

	// Select exit node (for now, just pick one randomly)
	exitNode := exitNodes[rand.Intn(len(exitNodes))]

	// Select intermediate hops (excluding self and exit node)
	var intermediateHops []PeerInfo
	for _, p := range peers {
		if p.PeerID != cm.host.ID().String() && p.PeerID != exitNode.PeerID {
			intermediateHops = append(intermediateHops, p)
		}
	}

	if len(intermediateHops) < 2 {
		log.Println("Not enough intermediate hops to build a circuit (need at least 2).")
		return
	}

	// Randomly select 2 intermediate hops
	rand.Shuffle(len(intermediateHops), func(i, j int) {
		intermediateHops[i], intermediateHops[j] = intermediateHops[j], intermediateHops[i]
	})
	hops := intermediateHops[:2]

	// Construct the circuit
	newCircuit := &Circuit{
		Hops: hops,
		ExitNode: exitNode,
	}
	cm.currentCircuit = newCircuit
	log.Printf("New circuit built: Hops: %v, Exit: %s\n", newCircuit.Hops, newCircuit.ExitNode.PeerID)
}