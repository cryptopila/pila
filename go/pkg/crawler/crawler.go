package crawler

import (
	"log"
	"net"
	"sync"

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

// Crawler manages outbound peer connections.
type Crawler struct {
	NodeID string

	mu    sync.Mutex
	peers map[string]*Peer
}

// New returns a new crawler instance.
func New(nodeID string) *Crawler {
	return &Crawler{NodeID: nodeID, peers: make(map[string]*Peer)}
}

// Connect adds a new peer connection to addr and stores it on success.
func (c *Crawler) Connect(addr string) (*Peer, error) {
	p, err := Connect(addr, c.NodeID)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	c.peers[addr] = p
	c.mu.Unlock()
	return p, nil
}

// Close terminates all peer connections.
func (c *Crawler) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for addr, p := range c.peers {
		_ = p.Conn.Close()
		delete(c.peers, addr)
	}
}
