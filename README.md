# SpeakGosy

# WORK IN PROGRESS:
## The structure of `internal/` is there, but it needs to be polished, checked to make sure it compiles and everything works, and comments need to be added. The `main` functions are missing.


> *"In an internet of bots and mass scanning, your ports are prohibited. You must speak the word to *go* in."*

Secure and lightweight Single Packet Authorization (SPA) daemon written in Go. 
Keeps a server with zero open ports until you send it the encrypted whisper. 
Evades restrictive networks, and is immune to port scanning, brute force, and 
replay attacks.

## How it works

Your server has a firewall that blocks all ports and drops all incoming traffic.
SpeakGosy remains active, monitoring the network.

When you want to connect:

1. The client builds a packet containing your IP, a one-time nonce and a TOTP code.
2. It encrypts everything with ChaCha20-Poly1305 and sends it as a UDP packet.
3. The firewall drops the packet, but SpeakGosy saw it.
4. SpeakGosy decrypts, verifies the nonce and TOTP.
5. If valid, it opens the firewall for your IP only and you can SSH in.


## Evading restrictive networks

At my university, when I'm connected to the Wi-Fi, it blocks me from connecting
to my server. SpeakGosy tries to solve this problem by masquerading as a DNS query,
and if that fails, it tries via HTTPS. The SSH connection is also made through 
port 443, avoiding using the prohibited port 22 in this public networks.

DNS 53 UDP: Payload encoded in the QNAME of a valid DNS query.
Raw 443 UDP: Passes as QUIC (HTTP/3) traffic, no disguise needed.

## Security properties

Zero ports open:
Nothing to scan, nothing to brute force

Authenticated encryption:
ChaCha20-Poly1305 ensures confidentiality and integrity.

Anti-replay:
Nonces are tracked with TTL matching the TOTP window.

Anti-expiry:
TOTP codes are valid for 30 seconds (±1 window for clock skew)

IP binding:
The firewall rule is opened only for the IP inside the encrypted payload.

No root required:
Runs with Linux capabilities `cap_net_raw` and `cap_net_admin`

## Structure

```
speakgosy/
├── cmd/
│   ├── client/
│   │   └── main.go          sends the SPA packet
│   └── server/
│       └── main.go          passive watcher, firewall manager
├── internal/
│   ├── firewall/
│   │   ├── manager.go       FirewallManager interface
│   │   └── nftables.go      nftables adapter
│   ├── replay/
│   │   └── cache.go         nonce cache
│   ├── totp/
│   │   └── totp.go          TOTP generation and validation
│   ├── crypt/
│   │   └── cipher.go        ChaCha20-Poly1305 encrypt/decrypt
│   ├── packet/
│   │   ├── payload.go       SPA payload struct
│   │   ├── builder.go       serialize and encrypt payload
│   │   └── parser.go        decrypt and deserialize payload
│   └── transport/
│       ├── transport.go     Transport interface
│       ├── dns.go           DNS transport (port 53, QNAME encoding)
│       └── raw.go           Raw UDP transport (port 443)
├── go.mod
└── go.sum
```

## Building
```bash
# Build both binaries
go build -o bin/spkgo-server ./cmd/server/
go build -o bin/spkgo-client ./cmd/client/

sudo setcap cap_net_raw,cap_net_admin+eip ./bin/spkgo-server
```
