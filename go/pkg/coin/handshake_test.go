package coin

import (
	"encoding/json"
	"net"
	"sync"
	"testing"
)

func TestHandshake(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()

	var srvID string
	var srvErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		srvID, srvErr = HandleHandshake(b, "server")
	}()

	cliID, err := PerformHandshake(a, "client")
	if err != nil {
		t.Fatalf("handshake failed: %v", err)
	}
	wg.Wait()
	if srvErr != nil {
		t.Fatalf("server handshake failed: %v", srvErr)
	}
	if cliID != "server" || srvID != "client" {
		t.Fatalf("unexpected ids cli=%s srv=%s", cliID, srvID)
	}
}

func TestHandshakeMismatch(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()

	var wg sync.WaitGroup
	var srvErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		dec := json.NewDecoder(b)
		enc := json.NewEncoder(b)

		var req handshakeMessage
		if err := dec.Decode(&req); err != nil {
			srvErr = err
			return
		}
		resp := handshakeMessage{Protocol: P2PProtocolVersion + 1, NodeID: "server"}
		srvErr = enc.Encode(&resp)
	}()

	id, err := PerformHandshake(a, "client")
	if err == nil {
		t.Fatalf("expected error, got id=%s", id)
	}
	wg.Wait()
	if srvErr != nil {
		t.Fatalf("server error: %v", srvErr)
	}
}
