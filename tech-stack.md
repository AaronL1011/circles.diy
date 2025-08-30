# circles.diy Tech Stack

A Go-based decentralized social platform designed with DARPA principles: no single point of failure, verifiable data, and true portability.

## Core Architecture

### Identity Layer
- **DIDs (Decentralized Identifiers)** - Per-device keys for user identity
- **WebAuthn Passkeys** - "Continue with device" onboarding experience
- **DNS Mapping** - Human-readable handles resolve to DIDs
- **Key Recovery** - Offline keys + trusted contacts

### Data Model
- **Append-only Event Log** - All user actions as signed events
- **Merkle DAG** - Content-addressed, auditable state transitions
- **CRDT Conflict Resolution** - Offline-first merge capabilities
- **Content Addressing** - Hash-based object references

### Storage
- **Local Storage** - SQLite/IndexedDB client cache
- **Distributed Blobs** - User-chosen storage (phone, home server, commodity)
- **Multi-home Replication** - Multiple mirrors for resilience
- **Erasure Coding** - Cost-effective backup strategy

## Transport & Networking

### Primary Transport (P2P First)
- **libp2p Core** - Primary networking stack from day one
  - **QUIC/TCP Transports** - Reliable connection layer
  - **Noise/TLS Security** - Encrypted communications
  - **Kademlia DHT** - Distributed routing and discovery
  - **GossipSub** - Efficient message propagation
  - **Circuit Relay** - NAT traversal and connectivity
- **Go libp2p Host** - Native Go implementation
- **IPFS Integration** - Content addressing and routing

### Protocol Interoperability Layer
- **Translation Gateway Services** - Protocol bridging infrastructure
  - **AT Protocol Bridge** - BlueSky network compatibility
  - **ActivityPub Bridge** - Mastodon/Fediverse interoperability  
  - **Scuttlebutt Bridge** - SSB network integration
  - **Nostr Bridge** - Decentralized social relay compatibility
  - **XMTP Bridge** - Web3 messaging protocol support
- **Schema Translation** - Cross-protocol data format conversion
- **Identity Mapping** - DID ↔ Protocol-specific identity resolution
- **Content Adaptation** - Feature set normalization across protocols

### Federated Fallback Layer
- **Go HTTP Services** - Backup reliability when P2P fails
- **gRPC Services** - High-performance service-to-service communication
- **Postgres** - Index storage for aggregators
- **Object Storage** - Blob storage backend

### Infrastructure
- **Distributed Hash Table** - Peer discovery and content routing
- **Content Addressing** - Trust-free replication via IPFS CIDs
- **Peer Routing** - Automatic network topology optimization

## Service Architecture

### Aggregators (P2P Nodes)
- **Feed Computation** - Home, trending, etc. computed over libp2p network
- **Social Graph Materialization** - Follow relationships via DHT queries
- **Subscription Management** - User log monitoring via GossipSub
- **API Layer** - libp2p services + HTTP gateway for web clients
- **Event Streaming** - Real-time updates via libp2p streams

### Personal Data Servers (P2P Nodes)
- **Repo Hosting** - User's signed event log published to IPFS
- **Mirror Coordination** - Multi-peer replication via libp2p
- **DNS Integration** - Handle-to-DID resolution + DHT publishing
- **Sync Protocols** - libp2p streaming for real-time replication
- **Protocol Translation** - Real-time conversion of incoming/outgoing cross-network content

### Moderation Services (P2P Nodes)
- **Label Providers** - Signed content assertions distributed via GossipSub
- **Blocklist Services** - Community-driven filtering via DHT
- **Abuse Reporting** - Signed message system over libp2p streams
- **Cross-Protocol Moderation** - Synchronized moderation actions across bridge networks

### Protocol Bridge Nodes
- **Multi-Protocol Gateways** - Specialized nodes running multiple protocol stacks
- **Event Translation** - Real-time conversion between native and foreign event formats
- **Identity Synchronization** - Cross-network identity verification and mapping
- **Content Normalization** - Feature parity handling (replies, reactions, media, etc.)
- **Rate Limiting & Abuse Prevention** - Cross-network spam and abuse mitigation

## Client Stack

