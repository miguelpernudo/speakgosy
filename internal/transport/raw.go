package transport

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type RawTransport struct {
	Port uint16
}

func NewRawTransport(port uint16) *RawTransport {
	return &RawTransport{Port: port}
}

func (r *RawTransport) ExtractPayload(packet gopacket.Packet) ([]byte, bool) {
	//Extracts the UDP layer.
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return nil, false
	}
	udp, _ := udpLayer.(*layers.UDP)


	if len(udp.Payload) == 0 {
    return nil, false
	}

	if uint16(udp.DstPort) != r.Port {
	   return nil, false
	}

	return udp.Payload, true

}
