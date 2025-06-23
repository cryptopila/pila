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
