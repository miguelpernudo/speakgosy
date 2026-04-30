//go:build proxy

package proxy

import (
	"fmt"
	"io"
	"net"
)

type Proxy struct {
    ListenAddr string  // ":443"
    TargetAddr string  // "localhost:22"
}

// Forward accepts one connection from allowedIP and pipes it to TargetAddr.
func (p *Proxy) Forward(allowedIP string) error {
	listener, err := net.Listen("tcp", p.ListenAddr) // Open TCP port for listening.
	if err != nil {
	    return err
	}
	defer listener.Close()
	
	conn, err := listener.Accept() 						// Waits for the connection.
	if err != nil {
	    return err
	}

	addr := conn.RemoteAddr().(*net.TCPAddr)	// The IP address is verified.
	if addr.IP.String() != allowedIP {
    conn.Close()
    return fmt.Errorf("proxy: rejected connection from %s", addr.IP)
	}

	// Communication itself.
	target, err := net.Dial("tcp", p.TargetAddr)
	if err != nil {
	    conn.Close()
	    return err
	}
	
	go io.Copy(target, conn)
	io.Copy(conn, target)	

	return nil
}
