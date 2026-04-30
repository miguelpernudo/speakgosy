# SpeakGosy

# WORK IN PROGRESS

> *"In an internet of bots and mass scanning, your ports are prohibited. You must speak the word to *go* in."*

Secure and lightweight Single Packet Authorization (SPA) daemon written in Go. 
Keeps a server with zero open ports until you send it the encrypted whisper. 
Evades restrictive networks, and is immune to port scanning, brute force, and 
replay attacks.

## How it works

Your server has a firewall that blocks all ports and drops all incoming traffic.
SpeakGosy remains active, monitoring the network.

## How it works

    CLIENT                                             SERVER
    ──────                                             ──────
generates nonce + TOTP               
  (totp/totp.go)                      
  
       │                            

builds Payload {ip, nonce, TOTP}     
  (packet/payload.go)                 

       │                            

    encrypts it                             
  (crypt/cipher.go)                   

       │                            

 encode + sends UDP                             
(transport/dns.go & raw.go)       ────►          captures the packet
                                            (transport/dns.go & raw.go)
  
                                                          │
                                                      
                                                  decrypts and parse
                                        (crypt/cipher.go & packet/parser.go)
                                        
                                                          │
                                                      
                                              verifies nonce, for anti-replay
                                                  (replay/cache.go)
                                                  
                                                          │
                                                      
                                                     verify TOTP
                                                    (totp/totp.go)
                                                    
                                                          │
                                              ┌───────────┴───────────┐
                                           firewall                 proxy
                                        (nftables.go)             (proxy.go)
                                        
                                           AllowIP                  Forward

                                              │                        │

                                           goroutine               TCP tunnel
                                           RevokeIP                until EOF
                                            (TTL)


## Evading restrictive networks

When you're connected to certain public networks, such as the university's Wi-Fi, If you try to connect to your server via SSH, the firewall on that network will most likely block that type of outbound traffic, so you won't be able to work from outside.

SpeakGosy tries to solve this problem by masquerading as a DNS query,
and if that fails, it tries via HTTPS. In the firewall mode, the SSH connection is also made through 
port 443, avoiding using the prohibited port 22 in this public networks.

DNS 53 UDP: Payload disguise as a valid DNS query.
Raw 443 UDP: Passes as QUIC traffic (no disguise).

## Security properties

Zero ports open:
Nothing to scan, nothing to brute force.

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
