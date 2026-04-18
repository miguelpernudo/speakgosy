package transport

import (
	"strings"
	"encoding/base32"
	
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type DNSTransport struct {
	Domain string
}

func NewDNSTransport(domain string) *DNSTransport {
	return &DNSTransport{Domain: domain}
}


func (d *DNSTransport) ExtractPayload(packet gopacket.Packet) ([]byte, bool) {
	// Extracts the DNS layer.
	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer == nil {
    return nil, false
	}
	dns, _ := dnsLayer.(*layers.DNS) 

	// Verifies if it's a query.
	if dns.QR {  
    return nil, false
	}

	// Extracts the QNAME.
	if len(dns.Questions) == 0 {
    return nil, false
	}
	qname := string(dns.Questions[0].Name)

	// Verifies if it ends in the correct domain we are using.
	suffix := d.Domain
	if !strings.HasSuffix(qname, suffix) {
    return nil, false
	}
	sub := strings.TrimSuffix(qname, suffix) 

	// Decode.
	data, err := base32.StdEncoding.DecodeString(strings.ToUpper(sub))
	if err != nil {
  	return nil, false
	}

return data, true
}
