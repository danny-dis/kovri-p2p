# Kovri P2P Mesh Network (Work in Progress)

This repository is being transformed into a decentralized, anonymous, and fast mesh networking system, inspired by GNUnet, Tor, and I2P principles.

## Current Status

We have set up the basic Go project structure and integrated `libp2p` for peer-to-peer communication. The application can now:

- Initialize a `libp2p` host.
- Utilize a Kademlia DHT for peer discovery.
- Connect to specified bootstrap peers for local network simulation.

## Getting Started (Local Development)

To run and test the mesh network locally, you will need to run multiple instances of the application. Each instance will act as a node in your mesh.

### Prerequisites

- Go (version 1.23.11 or higher)

### Running Multiple Nodes Locally

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/danny-dis/kovri-p2p.git
    cd kovri-p2p
    ```

2.  **Build the application:**

    ```bash
    /usr/local/go/bin/go build -o mesh-node .
    ```

3.  **Run the first node:**

    Open your first terminal and run:

    ```bash
    ./mesh-node
    ```

    This node will print its Peer ID and listening addresses. Copy one of the `/ip4/.../tcp/.../p2p/...` addresses. This will be your bootstrap peer address.

    Example output:
    ```
    Host created with ID: 12D3KooW... (your peer ID)
    Listening on addresses: [/ip4/127.0.0.1/tcp/4001/p2p/12D3KooW... /ip4/192.168.1.100/tcp/4001/p2p/12D3KooW...]
    Searching for peers...
    ```

4.  **Configure and run subsequent nodes:**

    Open a new terminal for each additional node you want to run. For each new node, create or modify its `config.yaml` to include the `bootstrap_peers` entry with the address from the first node.

    Example `config.yaml` for subsequent nodes:

    ```yaml
    exit_country: ""
    circuit_rotation_interval: "15m"
    bootstrap_peers:
      - "/ip4/127.0.0.1/tcp/4001/p2p/12D3KooW..." # Replace with the actual address from your first node
    ```

    Then, run the node:

    ```bash
    ./mesh-node
    ```

    You should see messages indicating that the nodes are discovering and connecting to each other.

## Next Steps

- Implement peer metadata advertisement (country code, bandwidth score, etc.).
- Develop circuit construction logic.
- Integrate TUN interface for traffic interception.