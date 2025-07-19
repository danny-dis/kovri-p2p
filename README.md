# Kovri P2P Mesh Network (Work in Progress)

This repository is being transformed into a decentralized, anonymous, and fast mesh networking system, inspired by GNUnet, Tor, and I2P principles. The goal is to create a robust and resilient network that prioritizes user privacy and data security.

## Table of Contents

- [Current Status](#current-status)
- [Project Goals](#project-goals)
- [Architecture Overview](#architecture-overview)
- [Getting Started (Local Development)](#getting-started-local-development)
  - [Prerequisites](#prerequisites)
  - [Project Structure](#project-structure)
  - [Configuration](#configuration)
  - [Running Multiple Nodes Locally](#running-multiple-nodes-locally)
- [Next Steps & Future Development](#next-steps--future-development)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Current Status

The project is currently in its foundational phase, focusing on establishing a stable and extensible peer-to-peer communication layer. We have successfully:

-   **Initialized a `libp2p` host:** The application can create and manage a `libp2p` host, enabling it to participate in the peer-to-peer network. This includes handling peer IDs, cryptographic key pairs, and network addresses.
-   **Utilized a Kademlia DHT for peer discovery:** A Kademlia Distributed Hash Table (DHT) is integrated to facilitate efficient peer discovery and routing within the network. Nodes can find each other and exchange information about available services.
-   **Connected to specified bootstrap peers for local network simulation:** For development and testing, nodes can connect to a predefined list of bootstrap peers. This allows for easy setup of a local mesh network without relying on external discovery mechanisms.

## Project Goals

The long-term vision for Kovri is to provide a comprehensive and secure mesh networking solution with the following key objectives:

-   **Decentralization:** Eliminate single points of failure and central authorities to enhance network resilience and censorship resistance.
-   **Anonymity:** Implement advanced routing and cryptographic techniques to protect user identities and communication patterns.
-   **Speed and Efficiency:** Optimize network performance for various applications, including file sharing, messaging, and general internet browsing.
-   **Extensibility:** Design a modular architecture that allows for easy integration of new features, protocols, and applications.
-   **Usability:** Strive for a user-friendly experience, making anonymous communication accessible to a broader audience.

## Architecture Overview

The Kovri mesh network is being built with a layered and modular architecture to ensure flexibility, scalability, and maintainability.

-   **Core Network Layer (Go/libp2p):** This layer handles fundamental peer-to-peer communication, including peer discovery, connection management, and data transfer. It leverages `libp2p` for its robust and flexible networking capabilities.
    -   **Peer Identity:** Each node has a unique Peer ID derived from its cryptographic public key.
    -   **Kademlia DHT:** Used for efficient peer discovery, content routing, and storing/retrieving network information.
    -   **Transport Protocols:** Support for various transport protocols (e.g., TCP, UDP) to ensure connectivity across different network environments.
-   **Routing and Anonymity Layer (C++/Riffle, I2P):** This layer is responsible for anonymous routing of traffic through the mesh network. It will incorporate principles from I2P's garlic routing and potentially Riffle's verifiable shuffle for enhanced anonymity and efficiency.
    -   **Garlic Routing (I2P):** Encapsulates messages in multiple layers of encryption, routed through a series of volunteer nodes (tunnels) to obscure the origin and destination.
    -   **Verifiable Shuffle (Riffle - Planned):** Exploration and potential integration of Riffle's verifiable shuffle mechanism to provide strong anonymity guarantees with improved efficiency for certain communication patterns.
-   **Application Layer:** This layer will support various decentralized applications built on top of the anonymous routing infrastructure. Examples include:
    -   **Anonymous File Sharing:** Secure and private file exchange among network participants.
    -   **Secure Messaging:** Encrypted and anonymous communication channels.
    -   **Decentralized Services:** Hosting and accessing services within the mesh network without revealing server identities or locations.

## Getting Started (Local Development)

To run and test the mesh network locally, you will need to run multiple instances of the application. Each instance will act as a node in your mesh, allowing you to simulate network behavior and test communication.

### Prerequisites

-   **Go:** Version 1.23.11 or higher. Ensure Go is correctly installed and configured in your system's PATH. You can download it from [golang.org](https://golang.org/dl/).
-   **Git:** For cloning the repository.
-   **C++ Compiler:** A C++17 compatible compiler (e.g., GCC, Clang) is required for building the C++ components.
-   **CMake:** Version 3.5 or higher for building the C++ components.
-   **Boost Libraries:** Required for C++ components. Ensure you have `Boost.Asio`, `Boost.System`, `Boost.Log`, and `Boost.Program_options` installed.
-   **Crypto++ Library:** Required for cryptographic operations in C++ components.

### Project Structure

The repository is organized as follows:

```
kovri-p2p/
├───.git/                   # Git repository metadata
├───.github/                # GitHub Actions workflows and issue templates
├───build/                  # Build artifacts (generated by CMake/Make)
├───cmake/                  # CMake modules and scripts for C++ components
├───contrib/                # Third-party contributions, utilities, and packaging scripts
│   ├───pgp/                # PGP keys for developers
│   ├───pkg/                # Packaging related files (Docker, Snap)
│   ├───python/             # Python utilities and bindings
│   └───testnet/            # Scripts and configurations for testnet deployment
├───deps/                   # External dependencies (e.g., cpp-netlib, cryptopp, miniupnp, webrtc)
├───docs/                   # Project documentation
├───pkg/                    # Go packages for core functionality (config, geoip, peer, router, tun)
├───src/                    # C++ source code
│   ├───app/                # Application-level logic
│   ├───client/             # Client-side components
│   ├───core/               # Core C++ components
│   │   ├───crypto/         # Cryptographic primitives and algorithms
│   │   ├───riffle/         # Riffle protocol implementation (new)
│   │   ├───router/         # Routing logic, identity management, and network database
│   │   └───util/           # General utilities
│   └───util/               # General utilities
├───tests/                  # Unit and integration tests
│   ├───fuzz_tests/         # Fuzzing tests
│   └───unit_tests/         # Unit tests
├───CMakeLists.txt          # Main CMake build script for C++ components
├───config.yaml             # Default configuration file
├───go.mod                  # Go module definition
├───go.sum                  # Go module checksums
├───LICENSE.md              # Project license information
├───main.go                 # Main Go application entry point
├───Makefile                # Makefile for build automation
├───mesh-node               # Compiled Go binary (after build)
└───README.md               # This README file
```

### Configuration

The `config.yaml` file is used to configure various aspects of the node's behavior. Key parameters include:

-   `exit_country`: (String) Specifies the country to exit the network from (e.g., "US", "DE"). Leave empty for no specific exit country.
-   `circuit_rotation_interval`: (Duration) How often the node should rotate its anonymous circuits (e.g., "15m" for 15 minutes, "1h" for 1 hour).
-   `bootstrap_peers`: (List of Strings) A list of multiaddresses for bootstrap peers. These are known, stable nodes that help new nodes join the network. In local development, this will be the address of your first running node.

### Running Multiple Nodes Locally

To simulate a mesh network, you will run multiple instances of the `mesh-node` binary, each configured to connect to at least one other node.

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/danny-dis/kovri-p2p.git
    cd kovri-p2p
    ```

2.  **Build the application:**

    This command compiles the Go and C++ components and creates the `mesh-node` executable.

    ```bash
    /usr/local/go/bin/go build -o mesh-node .
    ```

3.  **Run the first node (Bootstrap Node):**

    Open your first terminal and run the `mesh-node` without any specific bootstrap peers. This node will act as the initial entry point for other nodes.

    ```bash
    ./mesh-node
    ```

    This node will print its Peer ID and listening addresses. **Copy one of the `/ip4/.../tcp/.../p2p/...` addresses.** This will be your bootstrap peer address for subsequent nodes.

    Example output from the first node:

    ```
    Host created with ID: 12D3KooW... (your peer ID)
    Listening on addresses: [/ip4/127.0.0.1/tcp/4001/p2p/12D3KooW... /ip4/192.168.1.100/tcp/4001/p2p/12D3KooW...]
    Searching for peers...
    ```

4.  **Configure and run subsequent nodes:**

    Open a new terminal for each additional node you want to run. For each new node, you need to create or modify its `config.yaml` to include the `bootstrap_peers` entry with the address obtained from the first node.

    Example `config.yaml` for subsequent nodes (create this file in the `kovri-p2p` directory for each new node, or modify the existing one):

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

    You should see messages indicating that the nodes are discovering and connecting to each other. Observe the logs for successful peer connections and DHT operations.

## Next Steps & Future Development

The development of Kovri is ongoing, with several key areas identified for future work:

-   **Peer Metadata Advertisement:** Implement mechanisms for peers to advertise additional metadata, such as country code, bandwidth capabilities, and reliability scores. This information can be used by the routing layer to build more efficient and robust anonymous paths.
-   **Circuit Construction Logic:** Develop the core logic for building and managing anonymous circuits (tunnels) through the network. This will involve selecting appropriate nodes, negotiating encryption keys, and establishing multi-hop paths.
-   **Integrate TUN Interface for Traffic Interception:** Implement a TUN (Tunnel) interface to intercept and redirect all network traffic through the Kovri mesh. This will enable transparent anonymity for any application running on the system.
-   **Riffle Protocol Integration (Advanced):** Further research and integrate Riffle's verifiable shuffle protocol into the C++ core. This could offer alternative anonymity properties and performance characteristics compared to traditional onion routing.
-   **Application Development:** Build and integrate various decentralized applications (e.g., secure messaging, file sharing) that leverage the underlying anonymous network.
-   **Security Audits and Hardening:** Conduct thorough security audits and implement hardening measures to protect against various attacks, including traffic analysis, deanonymization, and denial-of-service.
-   **Cross-Platform Compatibility:** Ensure the project compiles and runs seamlessly on various operating systems (Linux, macOS, Windows).

## Contributing

We welcome contributions from the community! If you're interested in contributing, please refer to our [CONTRIBUTING.md](CONTRIBUTING.md) (coming soon) for guidelines on how to submit bug reports, feature requests, and pull requests.

## License

This project is licensed under the [MIT License](LICENSE.md).

## Contact

For any questions or inquiries, please open an issue on GitHub or contact the project maintainers.
