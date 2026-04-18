package transport

import (
	"github.com/google/gopacket"
)

// Transport extracts a raw SPA payload from a network packet.
type Transport interface {
    ExtractPayload(packet gopacket.Packet) ([]byte, bool)
}