### Web Client
- **js-libp2p** - Browser P2P networking via WebRTC/WebTransport
- **Service Worker** - Offline-first operation + libp2p node
- **Local Cache** - SQLite/IndexedDB + IPFS block store
- **CRDT Sync** - Conflict resolution over libp2p streams
- **DHT Queries** - Direct peer discovery and content resolution
- **Protocol Adapters** - Client-side translation for cross-network interactions
- **Unified Feed** - Merged timeline from multiple protocol sources

### Key Management
- **Signature Generation** - All actions signed client-side
- **Key Rotation** - Built-in key lifecycle
- **Recovery Flows** - Offline + social recovery
- **Cross-Protocol Keys** - Identity bridging with protocol-specific keypairs

## Discovery & Indexing

### Name Resolution
- **DNS-based Handles** - Human-readable to DID mapping
- **Kademlia DHT** - Primary content and peer discovery
- **IPNS Records** - Mutable name resolution via libp2p
- **Directory Services** - Competing discovery providers as P2P nodes
- **Cross-Protocol Identity** - Unified identity resolution across networks
  - **AT Protocol DIDs** - BlueSky identity integration
  - **ActivityPub Actors** - Mastodon/Fediverse handle mapping  
  - **SSB Keys** - Scuttlebutt identity bridging
  - **Nostr Public Keys** - Nostr identity translation

### Search & Discovery
- **Distributed Indexers** - Multiple P2P indexing nodes
- **DHT Queries** - Direct peer-to-peer content discovery
- **Cross-Network Search** - Federated search across protocol bridges
- **Result Blending** - Client-side aggregation from multiple DHT sources and networks
- **Verifiable Results** - Signature-backed content with IPFS addressing
- **Protocol-Specific Crawlers** - Specialized indexing for each bridged network

## MVP Implementation Path

### Phase 1: P2P Foundation + Core Interop
- **Go libp2p Host** - Core networking and DHT participation
- **IPFS Content Store** - Content-addressed data layer
- **Basic P2P Aggregator** - Chronological feed via GossipSub
- **Simple Web Client** - js-libp2p + WebRTC connectivity
- **AT Protocol Bridge** - BlueSky network compatibility (highest adoption)
- **ActivityPub Bridge** - Mastodon/Fediverse interoperability (mature ecosystem)
- **Identity Translation Layer** - DID ↔ AT Protocol ↔ ActivityPub mapping

### Phase 2: Full P2P Network + Extended Interop
- **Advanced Routing** - Optimized DHT and peer discovery
- **Web Client Enhancement** - Full offline-first P2P operation with unified feeds
- **Moderation Layer** - Distributed content labeling with cross-protocol sync
- **DNS Integration** - Handle resolution + IPNS publishing
- **Scuttlebutt Bridge** - SSB network integration for offline-first communities
- **Nostr Bridge** - Decentralized relay network compatibility
- **Cross-Protocol Follows** - Follow users across different networks seamlessly

### Phase 3: Advanced Interop & Optimization
- **XMTP Bridge** - Web3 messaging protocol integration
- **Matrix Bridge** - Chat/communities protocol compatibility
- **Performance Tuning** - Connection pooling, caching, routing optimization
- **Advanced CRDTs** - Rich conflict resolution for complex cross-protocol data types
- **Protocol Plugin System** - Extensible architecture for future protocol integration
- **Multi-Protocol Identity Verification** - Cross-network reputation and verification
- **Federation Fallback** - HTTP services for non-P2P clients

## Operational Considerations

### Security
- **Minimal Trust Surface** - Signatures everywhere, no trusted intermediaries
- **Peer Authentication** - Cryptographic peer identity verification
- **Content Verification** - IPFS hashing + signature validation
- **Key Rotation** - Automated key lifecycle with DHT publication

### Scalability
- **DHT Scalability** - Logarithmic routing, unlimited peer capacity
- **Content Addressing** - Efficient caching and deduplication
- **Peer Distribution** - Global P2P network with natural load balancing
- **Bandwidth Optimization** - BitSwap protocol for efficient data exchange

### Business Model
- **P2P Infrastructure** - Managed libp2p bootstrap and relay nodes
- **Exit Guarantees** - Data always portable via IPFS/DHT + protocol bridges
- **Premium Services** - Enhanced peer connectivity and storage pinning
- **Protocol Bridge Services** - Managed translation nodes for enterprise/communities
- **Gateway Services** - HTTP bridges for legacy client compatibility
- **Multi-Network Identity** - Premium cross-protocol identity verification