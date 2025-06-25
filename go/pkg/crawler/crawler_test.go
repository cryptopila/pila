package crawler

import (
	"net"
	"testing"
)

func TestConnect(t *testing.T) {
	ln, err := ListenAndServe("127.0.0.1:0", "server")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	peer, err := Connect(ln.Addr().String(), "client")
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer peer.Conn.Close()

	if peer.ID != "server" {
		t.Fatalf("expected server id, got %s", peer.ID)
	}
}

func TestListenAndServeHandshakeError(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	go func() {
		c, _ := ln.Accept()
		defer c.Close()
		// send invalid handshake
		c.Write([]byte("invalid"))
	}()

	_, err = Connect(ln.Addr().String(), "client")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCrawlerConnectAndClose(t *testing.T) {
	ln, err := ListenAndServe("127.0.0.1:0", "server")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	c := New("client")
	peer, err := c.Connect(ln.Addr().String())
	if err != nil {
		t.Fatalf("connect: %v", err)
	}

	if len(c.peers) != 1 {
		t.Fatalf("expected 1 peer, got %d", len(c.peers))
	}
	if peer.ID != "server" {
		t.Fatalf("unexpected peer id %s", peer.ID)
	}

	c.Close()
	if _, err := peer.Conn.Write([]byte{0}); err == nil {
		t.Fatalf("connection should be closed")
	}
}
