package crawler

import (
	"log"
	"net"

	"pila/pkg/coin"
)

// Peer represents a remote peer connection.
type Peer struct {
	Conn net.Conn
	ID   string
}

// Connect dials the address and performs a handshake using nodeID.
func Connect(addr, nodeID string) (*Peer, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	id, err := coin.PerformHandshake(c, nodeID)
	if err != nil {
		c.Close()
		return nil, err
	}
	return &Peer{Conn: c, ID: id}, nil
}

// ListenAndServe listens on addr and handles incoming handshake connections.
// For each successful connection, the remote node ID is logged.
func ListenAndServe(addr, nodeID string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()
				id, err := coin.HandleHandshake(conn, nodeID)
				if err != nil {
					log.Printf("handshake error: %v", err)
					return
				}
				log.Printf("connected peer %s", id)
			}(c)
		}
	}()
	return ln, nil
}
